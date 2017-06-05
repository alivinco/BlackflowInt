package restmqttproxy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"encoding/base64"
	"io/ioutil"
	"time"
	"net/http"
)

type Process struct {
	mqttTransport *fimpgo.MqttTransport
	Config      *ProcessConfig
	State       string
	LastError   string
	RunningRequest []RunningRequest
}

// NewProcess is a constructor
func NewProcess(config *ProcessConfig) *Process {
	proc := Process{Config: config}
	proc.State = "LOADED"
	return &proc
}


// Init doing the process bootrstrap .
func (pr *Process) Init() error {
	pr.State = "INIT_FAILED"
	log.Info("Initializing MQTT adapter. Broker address : ",pr.Config.MqttBrokerAddr)
	//"tcp://localhost:1883", "blackflowint", "", ""
	pr.mqttTransport = fimpgo.NewMqttTransport(pr.Config.MqttBrokerAddr, pr.Config.MqttClientID, pr.Config.MqttBrokerUsername, pr.Config.MqttBrokerPassword,true,0,0)
	pr.mqttTransport.SetMessageHandler(pr.OnMessage)
	log.Info("MQTT adapter initialization completed.")

	pr.State = "INITIALIZED"
	return nil
}


func (pr *Process) ConfigureHttpServer () {
	log.Info("Starting HTTP proxy")
	secret, err := base64.URLEncoding.DecodeString(pr.Config.ClientSecret)
	if err != nil {
		log.Error("Fail decode secret :",err)
	}
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Restricted group
	r := e.Group("/bfint/rest-mqtt-proxy")
	r.Use(middleware.JWT(secret))
	r.POST("/publish", pr.httpRequestHandler)
	e.Start(pr.Config.ProxyBindHost)
	log.Info("HTTP proxy started")

}

func (pr *Process) registerRequest(responseTopic string , responseMsgType string) chan []byte {
	respChan := make(chan []byte)
	runReq := RunningRequest{RespTopic:responseTopic,RespMsgType:responseMsgType,ResponseChan:respChan}
	pr.RunningRequest = append(pr.RunningRequest,runReq)
	//log.Info("Register.Number of in flight transactions = ",len(pr.RunningRequest))
	return respChan
}

func (pr *Process) unregisterRequest(responseTopic string , responseMsgType string) {
	var result []RunningRequest
	for i := range pr.RunningRequest {
		if pr.RunningRequest[i].RespTopic == responseTopic && pr.RunningRequest[i].RespMsgType == responseMsgType {
			result = append(pr.RunningRequest[:i],pr.RunningRequest[i+1:]...)
			break
		}
	}
	if result != nil {
		pr.RunningRequest = result
	}
	//log.Info("Unregister .Number of in flight transactions = ",len(pr.RunningRequest))
}

func (pr *Process) httpRequestHandler(c echo.Context) error {
	log.Info("New request ")
	user := c.Get("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	//domains , ok := claims["app_metadata"].(map[string][]string)
	username , ok := claims["email"].(string)
	if ok {
		log.Info("Authenticated domains ",username)
	}
	requestTopic := c.Request().Header.Get("ReqTopic")
	responseTopic := c.Request().Header.Get("RespTopic")
	responseMsgType := c.Request().Header.Get("RespMsgType")
	domain := ""
	for i := range pr.Config.UserSettings {
		if pr.Config.UserSettings[i].Name == username {
			domain = pr.Config.UserSettings[i].ActiveDomain;
		}
	}
	log.Debug("Request topic  : ",requestTopic)
	log.Debug("Response topic : ",responseTopic)
	log.Debug("ResponseMsgType: ",responseMsgType)
	var body []byte
	var err error
	if requestTopic != "" && (domain != "" || pr.Config.SkipDomain) {
		// Publish and forget
		if domain != ""{
			requestTopic = domain+"/"+requestTopic
		}
		body , err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Error("Request body is empty")
			return err
		}

	}else {
		return nil
	}
	if responseTopic != "" {
		// Publish and wait for response
		if domain != "" {
			responseTopic = domain + "/" + responseTopic
		}
		fimpResponseCh := pr.registerRequest(responseTopic,responseMsgType)
		err := pr.mqttTransport.Subscribe(responseTopic)
		if err != nil {
			log.Debug(err)
		}
		log.Debugf("Publishing message to topic : %s , reply topic: %s , msg : %s",requestTopic,responseTopic,string(body) )
		pr.mqttTransport.PublishRaw(requestTopic,body)
		select {
		case fimpResponse := <- fimpResponseCh:
			c.JSONBlob(http.StatusOK,fimpResponse)

		case <- time.After(time.Second*10):
			log.Info("No response from queue for 10 seconds")

		}
		pr.mqttTransport.Unsubscribe(responseTopic)

	}else {
		log.Debugf("Publishing message to topic : %s , msg : %s",requestTopic,string(body) )
		pr.mqttTransport.PublishRaw(requestTopic,body)
	}


	return nil
}



// OnMessage is invoked by an adapter on every new message
// The code is executed in callers goroutine
func (pr *Process) OnMessage(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage , rawMessage []byte) {
	log.Debugf("New message from topic %s",topic)
	for i := range pr.RunningRequest {
		if pr.RunningRequest[i].RespMsgType == iotMsg.Type && pr.RunningRequest[i].RespTopic == topic {
			pr.RunningRequest[i].ResponseChan <- rawMessage
			pr.unregisterRequest(topic,iotMsg.Type)
		}
	}

}


// Start starts the process by starting MQTT adapter ,
// starting scheduler
func (pr *Process) Start() error {
	// try to initialize process first if current state is not INITIALIZED
	if pr.State == "INIT_FAILED" || pr.State == "LOADED" {
		if err := pr.Init(); err != nil {
			return err
		}
	}
	err := pr.mqttTransport.Start()
	if err != nil {
		log.Fatalln("Error: ", err)
		return err
	}
	go pr.ConfigureHttpServer()
	pr.State = "RUNNING"
	return nil

}

// Stop stops the process by unsubscribing from all topics ,
// stops scheduler and stops adapter.
func (pr *Process) Stop() error {
	pr.mqttTransport.Stop()
	pr.State = "STOPPED"
	return nil
}

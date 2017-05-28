package restmqttproxy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"encoding/base64"
)

type Process struct {
	mqttTransport *fimpgo.MqttTransport
	Config      *ProcessConfig
	State       string
	LastError   string
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
	log.Info("Initializing MQTT adapter.")
	//"tcp://localhost:1883", "blackflowint", "", ""
	pr.mqttTransport = fimpgo.NewMqttTransport(pr.Config.MqttBrokerAddr, pr.Config.MqttClientID, pr.Config.MqttBrokerUsername, pr.Config.MqttBrokerPassword,true,0,0)
	pr.mqttTransport.SetMessageHandler(pr.OnMessage)
	log.Info("MQTT adapter initialization completed.")

	pr.State = "INITIALIZED"
	return nil
}


func (pr *Process) ConfigureHttpServer () {
	log.Info("Starting HTTP proxy")
	secret64 := ""
	secret, err := base64.URLEncoding.DecodeString(secret64)
	if err != nil {
		log.Error("Fail decode secret :",err)
	}
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Restricted group
	r := e.Group("/rest-mqtt-proxy")
	r.Use(middleware.JWT(secret))
	r.POST("/publish", pr.httpRequestHandler)
	e.Start(pr.Config.ProxyBindHost)
	log.Info("HTTP proxy started")

}

func (pr *Process) httpRequestHandler(c echo.Context) error {
	log.Info("New request ")
	user := c.Get("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	domains , ok := claims["app_metadata"].(map[string][]string)
	if ok {
		log.Info("Authenticated domains ",domains)
	}
	requestTopic := c.Request().Header.Get("FimpReqTopic")
	responseTopic := c.Request().Header.Get("FimpRespTopic")
	responseMsgType := c.Request().Header.Get("FimpRespMsgType")
	log.Info("Request topic  : ",requestTopic)
	log.Info("Response topic : ",responseTopic)
	log.Info("ResponseMsgType: ",responseMsgType)
	return nil
}



// OnMessage is invoked by an adapter on every new message
// The code is executed in callers goroutine
func (pr *Process) OnMessage(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage) {


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

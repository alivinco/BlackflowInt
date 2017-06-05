package restmqttproxy

import (
	"sync"
	"github.com/labstack/echo"
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/models"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/spf13/viper"
	"strings"
)

type Integration struct {
	processes []*Process
	// in memmory copy of config file
	processConfigs  []ProcessConfig
	StoreLocation   string
	storeFullPath   string
	Name            string
	configSaveMutex *sync.Mutex
}
// GetProcessByID returns process by it's ID
func (it *Integration) GetProcessByID(ID IDt) *Process {
	for i := range it.processes {
		if it.processes[i].Config.ID == ID {
			return it.processes[i]
		}
	}
	return nil
}

// Init initilizes integration app
func (it *Integration) Init() {
	it.storeFullPath = filepath.Join(it.StoreLocation, it.Name+".json")
}
func (it *Integration) LoadConfig() error {
	log.Info("Loading config from ",it.storeFullPath)
	// Check if configuration are set via ENV variable .
	if  viper.GetString("mqtt_broker_addr") != "" {
		it.EnvVarConfig()
		return nil
	}

	if _, err := os.Stat(it.storeFullPath); os.IsNotExist(err) {

		log.Info("Integration configuration is loaded from default.")
		payload, err := json.Marshal(it.processConfigs)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(it.storeFullPath, payload, 0777)
		panic("No config file , template was created ")
	}

	payload, err := ioutil.ReadFile(it.storeFullPath)
	if err != nil {
		log.Errorf("Integration can't load configuration file from %s, Errro:%s", it.storeFullPath, err)
		return err
	}
	err = json.Unmarshal(payload, &it.processConfigs)
	if err != nil {
		log.Error("Can't load the integration cofig.Unmarshall error :", err)
	}

	return nil
}

// BrokerAutoConfig configures broker using ENV variables set by BlackTower
func (it *Integration) EnvVarConfig() {
	log.Info("Loading configs from ENV variables")
	conf := ProcessConfig{ID:1}
	conf.MqttBrokerAddr = "tcp://" + viper.GetString("mqtt_broker_addr")
	conf.MqttBrokerUsername = viper.GetString("mqtt_username")
	conf.MqttBrokerPassword = viper.GetString("mqtt_password")
	//conf.MqttClientID = fmt.Sprintf("%s-%d", viper.GetString("mqtt_clientid"), 1)
	conf.MqttClientID = "rest-mqtt-proxy-1"
	conf.ClientSecret = viper.GetString("int_app_auth_client_secret")
	conf.ProxyBindHost = viper.GetString("int_proxy_bind_host")
	if conf.ProxyBindHost == "" {
		conf.ProxyBindHost = ":5050"
	}
	// format username:active_domain,username:active_domain
	userConfigs := strings.Split(viper.GetString("int_user_configs"),",")
	log.Info("Use configs = ",viper.GetString("int_user_configs"))
	for _ , uconf := range userConfigs {
		userSet := UserSettings{}
		// format username:active_domain
		confKeyVal := strings.Split(uconf,":")
		userSet.Name = confKeyVal[0]
		userSet.ActiveDomain = confKeyVal[1]
		conf.UserSettings = append(conf.UserSettings,userSet)
	}
	it.processConfigs = append(it.processConfigs,conf)

}

func (it *Integration) InitProcesses() error {
	it.Init()
	it.LoadConfig()
	log.Info("Process configs :",it.processConfigs)
	proc := NewProcess(&it.processConfigs[0])
	it.processes = append(it.processes,proc)
	return  proc.Start()
}



// Boot initializes integration
func Boot(mainConfig *models.MainConfig, restHandler *echo.Echo) *Integration {
	log.Info("Booting RestMqttProxy integration ")
	integr := Integration{Name: "rest-mqtt-proxy", StoreLocation: mainConfig.StorageLocation}
	integr.InitProcesses()
	//if restHandler != nil {
	//	restAPI := IntegrationAPIRestEndp{&integr, restHandler}
	//	restAPI.SetupRoutes()
	//}
	return &integr
}
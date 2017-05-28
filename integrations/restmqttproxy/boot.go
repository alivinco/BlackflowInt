package restmqttproxy

import (
	"sync"
	"github.com/labstack/echo"
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/models"
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

func (it *Integration) InitProcesses() error {
	conf := ProcessConfig{Name:"test",
		MqttBrokerAddr:"tcp://localhost:1883",
		MqttClientID:"rest-mqtt-proxy-proc-1",
		ProxyBindHost:":5050"}
	proc := NewProcess(&conf)
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
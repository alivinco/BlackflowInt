package influxdb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/models"
	"github.com/labstack/echo"
)

// Integration is root level container
type Integration struct {
	processes []*Process
	// in memmory copy of config file
	processConfigs []ProcessConfig
	StoreLocation  string
	storeFullPath  string
	Name           string
}

// GetProcessByID returns process by it's ID
func (it *Integration) GetProcessByID(ID IDt) *Process {
	for _, proc := range it.processes {
		if proc.Config.ID == ID {
			return proc
		}
	}
	return nil
}

// GetDefaultIntegrConfig returns default config .
func (it *Integration) GetDefaultIntegrConfig() []ProcessConfig {
	selector := []Selector{
		Selector{Topic: "*/jim1/evt*"},
	}
	filters := []Filter{
		Filter{
			ID:       1,
			MsgClass: "sensor",
			IsAtomic: true,
		},
		Filter{
			ID:       2,
			MsgClass: "binary",
			IsAtomic: true,
		},
	}
	config := ProcessConfig{
		MqttBrokerAddr:     "tcp://localhost:1883",
		MqttBrokerUsername: "",
		MqttBrokerPassword: "",
		MqttClientID:       "blackflowint_sub_test",
		InfluxAddr:         "http://localhost:8086",
		InfluxUsername:     "",
		InfluxPassword:     "",
		InfluxDB:           "iotmsg_test",
		BatchMaxSize:       1000,
		SaveInterval:       1000,
		Filters:            filters,
		Selectors:          selector,
	}

	return []ProcessConfig{config}

}

// LoadConfig loads configs from json file and saves it into ProcessConfigs
func (it *Integration) LoadConfig() error {
	if _, err := os.Stat(it.storeFullPath); os.IsNotExist(err) {
		it.processConfigs = it.GetDefaultIntegrConfig()
		log.Info("Integration configuration is loaded from default.")
		return it.SaveConfigs()
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
	return err

}

// SaveConfigs saves configs to json file
func (it *Integration) SaveConfigs() error {
	for _, prc := range it.processes {
		for _, prcConf := range it.processConfigs {
			if prc.Config.ID == prcConf.ID {
				prcConf.Filters = prc.GetFilters()
				prcConf.Selectors = prc.GetSelectors()
			}
		}
	}
	payload, err := json.Marshal(it.processConfigs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(it.storeFullPath, payload, 0777)
}

// InitProcesses starts processes based on ProcessConfigs
func (it *Integration) InitProcesses(autoStart bool) error {
	if it.processConfigs == nil {
		return errors.New("Load configurations first.")
	}
	for _, procConf := range it.processConfigs {
		proc := NewProcess(&procConf)
		proc.Init()
		log.Infof("Process ID=%d was initialized.", procConf.ID)
		if autoStart {
			err := proc.Start()
			if err != nil {
				log.Errorf("Process ID=%d failed to start . Error : %s", procConf, err)
			}
		}
		it.processes = append(it.processes, proc)
	}
	return nil
}

// Boot initializes integration
func Boot(mainConfig *models.MainConfig, restHandler *echo.Echo) *Integration {
	log.Info("Booting InfluxDB integration ")
	integr := Integration{Name: "influxdb", StoreLocation: mainConfig.StorageLocation}
	integr.storeFullPath = filepath.Join(integr.StoreLocation, integr.Name+".json")
	integr.LoadConfig()
	integr.InitProcesses(true)
	return &integr
}

package influxdb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/models"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

// Integration is root level container
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

// GetDefaultIntegrConfig returns default config .
func (it *Integration) GetDefaultIntegrConfig() []ProcessConfig {

	selector := []Selector{
		Selector{ID: 1, Topic: "*/jim1/evt*"},
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

	measurements := []Measurement{
		Measurement{
			Name: "sensor",
			RetentionPolicyDuration: "8w",
			RetentionPolicyName:     "bf_sensor",
		},
		Measurement{
			Name: "level",
			RetentionPolicyDuration: "8w",
			RetentionPolicyName:     "bf_level",
		},
		Measurement{
			Name: "binary",
			RetentionPolicyDuration: "8w",
			RetentionPolicyName:     "bf_binary",
		},
		Measurement{
			Name: "default",
			RetentionPolicyDuration: "8w",
			RetentionPolicyName:     "bf_default",
		},
	}

	config := ProcessConfig{
		ID:                 1,
		MqttBrokerAddr:     "tcp://" + viper.GetString("mqtt_broker_addr"),
		MqttBrokerUsername: viper.GetString("mqtt_username"),
		MqttBrokerPassword: viper.GetString("mqtt_password"),
		MqttClientID:       viper.GetString("mqtt_clientid") + "-1",
		InfluxAddr:         "http://localhost:8086",
		InfluxUsername:     "",
		InfluxPassword:     "",
		InfluxDB:           "iotmsg_test",
		BatchMaxSize:       1000,
		SaveInterval:       1000,
		Filters:            filters,
		Selectors:          selector,
		Measurements:       measurements,
	}

	return []ProcessConfig{config}

}

// BrokerAutoConfig configures broker using ENV variables set by BlackTowe
func (it *Integration) BrokerAutoConfig(procID IDt) {
	proc := it.GetProcessByID(procID)
	proc.Config.MqttBrokerAddr = "tcp://" + viper.GetString("mqtt_broker_addr")
	proc.Config.MqttBrokerUsername = viper.GetString("mqtt_username")
	proc.Config.MqttBrokerPassword = viper.GetString("mqtt_password")
	proc.Config.MqttClientID = fmt.Sprintf("%s-%d", viper.GetString("mqtt_clientid"), procID)
	it.SaveConfigs()
}

// LoadConfig loads configs from json file and saves it into ProcessConfigs
func (it *Integration) LoadConfig() error {
	// ENV variables bindig.
	viper.AutomaticEnv()
	viper.SetEnvPrefix("zm")
	viper.SetDefault("mqtt_broker_addr", "localhost:1883")
	viper.SetDefault("mqtt_username", "")
	viper.SetDefault("mqtt_password", "")
	viper.SetDefault("mqtt_clientid", "bfint-influxdb")

	if it.configSaveMutex == nil {
		it.configSaveMutex = &sync.Mutex{}
	}
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
	// for _, prc := range it.processes {
	// 	for _, prcConf := range it.processConfigs {
	// 		if prc.Config.ID == prcConf.ID {
	// 			prcConf.Filters = prc.GetFilters()
	// 			prcConf.Selectors = prc.GetSelectors()
	// 		}
	// 	}
	// }
	it.configSaveMutex.Lock()
	defer func() {
		it.configSaveMutex.Unlock()
	}()
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
	for i := range it.processConfigs {
		proc := NewProcess(&it.processConfigs[i])
		proc.Init()
		log.Infof("Process ID=%d was initialized.", it.processConfigs[i].ID)
		if autoStart {
			err := proc.Start()
			if err != nil {
				log.Errorf("Process ID=%d failed to start . Error : %s", it.processConfigs[i], err)
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
	restAPI := IntegrationAPIRestEndp{&integr, restHandler}
	restAPI.SetupRoutes()
	return &integr
}

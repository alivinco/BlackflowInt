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
		InfluxDB:           "iotmsg",
		BatchMaxSize:       1000,
		SaveInterval:       1000,
		Filters:            filters,
		Selectors:          selector,
		Measurements:       measurements,
	}

	return []ProcessConfig{config}

}

// Init initilizes integration app
func (it *Integration) Init() {
	it.storeFullPath = filepath.Join(it.StoreLocation, it.Name+".json")
}

// BrokerAutoConfig configures broker using ENV variables set by BlackTower
func (it *Integration) BrokerAutoConfig(procID IDt) {
	proc := it.GetProcessByID(procID)
	proc.Config.MqttBrokerAddr = "tcp://" + viper.GetString("mqtt_broker_addr")
	proc.Config.MqttBrokerUsername = viper.GetString("mqtt_username")
	proc.Config.MqttBrokerPassword = viper.GetString("mqtt_password")
	proc.Config.MqttClientID = fmt.Sprintf("%s-%d", viper.GetString("mqtt_clientid"), procID)
	it.SaveConfigs()
}

// SetConfig config setter
func (it *Integration) SetConfig(processConfigs []ProcessConfig) {
	it.processConfigs = processConfigs
}

// UpdateProcConfig update process configurations
func (it *Integration) UpdateProcConfig(ID IDt, procConfig ProcessConfig, doRestart bool) error {
	proc := it.GetProcessByID(ID)
	err := proc.Configure(procConfig, doRestart)
	if err != nil {
		return err
	}
	err = it.SaveConfigs()
	return err
}

// LoadConfig loads integration configs from json file and saves it into ProcessConfigs
func (it *Integration) LoadConfig() error {
	// ENV variables bindig.
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
	if it.StoreLocation != "" {

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
	log.Info("Save to disk was skipped , StoreLocation is empty")
	return nil
}

// InitProcesses loads and starts ALL processes based on ProcessConfigs
func (it *Integration) InitProcesses() error {
	if it.processConfigs == nil {
		return errors.New("Load configurations first.")
	}
	for i := range it.processConfigs {
		it.InitNewProcess(&it.processConfigs[i])
	}
	return nil
}

// InitNewProcess initialize and start single process
func (it *Integration) InitNewProcess(procConfig *ProcessConfig) error {

	proc := NewProcess(procConfig)
	it.processes = append(it.processes, proc)
	if procConfig.Autostart {
		err := proc.Init()
		if err == nil {
			log.Infof("Process ID=%d was initialized.", procConfig.ID)
			err := proc.Start()
			if err != nil {
				log.Errorf("Process ID=%d failed to start . Error : %s", procConfig, err)
			}

		} else {
			log.Errorf("Initialization of Process ID=%d FAILED . Error : %s", procConfig.ID, err)
			return err
		}
	}
	return nil
}

// AddProcess adds new process .
func (it *Integration) AddProcess(procConfig ProcessConfig) (IDt, error) {
	defaultProc := it.GetDefaultIntegrConfig()
	procConfig.ID = GetNewID(it.processConfigs)
	if len(procConfig.Filters) == 0 {
		procConfig.Filters = defaultProc[0].Filters
	}
	if len(procConfig.Selectors) == 0 {
		procConfig.Selectors = defaultProc[0].Selectors
	}
	if len(procConfig.Measurements) == 0 {
		procConfig.Measurements = defaultProc[0].Measurements
	}
	it.processConfigs = append(it.processConfigs, procConfig)
	it.SaveConfigs()
	return procConfig.ID, it.InitNewProcess(&procConfig)
}

// RemoveProcess stops process , removes it from config file and removes instance .
func (it *Integration) RemoveProcess(ID IDt) error {
	var err error
	// removing process instance
	for i := range it.processes {
		if it.processes[i].Config.ID == ID {
			err = it.processes[i].Stop()
			it.processes = append(it.processes[:i], it.processes[i+1:]...)
			break
		}
	}
	// removing from config file
	for ic := range it.processConfigs {
		if it.processConfigs[ic].ID == ID {
			it.processConfigs = append(it.processConfigs[:ic], it.processConfigs[ic+1:]...)
			break
		}
	}
	if err == nil {
		it.SaveConfigs()

	}
	return err
}

// Boot initializes integration
func Boot(mainConfig *models.MainConfig, restHandler *echo.Echo) *Integration {
	log.Info("Booting InfluxDB integration ")
	integr := Integration{Name: "influxdb", StoreLocation: mainConfig.StorageLocation}
	integr.Init()
	integr.LoadConfig()
	integr.InitProcesses()
	if restHandler != nil {
		restAPI := IntegrationAPIRestEndp{&integr, restHandler}
		restAPI.SetupRoutes()
	}
	return &integr
}

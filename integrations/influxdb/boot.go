package influxdb

import "io/ioutil"
import "path/filepath"
import "encoding/json"

// Integration is root level container
type Integration struct {
	Processes      []Process
	ProcessConfigs ProcessConfigs
	StoreLocation  string
	storeFullPath  string
	Name           string
}

func (it *Integration) LoadConfig() error {
	payload, err := ioutil.ReadFile(it.storeFullPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(payload, it.ProcessConfigs)

	return err

}
func (it *Integration) SaveConfigs() error {
	payload, err := json.Marshal(it.ProcessConfigs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(it.storeFullPath, payload, 0777)
}

func (it *Integration) boot() {
	it.Name = "influxdb"
	it.storeFullPath = filepath.Join(it.StoreLocation, it.Name+".json")
}

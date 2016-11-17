package influxdb

import (
	"time"

	"reflect"

	"github.com/alivinco/iotmsglibgo"
	influx "github.com/influxdata/influxdb/client/v2"
)

// Transform defines function which converts IotMsg into influx data point
type Transform func(topic string, iotMsg *iotmsglibgo.IotMsg, domain string) (*influx.Point, error)

// IDt defines type of struct ID
type IDt int

// Selector defines message selector.
type Selector struct {
	ID        IDt
	ProcessID IDt
	Topic     string
}

// Filter defines message filter.
// Emty string - means no filter.
type Filter struct {
	ID          IDt
	ProcessID   IDt
	Name        string
	Topic       string
	Domain      string
	MsgType     string
	MsgClass    string
	MsgSubClass string
	// If true then returns everythin except matching value
	Negation bool
	// Boolean operation between 2 filters , supported values : and , or
	LinkedFilterBooleanOperation string
	LinkedFilterID               IDt
	IsAtomic                     bool
}

// ProcessConfig process configuration
type ProcessConfig struct {
	ID                 IDt
	Name               string
	MqttBrokerAddr     string
	MqttClientID       string
	MqttBrokerUsername string
	MqttBrokerPassword string
	InfluxAddr         string
	InfluxUsername     string
	InfluxPassword     string
	InfluxDB           string
	// DataPoints are saved in batches .
	// Batch is sent to DB once it reaches BatchMaxSize or SaveInterval .
	// depends on what comes first .
	BatchMaxSize int
	// Interval in miliseconds
	SaveInterval time.Duration
}

// GetNewID returns new ID for any struct which implements CommonInt
func GetNewID(dsi interface{}) IDt {
	var newID int64
	dsSlice := reflect.ValueOf(dsi)
	if dsSlice.Kind() == reflect.Slice {
		for i := 0; i < dsSlice.Len(); i++ {
			id := dsSlice.Index(i).FieldByName("ID").Int()
			if id > newID {
				newID = id
			}
		}

		newID++
		return IDt(newID)
	}
	return 0

}

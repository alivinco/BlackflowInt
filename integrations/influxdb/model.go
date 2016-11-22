package influxdb

import (
	"time"

	"reflect"

	"github.com/alivinco/iotmsglibgo"
	influx "github.com/influxdata/influxdb/client/v2"
)

// Transform defines function which converts IotMsg into influx data point
type Transform func(context *MsgContext, topic string, iotMsg *iotmsglibgo.IotMsg, domain string) (*influx.Point, error)

// IDt defines type of struct ID
type IDt int

// MsgContext describes metadata collected along message goint throughout the pipeline  .
type MsgContext struct {
	FilterID        IDt
	MeasurementName string
}

// Selector defines message selector.
type Selector struct {
	ID       IDt
	Topic    string
	InMemory bool
}

// Filter defines message filter.
// Emty string - means no filter.
type Filter struct {
	ID          IDt
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
	// Optional field , all tags defined here will be converted into influxDb tags
	Tags map[string]string
	// If set , then the value will overrride default measurement name defined in transformation
	MeasurementName string
	// definies if filter is temporary and can be stored in memory or persisted to disk
	InMemory bool
}

// Measurement stores measurement specific configs like retention policy
type Measurement struct {
	Name string
	// Normally should be combination of Name + RetentionPolicy
	RetentionPolicyName string
	// Have to be in format 1d , 1w , 1h .
	RetentionPolicyDuration string
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
	Filters      []Filter
	Selectors    []Selector
	Measurements []Measurement
}

//ProcessConfigs is a collection of ProcessConfigs
type ProcessConfigs []ProcessConfig

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

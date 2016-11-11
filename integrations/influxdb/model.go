package influxdb

import (
	"time"

	"github.com/alivinco/iotmsglibgo"
	influx "github.com/influxdata/influxdb/client/v2"
)

// Transform defines function which converts IotMsg into influx data point
type Transform func(topic string, iotMsg *iotmsglibgo.IotMsg, domain string) (*influx.Point, error)

// Selector defines message selector.
type Selector struct {
	Topic string
}

// Filter defines message filter.
// Emty string - means no filter.
type Filter struct {
	ID          int
	Name        string
	Topic       string
	Domain      string
	MsgType     string
	MsgClass    string
	MsgSubClass string
	//
	Negation bool
	// Boolean operation between 2 filters , the value can one of : and , or
	LinkedFilterBooleanOperation string
	LinkedFilterID               int
	IsAtomic                     bool
}

// Config process configuration
type Config struct {
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

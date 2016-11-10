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
type Filter struct {
	Topic       string
	Domain      string
	MsgType     string
	MsgClass    string
	MsgSubClass string
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

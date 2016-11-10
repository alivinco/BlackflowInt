package influxdb

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"sync"

	"fmt"

	"github.com/alivinco/blackflowint/adapters"
	iotmsg "github.com/alivinco/iotmsglibgo"
	influx "github.com/influxdata/influxdb/client/v2"
)

// Process implements integration flow between messaging system and influxdb timeseries database.
// It inserts events into db
type Process struct {
	mqttAdapter *adapters.MqttAdapter
	influxC     influx.Client
	selectors   []Selector
	filters     []Filter
	Config      *Config
	batchPoints influx.BatchPoints
	ticker      *time.Ticker
	writeMutex  *sync.Mutex
	transform   Transform
}

// NewProcess is a constructor
func NewProcess(config *Config, selectors []Selector, filters []Filter) *Process {
	proc := Process{Config: config, selectors: selectors, filters: filters, transform: DefaultTransform}
	return &proc
}

// Init doing the process bootrstrap .
func (pr *Process) Init() error {
	var err error
	log.Info("Initializing influx client.")
	pr.influxC, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     pr.Config.InfluxAddr, //"http://localhost:8086",
		Username: pr.Config.InfluxUsername,
		Password: pr.Config.InfluxPassword,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
		return err
	}
	// Creating database
	q := influx.NewQuery(fmt.Sprintf("CREATE DATABASE %s", pr.Config.InfluxDB), "", "")
	if response, err := pr.influxC.Query(q); err == nil && response.Error() == nil {
		log.Infof("Database %s was created with status :%s", pr.Config.InfluxDB, response.Results)
	}
	// Create a new point batch
	pr.batchPoints, err = influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  pr.Config.InfluxDB,
		Precision: "ns",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}
	log.Info("Initialization completed.")
	log.Info("Initializing MQTT adapter.")
	//"tcp://localhost:1883", "blackflowint", "", ""
	pr.mqttAdapter = adapters.NewMqttAdapter(pr.Config.MqttBrokerAddr, pr.Config.MqttClientID, pr.Config.MqttBrokerUsername, pr.Config.MqttBrokerPassword)
	pr.mqttAdapter.SetMessageHandler(pr.OnMessage)
	log.Info("Initialization completed.")

	pr.writeMutex = &sync.Mutex{}
	return nil
}

// OnMessage is invoked by an adapter on every new message
// The code is executed in callers goroutine
func (pr *Process) OnMessage(topic string, iotMsg *iotmsg.IotMsg, domain string) {
	log.Debugf("New msg of class = %s", iotMsg.Class)
	if pr.filter(topic, iotMsg) {
		msg, err := pr.transform(topic, iotMsg, domain)
		if err != nil {
			log.Errorf("Transformation error: %s", err)
		} else {
			if msg != nil {
				pr.write(msg)
			} else {
				log.Debug("Message can't be mapped .Skipping .")
			}

		}
	} else {
		log.Debugf("Message from topic %s is skiped .", topic)
	}
}

// Filter - transforms IotMsg into DB compatable struct
func (pr *Process) filter(topic string, iotMsg *iotmsg.IotMsg) bool {
	return true
}

func (pr *Process) write(point *influx.Point) {
	pr.writeMutex.Lock()
	pr.batchPoints.AddPoint(point)
	pr.writeMutex.Unlock()
	if len(pr.batchPoints.Points()) >= pr.Config.BatchMaxSize {
		pr.WriteIntoDb()
	}
}

// Configure - configures the process and subscribes for all topics
func (pr *Process) Configure(selectors []Selector, filters []Filter) {
	pr.selectors = selectors
	pr.filters = filters
}

// WriteIntoDb - inserts record into db
func (pr *Process) WriteIntoDb() {
	// Mutex is needed to fix condition when the function is invoked by timer and batch size almost at the same time
	defer func() {
		pr.writeMutex.Unlock()
	}()
	pr.writeMutex.Lock()
	if len(pr.batchPoints.Points()) == 0 {
		return
	}

	log.Debugf("Writing batch of size = %d", len(pr.batchPoints.Points()))
	var err error
	err = pr.influxC.Write(pr.batchPoints)
	if err != nil {
		log.Error("Error: ", err)
	}
	// Create a new point batch
	// Create a new point batch
	pr.batchPoints, err = influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  pr.Config.InfluxDB,
		Precision: "ns",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}
	// points = []*influx.Point{}

}

// Start starts the process by starting MQTT adapter ,
// starting scheduler
func (pr *Process) Start() error {
	err := pr.Init()
	if err != nil {
		log.Fatalln("Initialization Error: ", err)
		return err
	}
	pr.ticker = time.NewTicker(time.Millisecond * pr.Config.SaveInterval)
	go func() {
		for _ = range pr.ticker.C {
			pr.WriteIntoDb()
		}
	}()
	err = pr.mqttAdapter.Start()
	if err != nil {
		log.Fatalln("Error: ", err)
		return err
	}
	for _, selector := range pr.selectors {
		pr.mqttAdapter.Subscribe(selector.Topic, 0)
	}
	return nil

}

// Stop stops the process by unsubscribing from all topics ,
// stops scheduler and stops adapter.
func (pr *Process) Stop() {
	pr.ticker.Stop()

	for _, selector := range pr.selectors {
		pr.mqttAdapter.Unsubscribe(selector.Topic)
	}
	pr.influxC.Close()
	pr.mqttAdapter.Stop()

}

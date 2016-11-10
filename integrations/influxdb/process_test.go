package influxdb

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/adapters"
	iotmsg "github.com/alivinco/iotmsglibgo"
	influx "github.com/influxdata/influxdb/client/v2"
)

func Setup() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.InfoLevel)
}

func MsgGenerator(config Config, numberOfMsg int) {
	r := rand.New(rand.NewSource(99))
	topics := []string{
		"jim1/evt/zw/2/sen_temp/1",
		"jim1/evt/zw/3/sen_temp/1",
		"jim1/evt/zw/4/sen_temp/1",
	}
	config.MqttClientID = "blackflowint_pub_test"
	mqttAdapter := adapters.NewMqttAdapter(config.MqttBrokerAddr, config.MqttClientID, config.MqttBrokerUsername, config.MqttBrokerPassword)
	mqttAdapter.Start()
	for i := 0; i < numberOfMsg; i++ {
		msg := iotmsg.NewIotMsg(iotmsg.MsgTypeEvt, "sensor", "temperature", nil)
		msg.SetDefaultFloat(r.Float64(), "C")
		mqttAdapter.Publish(topics[r.Intn(len(topics))], msg, 0, "testDomain")
	}
	time.Sleep(time.Second * 3)
	mqttAdapter.Stop()

}

func CleanUpDB(influxC influx.Client, config *Config) {
	// Delete measurments
	q := influx.NewQuery(fmt.Sprintf("DELETE from sensors"), config.InfluxDB, "")
	if response, err := influxC.Query(q); err == nil && response.Error() == nil {
		log.Infof("Datebase was created with status :", response.Results)

	}

}

func Count(influxC influx.Client, config *Config) int {
	q := influx.NewQuery(fmt.Sprintf("select count(value) from sensors"), config.InfluxDB, "")
	if response, err := influxC.Query(q); err == nil && response.Error() == nil {
		countN, ok := response.Results[0].Series[0].Values[0][1].(json.Number)
		count, _ := countN.Int64()
		if !ok {
			log.Errorf("Type assertion failed , type is = %s", reflect.TypeOf(response.Results[0].Series[0].Values[0][1]))
		}
		log.Info("Number of received messages = ", count)
		return int(count)
	}
	return 0
}

func TestProcess(t *testing.T) {
	Setup()
	//Start container : docker run --name influxdb -d -p 8084:8083 -p 8086:8086 -v influxdb:/var/lib/influxdb influxdb:1.1.0-rc1-alpine
	//Start mqtt broker
	CountOfMessagesToSend := 10000
	config := Config{
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
	}
	selector := []Selector{
		Selector{Topic: "testDomain/jim1/evt/zw/2/sen_temp/1"},
		Selector{Topic: "testDomain/jim1/evt/zw/3/sen_temp/1"},
		Selector{Topic: "testDomain/jim1/evt/zw/4/sen_temp/1"},
		Selector{Topic: "testDomain/jim1/evt/zw/3/bin_switch/1"},
	}

	influxC, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     config.InfluxAddr, //"http://localhost:8086",
		Username: config.InfluxUsername,
		Password: config.InfluxPassword,
	})
	if err != nil {
		t.Fatal("Error: ", err)
	}

	CleanUpDB(influxC, &config)
	filters := []Filter{}
	proc := NewProcess(&config, selector, filters)
	err = proc.Start()
	if err != nil {
		t.Fatal(err)
	}
	go MsgGenerator(config, CountOfMessagesToSend)

	time.Sleep(time.Second * 5)
	CountOfSavedEvents := Count(influxC, &config)
	if CountOfMessagesToSend != CountOfSavedEvents {
		t.Errorf("Count of sent messages doesn't match count of saved messages. Number of sent messages = %d , number of saved events = %d", CountOfMessagesToSend, CountOfSavedEvents)
	}
	proc.Stop()
	influxC.Close()

}

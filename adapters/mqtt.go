package adapters

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/iotmsglibgo"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttAdapter , mqtt adapter .
type MqttAdapter struct {
	client     MQTT.Client
	msgHandler MessageHandler
}

type MessageHandler func(topic string, iotMsg *iotmsglibgo.IotMsg, domain string)

// NewMqttAdapter constructor
//serverUri="tcp://localhost:1883"
func NewMqttAdapter(serverURI string, clientID string, username string, password string) *MqttAdapter {
	mh := MqttAdapter{}
	opts := MQTT.NewClientOptions().AddBroker(serverURI)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(mh.onMessage)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(mh.onConnectionLost)
	opts.SetOnConnectHandler(mh.onConnect)
	//create and start a client using the above ClientOptions
	mh.client = MQTT.NewClient(opts)
	return &mh
}

// SetMessageHandler message handler setter
func (mh *MqttAdapter) SetMessageHandler(msgHandler MessageHandler) {
	mh.msgHandler = msgHandler
}

// Start , starts adapter async.
func (mh *MqttAdapter) Start() error {
	log.Info("Connecting to MQTT broker ")
	if token := mh.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Stop , stops adapter.
func (mh *MqttAdapter) Stop() {
	mh.client.Disconnect(250)
}

// Subscribe - subscribing for topic
func (mh *MqttAdapter) Subscribe(topic string, qos byte) error {
	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	log.Debug("Subscribing to topic:", topic)
	if token := mh.client.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return token.Error()
	}
	return nil
}

// Unsubscribe , unsubscribing from topic
func (mh *MqttAdapter) Unsubscribe(topic string) error {
	log.Debug("Unsubscribing from topic:", topic)
	if token := mh.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mh *MqttAdapter) onConnectionLost(client MQTT.Client, err error) {
	log.Errorf("Connection lost with MQTT broker . Error : %v", err)
}

func (mh *MqttAdapter) onConnect(client MQTT.Client) {
	log.Infof("Connection established with MQTT broker .")
}

//define a function for the default message handler
func (mh *MqttAdapter) onMessage(client MQTT.Client, msg MQTT.Message) {
	log.Debugf(" New msg from TOPIC: %s", msg.Topic())
	// log.Debug("MSG: %s\n", msg.Payload())
	domain, topic := DetachDomainFromTopic(msg.Topic())
	iotMsg, err := iotmsglibgo.ConvertBytesToIotMsg(topic, msg.Payload(), nil)
	if err == nil {
		mh.msgHandler(topic, iotMsg, domain)
	} else {
		log.Error(err)

	}
}

// Publish iotMsg
func (mh *MqttAdapter) Publish(topic string, iotMsg *iotmsglibgo.IotMsg, qos byte, domain string) error {
	bytm, err := iotmsglibgo.ConvertIotMsgToBytes(topic, iotMsg, nil)
	topic = AddDomainToTopic(domain, topic)
	if err == nil {
		log.Debug("Publishing msg to topic:", topic)
		mh.client.Publish(topic, qos, false, bytm)
		return nil
	}
	return err

}

// AddDomainToTopic , adds prefix to topic .
func AddDomainToTopic(domain string, topic string) string {
	// Check if topic is already prefixed with  "/" if yes then concat without adding "/"
	// 47 is code of "/"
	if topic[0] == 47 {
		return domain + topic
	}
	return domain + "/" + topic
}

// DetachDomainFromTopic detaches domain from topic
func DetachDomainFromTopic(topic string) (string, string) {
	spt := strings.Split(topic, "/")
	// spt[0] - domain
	var top string
	if strings.Contains(spt[1], "jim") {
		top = strings.Replace(topic, spt[0]+"/", "", 1)
	} else {
		top = strings.Replace(topic, spt[0], "", 1)
	}
	// returns domain , topic
	return spt[0], top

}

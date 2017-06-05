package restmqttproxy


// IDt defines type of struct ID
type IDt int

type UserSettings struct {
	Name string
	ActiveDomain string
}

// ProcessConfig process configuration
type ProcessConfig struct {
	ID                 IDt
	Name               string
	MqttBrokerAddr     string
	MqttClientID       string
	MqttBrokerUsername string
	MqttBrokerPassword string
	ProxyBindHost      string
	Autostart    bool
	UserSettings []UserSettings
	ClientSecret 	   string
	SkipDomain	   bool

}

type RunningRequest struct {
	RespTopic string
	RespMsgType string
	ResponseChan chan []byte
}

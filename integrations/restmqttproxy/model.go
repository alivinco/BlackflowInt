package restmqttproxy


// IDt defines type of struct ID
type IDt int


// ProcessConfig process configuration
type ProcessConfig struct {
	ID                 IDt
	Name               string
	MqttBrokerAddr     string
	MqttClientID       string
	MqttBrokerUsername string
	MqttBrokerPassword string
	ProxyBindHost      string
	ProxyBindPort      string
	Autostart    bool
}
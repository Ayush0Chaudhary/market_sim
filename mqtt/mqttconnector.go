package mqttclient


import mqtt "github.com/eclipse/paho.mqtt.golang"

type MqttConnector interface {
	Subscribe(string) mqtt.Token
	SubscribeMultipleInstrumentToken([]int) mqtt.Token
	PublishJson(string, any) mqtt.Token
	PublishString(string, string) mqtt.Token
}

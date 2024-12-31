package mqttclient

import (
	"encoding/json"
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Mqtt struct {
	Client      mqtt.Client
	Environment string
}

func (m *Mqtt) Subscribe(topic string) mqtt.Token {
	return m.Client.Subscribe(topic, 0, nil)
}

func (m *Mqtt) SubscribeMultipleInstrumentToken(tokens []int) mqtt.Token {
	fmt.Println("subscribing to multiple topics")

	filters := make(map[string]byte)
	for _, token := range tokens {
		topic := m.Environment + "/gokicker/" + strconv.FormatUint(uint64(token), 10)
		filters[topic] = 0
	}

	token := m.Client.SubscribeMultiple(filters, nil)

	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			fmt.Println("error while subscribing to topic, %v", token.Error())
		}
	}()
	return token
}

func (m *Mqtt) PublishJson(topic string, value any) mqtt.Token {

	b, err := json.Marshal(value)
	if err != nil {
		fmt.Println("error while marshalling value for publishing json, %v", err)
		return nil
	}

	token := m.Client.Publish(topic, 0, true, b)

	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			fmt.Println("error while publish, %v", token.Error())
		}
	}()
	return token
}

func (m *Mqtt) PublishString(topic string, value string) mqtt.Token {
	return m.Client.Publish(topic, 0, true, value)
}

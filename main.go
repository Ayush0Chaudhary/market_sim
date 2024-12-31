package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqttclient "market_simulator/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type LighthouseResponse struct {
	Net []struct {
		InstrumentToken int `json:"instrument_token"`
		// Other fields as necessary
	}
	Day []struct {
		// Structure of Day array...
	}
}

var (
	ticks             = make(chan string)
	LighthousePostion []string
	volatility        = 3
)

func main() {
	// Call the function
	Hello()
	Process()
}

func Process() {
	for tick := range ticks {
		fmt.Println(tick)
	}
}

func Hello() {
	// Print the message
	fmt.Println("Hello, World!")
	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.SetConnectRetry(true)
	mqttOpts.AddBroker("tcp://localhost:1883")
	mqttOpts.SetAutoReconnect(true)

	// mqttOpts.SetUsername(helper.MustGetEnv("mqtt_username"))
	// mqttOpts.SetPassword(helper.MustGetEnv("mqtt_password"))
	mqttOpts.SetDefaultPublishHandler(func(c mqtt.Client, msg mqtt.Message) {
		// loop over known topics and check if the message is from any of them
		// if not, log the message

		log.Println("message received from unknown topic : ", msg.Topic())
	})

	mqttOpts.OnConnect = func(c mqtt.Client) {
		log.Println("connected to MQTT input server")
		m := mqttclient.Mqtt{
			Client:      c,
			Environment: "dev",
		}

		m.Subscribe("dev" + "/greekservice/startservice")
		m.Subscribe("dev" + "/greekservice/individualvalues_send")
		// m.Subscribe("dev" + "/gokicker/+")

		m.PublishString("dev"+"/gokicker/1", "start")
		StartPublish(m)
	}

	mqttOpts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("connection lost to mqtt broker at time %v :%v", time.Now(), err)
	}

	mqttClient := mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return
	}
}

func GetMockPositionFromLighthouse() (LighthouseResponse, error) {
	data, err := os.ReadFile("/home/panda/pr/market_simulator/positions.json")
	if err != nil {
		log.Fatal("error reading local file:", err)
		return LighthouseResponse{}, err
	}

	var response LighthouseResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println("error parsing json:", err)
		return LighthouseResponse{}, err
	}

	return response, nil
}

func StartPublish(mqtt mqttclient.Mqtt) {
	positions, err := GetMockPositionFromLighthouse()

	if err != nil {
		log.Println("error getting mock position from lighthouse:", err)
		return
	} else {
		for _, position := range positions.Net {
			LighthousePostion = append(LighthousePostion, fmt.Sprintf("%d", position.InstrumentToken))
		}
	}
	GoKicker(mqtt)
}

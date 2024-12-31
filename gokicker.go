package main

import (
	"fmt"
	mqttClient "market_simulator/mqtt"
	"math/rand"
	"strconv"
	"time"
)

type TickData struct {
	Mode            string
	InstrumentToken uint32
	IsTradable      bool
	IsIndex         bool

	// Timestamp represents Exchange timestamp
	Timestamp          string
	LastTradeTime      string
	LastPrice          float64
	LastTradedQuantity uint32
	TotalBuyQuantity   uint32
	TotalSellQuantity  uint32
	VolumeTraded       uint32
	TotalBuy           uint32
	TotalSell          uint32
	AverageTradePrice  float64
	OI                 uint32
	OIDayHigh          uint32
	OIDayLow           uint32
	NetChange          float64

	Symbol string
}

func GoKicker(mqtt mqttClient.Mqtt) {
	
	for {
		// Generate a random value to publish

		// Publish the value
		skipped := 0
		for _, position := range LighthousePostion {
			if skipped < 5 && rand.Intn(2) == 0 {
				skipped++
				continue
			}
			go publishOnPosition(position, mqtt)
		}
		

		sleepDuration := time.Duration(float64(volatility) * float64(time.Second))
		time.Sleep(sleepDuration)
		fmt.Println(time.Now())
		// Calculate the sleep duration based on volatility
	}
}

func publishOnPosition(token string, mqtt mqttClient.Mqtt) {
	// Publish the position
	tokenInt, err := strconv.Atoi(token)
	if err != nil {
		// Handle error
		return
	}
	mocktick := TickData{
		Mode:               "full",
		InstrumentToken:    uint32(tokenInt), // Convert token to uint32
		IsTradable:         true,
		IsIndex:            false,
		Timestamp:          time.Now().Format(time.RFC3339),
		LastTradeTime:      time.Now().Format(time.RFC3339),
		LastPrice:          rand.Float64(),
		LastTradedQuantity: uint32(rand.Intn(100)),
		TotalBuyQuantity:   uint32(rand.Intn(100)),
		TotalSellQuantity:  uint32(rand.Intn(100)),
		VolumeTraded:       uint32(rand.Intn(100)),
		TotalBuy:           uint32(rand.Intn(100)),
		TotalSell:          uint32(rand.Intn(100)),
		AverageTradePrice:  rand.Float64(),
		OI:                 uint32(rand.Intn(100)),
		OIDayHigh:          uint32(rand.Intn(100)),
		OIDayLow:           uint32(rand.Intn(100)),
		NetChange:          rand.Float64(),
		Symbol:             "NIFTY",
	}

	mqtt.PublishJson("dev/gokicker/"+token, mocktick)
}

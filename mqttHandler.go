package main

import (
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pmcgn/weatherreceiver/shared"
)

var config *shared.Configuration
var opts *mqtt.ClientOptions
var client mqtt.Client

func Init(c *shared.Configuration) {
	config = c
	log.Println("Connect to Broker", config.BrokerUrl, "as", "WeatherReceiver")

	opts = mqtt.NewClientOptions().
		AddBroker(config.BrokerUrl).
		SetClientID("WeatherReceiver").
		SetWill(config.Prefix+"/bridge/connectionState", "offline", config.Qos, config.Retained).
		SetAutoReconnect(true).
		SetConnectRetryInterval(15)

	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sendToBroker("online", "bridge/connectionState")
}

func Disconnect() {
	log.Println("Disconnect from Broker")
	sendToBroker("offline", "bridge/connectionState")
	client.Disconnect(100)
}

func MqttPublisher(events chan shared.WeatherStationEvent) {
	for {
		event := <-events
		parseStationEventAndSendToMqtt(event)
	}
}

func parseStationEventAndSendToMqtt(event shared.WeatherStationEvent) {
	sendToBroker(time.Now().Format("2006-01-02T15:04:05"), "station/lastSeen")

	for key, values := range event.RequestFormData {
		sendToBroker(values[0], key)

		switch key {
		case "tempf":
			sendToBroker(convertFahrenheitToCelsius(values[0]), "tempc")
		case "tempinf":
			sendToBroker(convertFahrenheitToCelsius(values[0]), "tempinc")
		case "rainratein":
			sendToBroker(convertInchToMillimeter(values[0]), "rainratemm")
		case "eventrainin":
			sendToBroker(convertInchToMillimeter(values[0]), "eventrainmm")
		case "hourlyrainin":
			sendToBroker(convertInchToMillimeter(values[0]), "hourlyrainmm")
		case "weeklyrainin":
			sendToBroker(convertInchToMillimeter(values[0]), "weeklyrainmm")
		case "yearlyrainin":
			sendToBroker(convertInchToMillimeter(values[0]), "yearlyrainmm")
		case "dailyrainin":
			sendToBroker(convertInchToMillimeter(values[0]), "dailyrainmm")
		case "monthlyrainin":
			sendToBroker(convertInchToMillimeter(values[0]), "monthlyrainmm")
		case "totalrainin":
			sendToBroker(convertInchToMillimeter(values[0]), "totalrainmm")
		case "windgustmph":
			sendToBroker(convertMphToKph(values[0]), "windgustkmh")
		case "windspeedmph":
			sendToBroker(convertMphToKph(values[0]), "windspeedkmh")
		}
	}
}

func convertFahrenheitToCelsius(value string) string {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "Error"
	}
	return strconv.FormatFloat((f-32)*5/9, 'f', 1, 32)
}

func convertInchToMillimeter(value string) string {
	in, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "Error"
	}
	return strconv.FormatFloat(in*25.4, 'f', 1, 32)
}

func convertMphToKph(value string) string {
	mph, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "Error"
	}
	return strconv.FormatFloat(mph*1.60934, 'f', 1, 32)
}

func sendToBroker(value string, topic string) bool {
	start := time.Now()
	fullTopic := config.Prefix + "/" + topic
	token := client.Publish(fullTopic, config.Qos, config.Retained, value)
	ok := token.WaitTimeout(2 * time.Second)
	err := token.Error()
	duration := time.Since(start)

	if ok && err == nil {
		return true
	}

	log.Println("ERROR: topic:", fullTopic, "with value: ", value, "Wait", ok, "Error:", err, duration)
	return false
}

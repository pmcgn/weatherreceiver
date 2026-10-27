package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pmcgn/weatherreceiver/shared"
)

func readEnviromentVariable(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseBool(value string) bool {
	return strings.ToLower(value) == "true" ||
		strings.ToLower(value) == "on" ||
		strings.ToLower(value) == "1" ||
		strings.ToLower(value) == "enabled"
}

func GetConfig() *shared.Configuration {
	qos, err := strconv.ParseInt(readEnviromentVariable("MQTT_QOS", "1"), 10, 64)
	if err != nil {
		fmt.Printf("Invalid value for QOS found: %d, must be a number. Falling back to default of 1", qos)
		qos = 1
	}

	retained := parseBool(readEnviromentVariable("MQTT_RETAIN", "false"))

	return &shared.Configuration{
		BrokerUrl: readEnviromentVariable("MQTT_BROKER_URL", "127.0.0.1:1883"),
		Prefix:    readEnviromentVariable("MQTT_PREFIX", "weather"),
		Qos:       byte(qos),
		Retained:  retained,
		HttpPort:  readEnviromentVariable("HTTP_PORT", "8090"),
	}
}

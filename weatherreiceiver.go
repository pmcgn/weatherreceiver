package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pmcgn/weatherreceiver/shared"
)

func main() {
	config := GetConfig()
	events := make(chan shared.WeatherStationEvent, 10)

	log.Printf("Starting receiver at HTTP port %v", config.HttpPort)
	Init(config)
	go startServer(events)
	go MqttPublisher(events)

	stop := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		Disconnect()
		os.Exit(0)
	}()
	<-done
}

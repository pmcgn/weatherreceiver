package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmcgn/weatherreceiver/shared"
)

var weatherEvents chan shared.WeatherStationEvent

func requestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	event := shared.WeatherStationEvent{
		RequestFormData: r.PostForm,
	}

	select {
	case weatherEvents <- event:
	default:
		log.Println("Could not send Weather Update to MQTT Broker")
	}
}

func startServer(events chan shared.WeatherStationEvent) {
	weatherEvents = events

	router := mux.NewRouter()
	router.HandleFunc("/data/report/", requestHandler)
	http.Handle("/", router)

	http.ListenAndServe(":"+config.HttpPort, nil)
}

package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	loadConfig()
	initMetrics()

	for _, sensor := range config.Sensors {
		go recordMetrics(sensor)
	}

	http.Handle("/metrics", loggingMiddleware(promhttp.Handler()))
	log.Printf("Starting server at %s", config.Server.Addr)
	log.Fatal(http.ListenAndServe(config.Server.Addr, nil))
}

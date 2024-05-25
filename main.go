package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type EnvironmentData struct {
	Temperature float64 `json:"TempC"`
	Humidity    float64 `json:"Humidity"`
}

var (
	temperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "box_temperature_celsius",
		Help: "Current box temperature in celsius",
	})
	humidityGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "box_humidity_percentage",
		Help: "Current box humidity in percentage",
	})
	dataFilePath = getEnv("DATA_FILE_PATH", "data.json")
	serverAddr   = getEnv("SERVER_ADDR", ":8080")
)

func init() {
	// Register the gauge with Prometheus's default registry
	prometheus.MustRegister(temperatureGauge)
	prometheus.MustRegister(humidityGauge)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func readEnvironmentData() (*EnvironmentData, error) {
	fileInfo, err := os.Stat(dataFilePath)
	if err != nil {
		return nil, err
	}

	// Check if the file is older than 5 minutes
	if time.Since(fileInfo.ModTime()) > 5*time.Minute {
		return nil, fmt.Errorf("data file is older than 5 minutes")
	}

	file, err := os.Open(dataFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dataBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var data EnvironmentData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func recordMetrics() {
	for {
		data, err := readEnvironmentData()
		if err != nil {
			log.Printf("Error reading environment data: %v", err)
			// Set gauges to an invalid state or clear them
			temperatureGauge.Set(math.NaN())
			humidityGauge.Set(math.NaN())
			time.Sleep(10 * time.Second)
			continue
		}
		temperatureGauge.Set(data.Temperature)
		humidityGauge.Set(data.Humidity)
		time.Sleep(30 * time.Second) // adjust the interval as needed
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	go recordMetrics()

	http.Handle("/metrics", loggingMiddleware(promhttp.Handler()))
	log.Printf("Starting server at %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

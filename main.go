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
	temperatureGaugeSensor1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "box_temperature_celsius",
		Help: "Current box temperature in celsius",
	})
	humidityGaugeSensor1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "box_humidity_percentage",
		Help: "Current box humidity in percentage",
	})

	temperatureGaugeSensor2 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "room_temperature_celsius",
		Help: "Current room temperature in celsius",
	})
	humidityGaugeSensor2 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "room_humidity_percentage",
		Help: "Current room humidity in percentage",
	})
	dataFilePathSensor1 = getEnv("DATA_FILE_PATH1", "/var/lib/dht22/data_4.json")
	dataFilePathSensor2 = getEnv("DATA_FILE_PATH2", "/var/lib/dht22/data_22.json")
	serverAddr          = getEnv("SERVER_ADDR", ":8080")
)

func init() {
	// Register the gauges with Prometheus's default registry
	prometheus.MustRegister(temperatureGaugeSensor1)
	prometheus.MustRegister(humidityGaugeSensor1)
	prometheus.MustRegister(temperatureGaugeSensor2)
	prometheus.MustRegister(humidityGaugeSensor2)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func readEnvironmentData(filePath string) (*EnvironmentData, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Check if the file is older than 5 minutes
	if time.Since(fileInfo.ModTime()) > 5*time.Minute {
		return nil, fmt.Errorf("data file %s is older than 5 minutes", filePath)
	}

	file, err := os.Open(filePath)
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

func recordMetrics(sensorID int, dataFilePath string, temperatureGauge, humidityGauge prometheus.Gauge) {
	for {
		data, err := readEnvironmentData(dataFilePath)
		if err != nil {
			log.Printf("Error reading environment data from sensor %d: %v", sensorID, err)
			temperatureGauge.Set(math.NaN())
			humidityGauge.Set(math.NaN())
			time.Sleep(10 * time.Second)
			continue
		} else {
			temperatureGauge.Set(data.Temperature)
			humidityGauge.Set(data.Humidity)
		}
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
	go recordMetrics(1, dataFilePathSensor1, temperatureGaugeSensor1, humidityGaugeSensor1)
	go recordMetrics(2, dataFilePathSensor2, temperatureGaugeSensor2, humidityGaugeSensor2)

	http.Handle("/metrics", loggingMiddleware(promhttp.Handler()))
	log.Printf("Starting server at %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

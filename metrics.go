package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var gauges map[int]struct {
	temperatureGauge prometheus.Gauge
	humidityGauge    prometheus.Gauge
}

func initMetrics() {
	gauges = make(map[int]struct {
		temperatureGauge prometheus.Gauge
		humidityGauge    prometheus.Gauge
	})

	for _, sensor := range config.Sensors {
		tempGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: sensor.TemperatureGaugeName,
			Help: fmt.Sprintf("Current temperature in celsius for sensor %d", sensor.ID),
		})
		humGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: sensor.HumidityGaugeName,
			Help: fmt.Sprintf("Current humidity in percentage for sensor %d", sensor.ID),
		})

		prometheus.MustRegister(tempGauge)
		prometheus.MustRegister(humGauge)

		gauges[sensor.ID] = struct {
			temperatureGauge prometheus.Gauge
			humidityGauge    prometheus.Gauge
		}{
			temperatureGauge: tempGauge,
			humidityGauge:    humGauge,
		}
	}
}

func recordMetrics(sensorConfig SensorConfig) {
	queryInterval := time.Duration(sensorConfig.QueryInterval) * time.Second
	if queryInterval == 0 {
		queryInterval = 30 * time.Second // default query interval
	}

	errorInterval := time.Duration(sensorConfig.ErrorInterval) * time.Second
	if errorInterval == 0 {
		errorInterval = 10 * time.Second // default error interval
	}

	for {
		data, err := readEnvironmentData(sensorConfig.FilePath)
		if err != nil {
			log.Printf("Error reading environment data from sensor %d: %v", sensorConfig.ID, err)
			gauges[sensorConfig.ID].temperatureGauge.Set(math.NaN())
			gauges[sensorConfig.ID].humidityGauge.Set(math.NaN())
			time.Sleep(errorInterval)
			continue
		} else {
			log.Printf("Read environment data from sensor %d: %+v", sensorConfig.ID, data)
			gauges[sensorConfig.ID].temperatureGauge.Set(data.Temperature)
			gauges[sensorConfig.ID].humidityGauge.Set(data.Humidity)
		}
		time.Sleep(queryInterval)
	}
}

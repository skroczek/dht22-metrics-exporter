# Environment Data Monitoring with Prometheus

This project collects environment data from multiple sensors and exposes the metrics for Prometheus to scrape. The data is read from JSON files, and the configuration for the sensors is specified in a YAML file.

## Features

- Reads temperature and humidity data from JSON files.
- Exposes metrics in Prometheus format.
- Configurable sensor settings via YAML file.
- Supports custom query and error intervals for each sensor.

## Project Structure

- `main.go`: The entry point of the application. It initializes the configuration and metrics, starts the HTTP server, and launches goroutines for each sensor.
- `config.go`: Handles loading and parsing the YAML configuration file.
- `metrics.go`: Defines and registers Prometheus metrics. Contains the logic for recording metrics from sensors.
- `environment_data.go`: Provides functions to read environment data from JSON files.
- `middleware.go`: Implements HTTP middleware for logging requests.

## Configuration

The configuration is done via a YAML file (`config.yaml`). Here is an example configuration:

```yaml
sensors:
  - id: 1
    file_path: "/var/lib/dht22/data_4.json"
    temperature_gauge_name: "box_temperature_celsius"
    humidity_gauge_name: "box_humidity_percentage"
    query_interval: 30  # Optional, in seconds
    error_interval: 10  # Optional, in seconds
  - id: 2
    file_path: "/var/lib/dht22/data_22.json"
    temperature_gauge_name: "room_temperature_celsius"
    humidity_gauge_name: "room_humidity_percentage"
    query_interval: 30  # Optional, in seconds
    error_interval: 10  # Optional, in seconds
server:
  addr: ":8080"
```

## Building the Project

To build the project, navigate to the project directory and run:

```bash
CC=arm-linux-gnueabihf-gcc CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o environment_monitor
```

## Running the Application

To run the application, simply execute the built binary:

```bash
./environment_monitor
```

The application will start an HTTP server on the address specified in the configuration file and expose Prometheus metrics at `/metrics`.

## Setting Up as a Systemd Service

To run this application as a systemd service, create a unit file at `/etc/systemd/system/environment_monitor.service` with the following content:

```ini
[Unit]
Description=Environment Data Monitoring Service
After=network.target

[Service]
Environment=CONFIG_FILE_PATH=/path/to/your/config.yaml
ExecStart=/path/to/your/project/environment_monitor
WorkingDirectory=/path/to/your/project
Restart=always
User=nobody
Group=nogroup

[Install]
WantedBy=multi-user.target
```

Replace `/path/to/your/project` with the actual path to your project directory.

### Enable and Start the Service

```bash
sudo systemctl enable environment_monitor
sudo systemctl start environment_monitor
```

### Checking the Service Status

```bash
sudo systemctl status environment_monitor
```

## License

This project is licensed under the MIT License.
```

This README provides an overview of the project, instructions for configuring and building it, and details on setting it up as a systemd service. Adjust paths and user/group settings in the systemd unit file as needed for your environment.
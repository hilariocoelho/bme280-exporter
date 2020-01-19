# BME280-Exporter
Golang Prometheus exporter for the Bosch BME280 sensor

## Usage

In order to support multiple devices/sensors, the metrics path contains the device and the I2C address of the sensor like the following: `/metrics/{device}/{register}`.  

Since this container runs with a low privilege user, you should change the group of device before mounting it in the container to `500` GID.

```shell script
sudo chgrp 500 /dev/i2c-1
docker pull spawn2kill/bme280-exporter:1.0.0
docker run -it --device /dev/i2c-1 -p 8080:8080 spawn2kill/bme280-exporter:1.0.0

curl -X GET localhost:8080/metrics/i2c-1/0x76
```

### Docker-Compose

```yaml
version: "3"

services:
  bme280:
    image: spawn2kill/bme280-exporter:1.0.0
    expose:
     - 8080
    devices:
     - /dev/i2c-1
```

### Prometheus Configuration

```yaml
scrape_configs:
  - job_name: 'bme280'
    scrape_interval: 5s
    static_configs:
    - targets:
      - 'bme280:8080'
      labels:
        alias: 'Room #1'
    metrics_path: '/metrics/i2c-1/0x76'
```
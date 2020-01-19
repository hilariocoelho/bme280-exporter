package main

import (
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
)

var (
	tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bme280_temperature_celsius",
			Help: "Temperature in celsius degree",
		},
		[]string{"device", "register"},
	)
	pressGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bme280_pressure_hpa",
			Help: "Barometric pressure in hPa",
		},
		[]string{"device", "register"},
	)
	humidityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bme280_humidity",
			Help: "Humidity in percentage of relative humidity",
		},
		[]string{"device", "register"},
	)
)

func main() {
	r := mux.NewRouter()
	r.Use(bme280collector)
	r.Handle("/metrics/{device}/{register}", promhttp.Handler())

	if err := prometheus.Register(tempGauge); err != nil {
		panic(err)
	}
	if err := prometheus.Register(pressGauge); err != nil {
		panic(err)
	}
	if err := prometheus.Register(humidityGauge); err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}

func bme280collector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		device := vars["device"]
		registerAddressStr := vars["register"]
		registerAddress, err := strconv.ParseInt(registerAddressStr, 0, 32)
		if err != nil {
			log.Println(err.Error())
			return
		}

		d, err := i2c.Open(&i2c.Devfs{Dev: path.Join("/dev/", device)}, int(registerAddress))
		if err != nil {
			log.Println(err.Error())
			return
		}

		b := bme280.New(d)
		err = b.Init()

		defer d.Close()

		t, p, h, err := b.EnvData()
		if err != nil {
			log.Println(err.Error())
			return
		}

		tempGauge.WithLabelValues(device, registerAddressStr).Set(t)
		pressGauge.WithLabelValues(device, registerAddressStr).Set(p)
		humidityGauge.WithLabelValues(device, registerAddressStr).Set(h)

		next.ServeHTTP(w, r)
	})
}

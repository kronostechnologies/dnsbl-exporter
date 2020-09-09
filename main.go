package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Addresses []string
	Blacklist string
	ListCodes map[string]string `yaml:"listCodes"`
	Interval int
}

var config Config

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("ok\n"))
}

func main() {
	yamlFile, ye := ioutil.ReadFile("config.yaml")
	if ye != nil {
		panic(ye)
	}

	ue := yaml.Unmarshal(yamlFile, &config)
	if ue != nil {
		panic(ue)
	}

	if config.Interval == 0 {
		config.Interval = 3600
	}

	prometheus.MustRegister(newBlCollector())

    http.HandleFunc("/health", health)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on :2112")
	serveError := http.ListenAndServe(":2112", nil)
	if serveError != nil {
		panic(serveError)
	}
}

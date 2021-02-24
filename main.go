package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	chassisParser()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe("localhost:9111", nil)
}

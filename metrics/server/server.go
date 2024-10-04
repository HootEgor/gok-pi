package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Listen(ip, port string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	address := ip + ":" + port
	return http.ListenAndServe(address, mux)
}

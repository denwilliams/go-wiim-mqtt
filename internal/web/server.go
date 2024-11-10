package web

import (
	"fmt"
	"net/http"

	"github.com/denwilliams/go-wiim-mqtt/internal/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func CreateHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		logging.Info("%s /", r.Method)
		fmt.Fprintf(w, "OK")
	})
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

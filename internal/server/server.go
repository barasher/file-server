package server

// https://dev.to/jinagamvasubabu/how-to-post-multipart-form-data-in-go-using-mux-22kp

import (
	"fmt"
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const keyParam = "key"
const logKeyKey = "key"

func Run(prov provider.Provider) {
	r := mux.NewRouter()

	// get
	getRequestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "file_server_get_request_duration_seconds",
			Help:        "Histogram concerning get request durations (seconds)",
			Buckets:     []float64{.0025, .005, .01, .025, .05, .1},
			ConstLabels: prometheus.Labels{"method": "GET", "path": "/key/{key}"},
		},
		[]string{},
	)
	getHandler := promhttp.InstrumentHandlerDuration(getRequestDuration, handlerGet{provider: prov})
	r.HandleFunc(fmt.Sprintf("/key/{%v}", keyParam), getHandler).Methods("GET")

	// set
	setRequestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "file_server_set_request_duration_seconds",
			Help:        "Histogram concerning set request durations (seconds)",
			Buckets:     []float64{.0025, .005, .01, .025, .05, .1},
			ConstLabels: prometheus.Labels{"method": "POST", "path": "/key/{key}"},
		},
		[]string{},
	)
	setHandler := promhttp.InstrumentHandlerDuration(setRequestDuration, handlerSet{provider: prov})
	r.HandleFunc(fmt.Sprintf("/key/{%v}", keyParam), setHandler).Methods("POST")

	// metrics
	r.Handle("/metrics", promhttp.Handler())

	// TODO spécifier le port
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	log.Info().Msg("Server running...")
	srv.ListenAndServe();

	// TODO graceful stop
}

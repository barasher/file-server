package server

// https://dev.to/jinagamvasubabu/how-to-post-multipart-form-data-in-go-using-mux-22kp

import (
	"fmt"
	"github.com/barasher/file-server/internal"
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const keyParam = "key"
const logKeyKey = "key"

type Server struct {
	router *mux.Router
	prov   provider.Provider
}

func (s *Server) Run() {
	// TODO spécifier le port
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}
	log.Info().Msg("Server running...")
	srv.ListenAndServe();
	// TODO graceful stop
}

func (s *Server) Close() {
	s.prov.Close()
}

func NewServer(c internal.ServerConf) (Server, error) {
	var err error
	s := Server{}

	switch c.Type {
	case internal.S3ProviderID:
		log.Info().Msg("Provider: S3")
		if s.prov, err = provider.NewS3Provider(c.S3Conf); err != nil {
			return s, err
		}
	case internal.LocalProviderID:
		log.Info().Msg("Provider: local")
		if s.prov, err = provider.NewLocalProvider(c.LocalConf); err != nil {
			return s, err
		}
	default:
		return s, fmt.Errorf("unknown provider type (%v)", c.Type)
	}

	s.router = mux.NewRouter()
	getRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "file_server_get_request_duration_seconds",
			Help:        "Histogram concerning get request durations (seconds)",
			Buckets:     []float64{.0025, .005, .01, .025, .05, .1},
			ConstLabels: prometheus.Labels{"method": "GET", "path": "/key/{key}"},
		},
		[]string{},
	)
	prometheus.Unregister(getRequestDuration)
	getHandler := promhttp.InstrumentHandlerDuration(getRequestDuration, handlerGet{provider: s.prov})
	setRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "file_server_set_request_duration_seconds",
			Help:        "Histogram concerning set request durations (seconds)",
			Buckets:     []float64{.0025, .005, .01, .025, .05, .1},
			ConstLabels: prometheus.Labels{"method": "POST", "path": "/key/{key}"},
		},
		[]string{},
	)
	prometheus.Unregister(setRequestDuration)
	prometheus.MustRegister(getRequestDuration, setRequestDuration)
	setHandler := promhttp.InstrumentHandlerDuration(setRequestDuration, handlerSet{provider: s.prov})
	s.router.HandleFunc(fmt.Sprintf("/key/{%v}", keyParam), setHandler).Methods("POST")
	s.router.HandleFunc(fmt.Sprintf("/key/{%v}", keyParam), getHandler).Methods("GET")
	s.router.Handle("/metrics", promhttp.Handler())

	return s, nil
}

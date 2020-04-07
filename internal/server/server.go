package server

import (
	"errors"
	"fmt"
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

const getKeyParam = "key"

type handlerGetKey struct {
	provider provider.Provider
}

func (h handlerGetKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k := mux.Vars(r)[getKeyParam]
	reader, err := h.provider.Provide(k)
	if err != nil {
		if errors.Is(err, provider.ErrKeyNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "", 503)
		return
	}
	defer reader.Close()
	w.WriteHeader(200)
	io.Copy(w, reader)
}

func Run(prov provider.Provider) {
	r := mux.NewRouter()

	// getKey
	requestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "s3_http_server_request_duration_seconds",
			Help: "Histogram concerning request durations (seconds)",
			Buckets:[]float64{.0025, .005, .01, .025, .05, .1},
			ConstLabels:prometheus.Labels{"method":"GET", "path":"/key/{key}"},
		},
		[]string{},
	)
	h := handlerGetKey{provider:prov}
	getKeyHandler := promhttp.InstrumentHandlerDuration(requestDuration, h)
	r.HandleFunc(fmt.Sprintf("/key/{%v}", getKeyParam) ,	getKeyHandler).Methods("GET")

	// metrics
	r.Handle("/metrics", promhttp.Handler())

	// TODO spécifier le port
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: r,
	}
	log.Info().Msg("Server running...")
	srv.ListenAndServe();

	// TODO graceful stop
}

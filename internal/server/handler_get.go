package server

import (
	"errors"
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"github.com/rs/zerolog/log"
)

type handlerGet struct {
	provider provider.Provider
}

func (h handlerGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k := mux.Vars(r)[keyParam]
	log.Info().Msgf("key: %v / params: %v", k, mux.Vars(r))
	reader, err := h.provider.Get(k)
	if err != nil {
		switch {
		case errors.Is(err, provider.ErrKeyNotFound):
			http.NotFound(w, r)
			return
		case errors.Is(err, provider.ErrChroot):
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Chroot escalation detected"))
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer reader.Close()
	w.WriteHeader(http.StatusOK)
	io.Copy(w, reader)
}
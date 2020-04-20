package server

import (
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type handlerSet struct {
	provider provider.Provider
}

func (h handlerSet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k := mux.Vars(r)[keyParam]
	subLog := log.With().Str("key", k).Logger()
	subLog.Debug().Msgf("Set key %v...", k)

	file, _, err := r.FormFile("file")
	if err != nil {
		subLog.Error().Msgf("%v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := h.provider.Set(k, file); err != nil {
		subLog.Error().Msgf("%v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package server

import (
	"fmt"
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	var tcs = []struct {
		tcID   string
		preProv provider.Provider
		inKey string
		expStatus int
		expContent string
	}{
		{"nominal", buildFakeProv("file content", nil, nil),"file.txt", http.StatusOK, "file content"},
		{ "unknown", buildFakeProv("", provider.ErrKeyNotFound, nil),"unknown", http.StatusNotFound, ""},
		{ "error", buildFakeProv("", fmt.Errorf("error"), nil),"unknown", http.StatusInternalServerError, ""},
	}

	for _, tc := range tcs {
		t.Run(tc.tcID, func(t *testing.T) {
			handler := handlerGet{tc.preProv}
			path := fmt.Sprintf("/key/%s", tc.inKey)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Handle("/key/{key}", handler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expStatus, rr.Code)
			if tc.expStatus== http.StatusOK {
				b, err := ioutil.ReadAll(rr.Body)
				assert.Nil(t, err)
				assert.Equal(t, tc.expContent, string(b))
			}
		})
	}
}



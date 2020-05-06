package server

import (
	"bytes"
	"fmt"
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetHandler(t *testing.T) {
	var tcs = []struct {
		tcID       string
		preProv    provider.Provider
		inField string
		expStatus  int
	}{
		{"nominal", buildFakeProv("",nil, nil),  "file",http.StatusNoContent},
		{"error", buildFakeProv("", nil, fmt.Errorf("error")), "file", http.StatusInternalServerError},
		{"wrongField", buildFakeProv("",nil, nil),  "wrongField",http.StatusInternalServerError},
		{"chrootError", buildFakeProv("",nil, provider.ErrChroot),  "file",http.StatusBadRequest},
	}

	for _, tc := range tcs {
		t.Run(tc.tcID, func(t *testing.T) {
			handler := handlerSet{tc.preProv}
			path := fmt.Sprintf("/key/%s", "blabla.txt")

			body := new(bytes.Buffer)
			mpart := multipart.NewWriter(body)
			part, err := mpart.CreateFormFile(tc.inField, "blabla.txt")
			assert.Nil(t, err)
			_, err = part.Write([]byte("upload"))
			assert.Nil(t, err)
			assert.Nil(t, mpart.Close())
			contentType := fmt.Sprintf("multipart/form-data; boundary=%s", mpart.Boundary())

			req, err := http.NewRequest("POST", path, body)
			req.Header.Set("Content-type", contentType)
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Handle("/key/{key}", handler)
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.expStatus, rr.Code)
		})
	}
}

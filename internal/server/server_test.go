package server

import (
	"bytes"
	"fmt"
	"github.com/barasher/file-server/internal"
	"github.com/barasher/file-server/internal/provider"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

func TestNewServer_ProvType(t *testing.T) {
	localConf := internal.ServerConf{
		Type: "local",
		LocalConf: provider.LocalConf{
			Folder: "../../testdata/local",
		},
	}

	s3Conf := internal.ServerConf{
		Type: "s3",
		S3Conf: provider.S3Conf{
			AccessId: "a", AccessSecret: "b", Bucket: "c", URL: "d",
		},
	}

	unkConf := internal.ServerConf{
		Type: "unknown",
	}

	var tcs = []struct {
		tcID       string
		inConf     internal.ServerConf
		expSuccess bool
	}{
		{"localConf", localConf, true},
		{"s3Conf", s3Conf, true},
		{"unknown", unkConf, false},
	}

	for _, tc := range tcs {
		t.Run(tc.tcID, func(t *testing.T) {
			p, err := NewServer(tc.inConf)
			assert.Equal(t, err == nil, tc.expSuccess)
			t.Logf("%v", err)
			if err == nil {
				p.Close()
			}
		})
	}
}

func TestNewServer_Mux(t *testing.T) {
	c := internal.ServerConf{
		Type: "local",
		LocalConf: provider.LocalConf{
			Folder: "../../testdata/local",
		},
	}
	s, err := NewServer(c)
	assert.Nil(t, err)
	defer s.Close()

	// get
	req, err := http.NewRequest("GET", "/key/file.txt", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// set
	defer os.Remove("../../testdata/local/file2.txt")
	body := new(bytes.Buffer)
	mpart := multipart.NewWriter(body)
	part, err := mpart.CreateFormFile("file", "file2.txt")
	assert.Nil(t, err)
	_, err = part.Write([]byte("upload"))
	assert.Nil(t, err)
	assert.Nil(t, mpart.Close())
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", mpart.Boundary())
	req, err = http.NewRequest("POST", "/key/file2.txt", body)
	req.Header.Set("Content-type", contentType)
	assert.Nil(t, err)
	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// metrics
	req, err = http.NewRequest("GET", "/metrics", nil)
	assert.Nil(t, err)
	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	metrics := rr.Body.Bytes()
	re := regexp.MustCompile(`.*file_server_get_request_duration_seconds_bucket.*`)
	assert.True(t, re.Match(metrics))
	re = regexp.MustCompile(`.*file_server_set_request_duration_seconds_bucket.*`)
	assert.True(t, re.Match(metrics))
}

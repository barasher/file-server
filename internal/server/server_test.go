package server

import (
	"fmt"
	"github.com/barasher/file-server/internal/provider"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeProv struct {
	valGet string
	errGet error
	errSet error
}

func (f fakeProv) Get(key string) (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(f.valGet)), f.errGet
}

func (f fakeProv) Set(k string, v io.Reader)  error {
	return f.errSet
}

func (f fakeProv) Close() {

}

func buildFakeProv(valGet string, errGet error, errSet error) fakeProv{
	return fakeProv{
		valGet: valGet,
		errGet: errGet,
		errSet:errSet,
	}
}

func TestHandlerGeyKey(t *testing.T) {
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

	conf := provider.LocalConf{ 		Folder: "../../testdata/local"	}
	prov, err := provider.NewLocalProvider(conf)
	assert.Nil(t, err)
	defer prov.Close()



	for _, tc := range tcs {
		t.Run(tc.tcID, func(t *testing.T) {
			handler := handlerGetKey{tc.preProv}
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



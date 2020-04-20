package server

import (
	"io"
	"io/ioutil"
	"strings"
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

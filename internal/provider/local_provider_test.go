package provider

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestProvideUnknownKey(t *testing.T) {
	conf := LocalConf{Folder:"../../testdata/local"}
	prov, err := NewLocalProvider(conf)
	assert.Nil(t, err)
	defer prov.Close()
	_, err = prov.Provide("unknown")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrKeyNotFound))
}

func TestProvideNominal(t *testing.T) {
	conf := LocalConf{Folder:"../../testdata/local"}
	prov, err := NewLocalProvider(conf)
	assert.Nil(t, err)
	defer prov.Close()
	reader, err := prov.Provide("file.txt")
	assert.Nil(t, err)
	defer reader.Close()
	b, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, "file content", string(b))
}

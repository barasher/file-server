package provider

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func createProvider(t *testing.T, f string) Provider {
	conf := LocalConf{Folder: "../../testdata/local"}
	prov, err := NewLocalProvider(conf)
	assert.Nil(t, err)
	return prov
}

func TestGetUnknownKey(t *testing.T) {
	prov := createProvider(t, "../../testdata/local")

	_, err := prov.Get("unknown")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrKeyNotFound))
}

func TestGetNominal(t *testing.T) {
	prov := createProvider(t, "../../testdata/local")
	checkKeyValue(t, prov, "file.txt", "file content")
}

func TestSetNewKey(t *testing.T) {
	k := "blabla.txt"
	v := "blabla"
	defer func() {
		os.Remove(fmt.Sprintf("../../testdata/local/%v", k))
	}()
	prov := createProvider(t, "../../testdata/local")

	r := strings.NewReader(v)
	err := prov.Set(k, r)
	assert.Nil(t, err)

	checkKeyValue(t, prov, k, v)
}

func TestSetExistingKey(t *testing.T) {
	k := "blabla.txt"
	v := "blabla"
	v2 := "blublu"
	defer func() {
		os.Remove(fmt.Sprintf("../../testdata/local/%v", k))
	}()
	prov := createProvider(t, "../../testdata/local")

	r := strings.NewReader(v)
	assert.Nil(t, prov.Set(k, r))

	r = strings.NewReader(v2)
	assert.Nil(t, prov.Set(k, r))

	checkKeyValue(t, prov, k, v2)
}

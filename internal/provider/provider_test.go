package provider

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func checkKeyValue(t *testing.T, p Provider, k string, v string) {
	r, err := p.Get(k)
	assert.Nil(t, err)
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, v, string(b))
}
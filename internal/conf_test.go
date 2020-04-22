package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfFile(t *testing.T) {
	var tcs = []struct {
		tcID   string
		inFile  string
		expSuccess bool
	}{
		{"Unknown", "unkknown.json", false},
		{"NonParsable", "../testdata/conf/unparsable.json", false},
		{"Nominal", "../testdata/conf/nominal.json", true},
	}

	for _, tc := range tcs {
		t.Run(tc.tcID, func(t *testing.T) {
			_ , err := LoadConfFile(tc.inFile)
			assert.Equal(t, tc.expSuccess, err == nil)
		})
	}
}
package geoip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	err := loadCityDB("./GeoLite2-City.mmdb")
	assert.NoError(t, err)
	m := &MaxMindReader{}
	result, err := m.Lookup("45.11.3.101", "en")
	assert.NoError(t, err)
	t.Log(result)
}

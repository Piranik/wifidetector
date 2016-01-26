package probecollector

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUniqueFunction(t *testing.T) {
	assert := assert.New(t)

	uniqueMac, err := net.ParseMAC("ca:cf:1a:8b:9b:d5")
	assert.Nil(err)
	assert.False(isUnique(uniqueMac))

	nonuniqueMac, err := net.ParseMAC("00:00:00:ea:95:ee")
	assert.Nil(err)
	assert.True(isUnique(nonuniqueMac))

}

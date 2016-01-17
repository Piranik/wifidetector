package store

import (
	"net"
	"testing"
	"time"

	"github.com/dereulenspiegel/wifidetector/probecollector"
	"github.com/stretchr/testify/assert"
)

func TestExpirationOfProbeRequests(t *testing.T) {
	assert := assert.New(t)

	pr1 := probecollector.ProbeRequest{
		HWAddr:    net.HardwareAddr{0x12, 0xAB, 0x05, 0xFF, 0x42, 0x23},
		Timestamp: time.Now(),
	}

	pr2 := probecollector.ProbeRequest{
		HWAddr:    net.HardwareAddr{0x12, 0xAB, 0x05, 0xFF, 0x42, 0x21},
		Timestamp: time.Now().Add(time.Minute * -2),
	}

	store := NewMemoryStore()
	store.PutProbeRequest(pr1)
	store.PutProbeRequest(pr2)
	assert.Equal(2, len(store.requestMap))

	store.ExpireOlderThan(time.Minute * 1)
	assert.Equal(1, len(store.requestMap))

}

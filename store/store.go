package store

import (
	"net"
	"time"

	"github.com/dereulenspiegel/wifidetector/probecollector"
)

func ConvertHWAddr(hwaddr net.HardwareAddr) [6]byte {
	out := new([6]byte)
	for i, _ := range out {
		out[i] = hwaddr[i]
	}
	return *out
}

type DataStore interface {
	PutProbeRequest(proberequest probecollector.ProbeRequest) bool
	ExpireOlderThan(span time.Duration) []probecollector.ProbeRequest
	FindLastProbeRequest(hwaddr net.HardwareAddr) *probecollector.ProbeRequest
	GetAllProbeRequests() []probecollector.ProbeRequest
}

type MemoryStore struct {
	requestMap map[[6]byte]probecollector.ProbeRequest
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		requestMap: make(map[[6]byte]probecollector.ProbeRequest),
	}
}

func (m *MemoryStore) PutProbeRequest(proberequest probecollector.ProbeRequest) bool {
	lastPR := m.FindLastProbeRequest(proberequest.HWAddr)
	m.requestMap[ConvertHWAddr(proberequest.HWAddr)] = proberequest
	return lastPR == nil
}

func (m *MemoryStore) ExpireOlderThan(span time.Duration) []probecollector.ProbeRequest {
	now := time.Now()
	expiredPRs := make([]probecollector.ProbeRequest, 0, 100)
	for key, pr := range m.requestMap {
		if now.Sub(pr.Timestamp) > span {
			expiredPRs = append(expiredPRs, pr)
			delete(m.requestMap, key)
		}
	}
	return expiredPRs
}

func (m *MemoryStore) FindLastProbeRequest(hwaddr net.HardwareAddr) *probecollector.ProbeRequest {
	searchKey := ConvertHWAddr(hwaddr)
	pr, exists := m.requestMap[searchKey]
	if !exists {
		return nil
	}
	return &pr
}

func (m *MemoryStore) GetAllProbeRequests() []probecollector.ProbeRequest {
	proberequests := make([]probecollector.ProbeRequest, 0, len(m.requestMap))

	for _, pr := range m.requestMap {
		proberequests = append(proberequests, pr)
	}
	return proberequests
}

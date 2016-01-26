package store

import (
	"net"
	"sync"
	"time"

	"github.com/dereulenspiegel/wifidetector/probecollector"
)

func ConvertHWAddr(hwaddr probecollector.HardwareAddr) [6]byte {
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
	lock       *sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		requestMap: make(map[[6]byte]probecollector.ProbeRequest),
		lock:       &sync.Mutex{},
	}
}

func (m *MemoryStore) PutProbeRequest(proberequest probecollector.ProbeRequest) bool {
	lastPR := m.FindLastProbeRequest(net.HardwareAddr(proberequest.HWAddr))
	if lastPR != nil {
		proberequest.Firstseen = lastPR.Firstseen
	} else {
		proberequest.Firstseen = time.Now()
	}
	m.lock.Lock()
	m.requestMap[ConvertHWAddr(proberequest.HWAddr)] = proberequest
	m.lock.Unlock()
	return lastPR == nil
}

func (m *MemoryStore) ExpireOlderThan(span time.Duration) []probecollector.ProbeRequest {
	now := time.Now()
	expiredPRs := make([]probecollector.ProbeRequest, 0, 100)
	m.lock.Lock()
	for key, pr := range m.requestMap {
		if now.Sub(pr.Timestamp) > span {
			expiredPRs = append(expiredPRs, pr)
			delete(m.requestMap, key)
		}
	}
	m.lock.Unlock()
	return expiredPRs
}

func (m *MemoryStore) FindLastProbeRequest(hwaddr net.HardwareAddr) *probecollector.ProbeRequest {
	searchKey := ConvertHWAddr(probecollector.HardwareAddr(hwaddr))
	m.lock.Lock()
	pr, exists := m.requestMap[searchKey]
	m.lock.Unlock()
	if !exists {
		return nil
	}
	return &pr
}

func (m *MemoryStore) GetAllProbeRequests() []probecollector.ProbeRequest {
	m.lock.Lock()
	proberequests := make([]probecollector.ProbeRequest, 0, len(m.requestMap))
	m.lock.Unlock()

	for _, pr := range m.requestMap {
		proberequests = append(proberequests, pr)
	}
	return proberequests
}

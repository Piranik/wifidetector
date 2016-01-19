package probecollector

/*
#include <stdint.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

import "net"
import "time"

type HardwareAddr net.HardwareAddr

type ProbeRequest struct {
	HWAddr         HardwareAddr
	SignalStrength int
	Timestamp      time.Time
}

func (h HardwareAddr) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, h.String())), nil
}

const hexDigit = "0123456789abcdef"

func (h HardwareAddr) String() string {
	if len(h) == 0 {
		return ""
	}
	buf := make([]byte, 0, len(h)*3-1)
	for i, b := range h {
		if i > 0 {
			buf = append(buf, ':')
		}
		buf = append(buf, hexDigit[b>>4])
		buf = append(buf, hexDigit[b&0xF])
	}
	return string(buf)
}

//export pb_callback
func pb_callback(mac *C.uint8_t, signal_strength C.int) {
	macSlice := C.GoBytes(unsafe.Pointer(mac), 6)
	//defer C.free(mac)
	hwAddr := net.HardwareAddr(macSlice)
	prequest := ProbeRequest{
		HWAddr:         HardwareAddr(hwAddr),
		SignalStrength: int(signal_strength),
		Timestamp:      time.Now(),
	}
	resultChannel <- prequest
}

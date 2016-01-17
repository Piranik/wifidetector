package probecollector

/*
#include <stdint.h>
*/
import "C"
import "unsafe"

import "net"
import "time"

type ProbeRequest struct {
	HWAddr         net.HardwareAddr
	SignalStrength int
	Timestamp      time.Time
}

//export pb_callback
func pb_callback(mac *C.uint8_t, signal_strength C.int) {
	macSlice := C.GoBytes(unsafe.Pointer(mac), 6)
	//defer C.free(mac)
	hwAddr := net.HardwareAddr(macSlice)
	prequest := ProbeRequest{
		HWAddr:         hwAddr,
		SignalStrength: int(signal_strength),
		Timestamp:      time.Now(),
	}
	resultChannel <- prequest
}

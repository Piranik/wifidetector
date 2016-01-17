package probecollector

/*
#include <stdint.h>
*/
import "C"
import "log"
import "net"
import "time"

type ProbeRequest struct {
	HWAddr         net.HardwareAddr
	SignalStrength int
	Timestamp      time.Time
}

//export pb_callback
func pb_callback(mac *C.uint8_t, signal_strength C.int) {
	log.Printf("Received mac %v with strength %d", mac, signal_strength)
}

package probecollector

/*
#cgo linux LDFLAGS: -lpcap
#cgo freebsd LDFLAGS: -lpcap
#cgo darwin LDFLAGS: -lpcap
#cgo windows CFLAGS: -I C:/WpdPack/Include
#cgo windows,386 LDFLAGS: -L C:/WpdPack/Lib -lwpcap
#cgo windows,amd64 LDFLAGS: -L C:/WpdPack/Lib/x64 -lwpcap

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <unistd.h>
#include <math.h>
#include <string.h>

#include <pcap/pcap.h>

#define BUFFSIZE 3000
#define PROMISCUOUS 1
#define IF_MACADDR   6

pcap_t *pcd = NULL;
struct bpf_program bpg;

typedef struct {
    uint8_t radiotap_stuff[34];

    uint8_t rssi;
    uint8_t antenna;
    //uint16_t rx_flags;

    uint16_t flags;
    uint16_t duration;
    uint8_t destination[IF_MACADDR];
    uint8_t source[IF_MACADDR];
    uint8_t bssid[IF_MACADDR];

    uint16_t sequence_control;
} prequest_t;

void pb_callback(uint8_t*, int);

void packet_view(unsigned char *arg, const struct pcap_pkthdr *h,
                 const unsigned char *p) {
    prequest_t *pr = (prequest_t *) p;
    int signal_strength = fabs((pr->rssi / 255.0 * 100) / 2 - 100);

    pb_callback(pr->source, signal_strength);
}

static int configure_pcap(char *dev) {
    char errbuf[PCAP_ERRBUF_SIZE];
    memset(errbuf, 0, PCAP_ERRBUF_SIZE);
    pcd = pcap_open_live(dev, BUFFSIZE, PROMISCUOUS, 1000, errbuf); // wait for 1 sec (1000 ms)

    if (!pcd) {
        return 3;
    }

    memset(&bpg, 0, sizeof(bpg));

    if (0 != pcap_compile(pcd, &bpg, "wlan type mgt subtype probe-req", 1,
                          PCAP_NETMASK_UNKNOWN)) {
      return 4;
    }

    if (0 != pcap_setfilter(pcd, &bpg)) {
        return 5;
    }
    return 0;
}

static int start_pcab() {
  int count = 0;
  if (pcap_loop(pcd, -1, packet_view, (u_char *) &count) == -1) {
    return 6;
  }
  pcap_freecode(&bpg);
  pcap_close(pcd);
  return 0;
}
*/
import "C"
import (
	"fmt"
	"os/exec"
	"unsafe"
)

var (
	resultChannel chan ProbeRequest
)

func init() {
	resultChannel = make(chan ProbeRequest, 100)
}

func SetupInterface(phyIface, interfaceName string) error {
	// sudo iw phy phy0 interface add mon0 type monitor
	// TODO test if this is workin as expected
	cmd := exec.Command("sudo", "iw", "phy", phyIface, "interface", "add", interfaceName, "type", "monitor")
	_, err := cmd.CombinedOutput()
	return err
}

func StartCollection(interfaceName string) (chan ProbeRequest, error) {
	var cname *C.char = C.CString(interfaceName)
	defer C.free(unsafe.Pointer(cname))
	rc := C.configure_pcap(cname)
	if rc != 0 {
		switch rc {
		case 3:
			return nil, fmt.Errorf("Can't open device %s", interfaceName)
		case 4:
			return nil, fmt.Errorf("Can't compile filter statement")
		case 5:
			return nil, fmt.Errorf("Can't set filter")
		}
	}
	go func() {
		C.start_pcab()
	}()

	return resultChannel, nil
}

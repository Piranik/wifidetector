package probecollector

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
)

var (
	pcapHandle          *pcap.Handle
	ProbeRequestChannel chan ProbeRequest
	packetSource        *gopacket.PacketSource
)

func Init(iface string) error {
	ProbeRequestChannel = make(chan ProbeRequest, 100)

	if handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever); err != nil {
		return err
	} else {
		pcapHandle = handle
	}

	if err := pcapHandle.SetBPFFilter("wlan type mgt subtype probe-req"); err != nil {
		return err
	}
	packetSource = gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())

	return nil
}

func Read() {
	for packet := range packetSource.Packets() {
		log.Println(packet.Dump())
	}
}

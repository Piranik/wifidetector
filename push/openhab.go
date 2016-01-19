package push

import (
	"net"

	"github.com/dereulenspiegel/openhab-cli/openhab"
	"github.com/dereulenspiegel/wifidetector/probecollector"
	"github.com/dereulenspiegel/wifidetector/store"
)

type OpenHABPusher struct {
	client        *openhab.Client
	monitoredMACs map[[6]byte]string
}

func NewOpenHABPusher(host string) *OpenHABPusher {
	openHABClient := openhab.NewClient(host)

	return &OpenHABPusher{
		client:        openHABClient,
		monitoredMACs: make(map[[6]byte]string),
	}
}

func (o *OpenHABPusher) AddMonitoredMAC(macAddress, itemName string) error {
	hwAddr, err := net.ParseMAC(macAddress)
	if err != nil {
		return err
	}
	key := store.ConvertHWAddr(probecollector.HardwareAddr(hwAddr))
	o.monitoredMACs[key] = itemName
	return nil
}

func (o *OpenHABPusher) DeviceFound(pr probecollector.ProbeRequest) {
	key := store.ConvertHWAddr(pr.HWAddr)
	if itemName, exists := o.monitoredMACs[key]; exists {
		o.client.SendCommand(itemName, "ON")
	}
}

func (o *OpenHABPusher) DeviceLost(pr probecollector.ProbeRequest) {
	key := store.ConvertHWAddr(pr.HWAddr)
	if itemName, exists := o.monitoredMACs[key]; exists {
		o.client.SendCommand(itemName, "OFF")
	}
}

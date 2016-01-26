package main

import (
	"flag"
	"log"
	"time"

	"github.com/dereulenspiegel/wifidetector/config"
	"github.com/dereulenspiegel/wifidetector/probecollector"
	"github.com/dereulenspiegel/wifidetector/push"
	"github.com/dereulenspiegel/wifidetector/rest"
	"github.com/dereulenspiegel/wifidetector/store"
)

var (
	configFile = flag.String("config", "", "Configuration file")
	db         store.DataStore
)

func main() {
	flag.Parse()
	config.ParseConfig(*configFile)

	monitorDevice := config.GlobalConfig.MonitorDevice
	if monitorDevice == "" {
		monitorDevice = "mon0"
	}

	db = store.NewMemoryStore()
	go rest.InitRestAPI(db)
	pusher := push.NewOpenHABPusher(config.GlobalConfig.OpenHABHost)
	for mac, item := range config.GlobalConfig.MonitoredMACs {
		pusher.AddMonitoredMAC(mac, item)
	}
	resultChan, err := probecollector.StartCollection(monitorDevice)
	if err != nil {
		log.Fatalf("Can't initialise probe collection: %v", err)
	}

	for pr := range resultChan {
		if db.PutProbeRequest(pr) {
			pusher.DeviceFound(pr)
		}
		expiredPRs := db.ExpireOlderThan(time.Minute * time.Duration(config.GlobalConfig.ExpireAfter))
		for _, expiredPR := range expiredPRs {
			pusher.DeviceLost(expiredPR)
		}
	}
}

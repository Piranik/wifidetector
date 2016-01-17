package main

import (
	"flag"
	"log"
	"time"

	"github.com/dereulenspiegel/wifidetector/config"
	"github.com/dereulenspiegel/wifidetector/probecollector"
	"github.com/dereulenspiegel/wifidetector/push"
	"github.com/dereulenspiegel/wifidetector/store"
)

var (
	configFile = flag.String("config", "", "Configuration file")
	db         store.DataStore
)

func main() {
	db = store.NewMemoryStore()
	flag.Parse()
	config.ParseConfig(*configFile)
	pusher := push.NewOpenHABPusher(config.GlobalConfig.OpenHABHost)
	for mac, item := range config.GlobalConfig.MonitoredMACs {
		pusher.AddMonitoredMAC(mac, item)
	}
	resultChan, err := probecollector.StartCollection(config.GlobalConfig.MonitorDevice)
	if err != nil {
		log.Fatalf("Can't initialise probe collection: %v", err)
	}

	for pr := range resultChan {
		if db.PutProbeRequest(pr) {
			pusher.DeviceFound(pr)
		}
		expiredPRs := db.ExpireOlderThan(time.Minute * 1)
		for _, expiredPR := range expiredPRs {
			pusher.DeviceLost(expiredPR)
		}
	}
}

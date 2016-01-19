package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dereulenspiegel/wifidetector/store"
)

func mustEncode(data interface{}) string {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(dataBytes)
}

func GetAllProbeRequests(db store.DataStore) (int, string) {
	probes := db.GetAllProbeRequests()
	return http.StatusOK, mustEncode(&probes)
}

func GetProbeRequestCount(db store.DataStore) (int, string) {
	probes := db.GetAllProbeRequests()
	return http.StatusOK, fmt.Sprintf("%d", len(probes))
}

package rest

import (
	"github.com/dereulenspiegel/wifidetector/store"
	"github.com/go-martini/martini"
)

var (
	m      *martini.Martini
	router martini.Router

	db store.DataStore
)

func configureMartini() {
	m = martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.MapTo(db, (*store.DataStore)(nil))
}

func createRoutes() {
	router = martini.NewRouter()
	router.Get("/proberequests", GetAllProbeRequests)
	router.Get("/count", GetProbeRequestCount)
	m.MapTo(router, (*martini.Routes)(nil))
	m.Action(router.Handle)
}

func InitRestAPI(datastore store.DataStore) {
	db = datastore
	configureMartini()
	createRoutes()
	m.Run()
}

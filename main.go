package main

import (
	"log"
	"net/http"

	"github.com/UrlShortener/src/pkg/api_router"
	"github.com/UrlShortener/src/pkg/config"
	"github.com/UrlShortener/src/pkg/utility"
	"github.com/gorilla/mux"
)

func main() {
	initialize()
	router := mux.NewRouter()
	webApp(router)
	log.Println("Server started")
	port := config.AppConfig.ServerPort
	http.ListenAndServe(port, router)
}

func webApp(router *mux.Router) {
	router.HandleFunc("/url/shorten", api_router.CreateShortUrl).Methods(http.MethodPost)
	router.HandleFunc("/{hash}", api_router.FetchRedirecUrl).Methods(http.MethodGet)
}

func initialize() {
	config.Initialize()
	utility.ConnectZookeeper()
	utility.CreateRangeNode("/range")
	if utility.NodeRange == nil {
		utility.NodeRange = &utility.SequenceRange{
			Start: 0,
			Curr:  0,
			End:   0,
		}
	}
}

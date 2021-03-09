package main

import (
	"github.com/gorilla/mux"
	"harvanir/terraform-cloud-webhook-demo/pkg/config"
)

func main() {
	appConfig := config.Load()
	router := mux.NewRouter()
	registerHandler(router, appConfig)
	startServer(router, appConfig.Server.Port)
}

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"harvanir/terraform-cloud-webhook-demo/pkg/config"
	"net/http"
)

func main() {
	appConfig := config.Load()
	router := mux.NewRouter()
	registerHandler(router)
	startServer(router, appConfig.Server.Port)
}

func startServer(router *mux.Router, port int) {
	address := fmt.Sprintf("127.0.0.1:%v", port)
	logrus.Info("starting server on address: ", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		logrus.Error("failed to cmd-run the server", err)
	}
}

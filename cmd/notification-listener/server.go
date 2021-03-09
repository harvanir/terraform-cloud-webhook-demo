package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func startServer(router *mux.Router, port int) {
	address := fmt.Sprintf("127.0.0.1:%v", port)
	logrus.Info("starting server on address: ", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		logrus.Error("failed to cmd-run the server", err)
	}
}

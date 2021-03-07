package main

import (
	"github.com/gorilla/mux"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest-api/notification"
	"net/http"
)

func registerHandler(router *mux.Router) {
	// register handler
	router.HandleFunc("/terraform-notification", notification.Handler).Methods(http.MethodPost)
}

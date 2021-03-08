package main

import (
	"github.com/gorilla/mux"
	"harvanir/terraform-cloud-webhook-demo/pkg/config"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_handler"
)

func registerHandler(router *mux.Router, config *config.AppConfig) {
	handler := NewHandlerAdapter(config)
	rest_handler.RegisterNotification(router, handler)
}

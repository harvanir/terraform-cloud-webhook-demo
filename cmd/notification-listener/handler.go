package main

import (
	"github.com/gorilla/mux"
	"harvanir/terraform-cloud-webhook-demo/pkg/config"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_handler"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_handler_adapter"
)

func registerHandler(router *mux.Router, config *config.AppConfig) {
	handler := rest_handler_adapter.NewHandlerAdapter(config)
	rest_handler.RegisterNotification(router, handler)
}

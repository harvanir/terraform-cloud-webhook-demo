package main

import (
	"github.com/gorilla/mux"
	"harvanir/terraform-cloud-webhook-demo/pkg/config"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_api/notification"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_handler"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_handler_adapter"
)

func registerHandler(router *mux.Router, appConfig *config.AppConfig) {
	handlerAdapter := rest_handler_adapter.NewHandlerAdapter()
	handlerAdapter.NotificationHandler = notification.NewNotificationCtx(appConfig).Notification
	rest_handler.RegisterNotification(router, handlerAdapter)
}

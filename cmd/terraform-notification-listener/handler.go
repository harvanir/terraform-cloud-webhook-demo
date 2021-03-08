package main

import (
	"github.com/sirupsen/logrus"
	"harvanir/terraform-cloud-webhook-demo/pkg/config"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest-api/notification"
	"net/http"
)

type HandlerAdapter struct {
	NotificationHandler func(rw http.ResponseWriter, r *http.Request)
}

func NewHandlerAdapter(config *config.AppConfig) *HandlerAdapter {
	notificationCtx := notification.NewNotificationCtx(config)
	return &HandlerAdapter{
		NotificationHandler: notificationCtx.Notification,
	}
}

func (handler *HandlerAdapter) Notification(rw http.ResponseWriter, r *http.Request) {
	if handler.NotificationHandler != nil {
		handler.NotificationHandler(rw, r)
		return
	}
	logrus.Warn("no Notification handler registered")
}

package rest_handler_adapter

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type HandlerAdapter struct {
	NotificationHandler func(rw http.ResponseWriter, r *http.Request)
}

func NewHandlerAdapter() *HandlerAdapter {
	return &HandlerAdapter{}
}

func (handler *HandlerAdapter) Notification(rw http.ResponseWriter, r *http.Request) {
	if handler.NotificationHandler != nil {
		handler.NotificationHandler(rw, r)
		return
	}
	logrus.Warn("no Notification handler registered")
}

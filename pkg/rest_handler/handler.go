package rest_handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	Notification(rw http.ResponseWriter, r *http.Request)
}

// RegisterNotification register notification handler
func RegisterNotification(router *mux.Router, handler Handler) {
	router.HandleFunc("/terraform-notification", handler.Notification).Methods(http.MethodPost)
}

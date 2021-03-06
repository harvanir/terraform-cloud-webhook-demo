package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func setJsonHeader(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
}

// WriteJsonResponse write http response
func WriteJsonResponse(bytes []byte, rw http.ResponseWriter) {
	setJsonHeader(rw)
	_, err := rw.Write(bytes)
	if err != nil {
		logrus.Error("write response error", err)
	}
}

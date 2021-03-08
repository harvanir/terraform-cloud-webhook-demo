package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func setJsonHeader(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
}

// WriteResponse write http response
func WriteResponse(bytes []byte, rw http.ResponseWriter) {
	setJsonHeader(rw)
	_, err := rw.Write(bytes)
	if err != nil {
		logrus.Error("write response error", err)
	}
}

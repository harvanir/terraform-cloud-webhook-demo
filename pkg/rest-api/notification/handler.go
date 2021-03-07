package notification

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"harvanir/terraform-cloud-webhook-demo/pkg/json"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest-api/common"
	"net/http"
)

// Handler handle notification API
func Handler(rw http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := json.Decode(&payload, r); err != nil {
		logrus.Error(fmt.Errorf("failed decode to type: %w", err))
		common.WriteResponse([]byte("{\"message\":\"error\"}"), rw)
		return
	}
	bytes, err := json.Marshal(&payload)
	if err != nil {
		logrus.Error(fmt.Errorf("failed marshal to bytes: %w", err))
		common.WriteResponse([]byte("{\"message\":\"error\"}"), rw)
		return
	}
	common.WriteResponse(bytes, rw)
	logrus.Info("writing object: \n", string(bytes))
}

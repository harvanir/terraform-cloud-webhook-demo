package notification

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"harvanir/terraform-cloud-webhook-demo/pkg/json"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_api/common"
	"net/http"
)

const (
	runApplyAPI     = "https://app.terraform.io/api/v2/runs/%v/apply"
	stateVersionAPI = "https://app.terraform.io/api/v2/state-versions/%v"
)

// Handler handle notification API
func (ctx *Context) Notification(rw http.ResponseWriter, r *http.Request) {
	var payload Payload
	if err := json.Decode(&payload, r.Body); err != nil {
		logrus.Error(fmt.Errorf("failed decode to type: %w", err))
		writeErrorResponse(rw)
		return
	}
	if !isRunStatusApplied(payload.Notifications) {
		writeDefaultResponse(&payload, rw)
		return
	}
	downloadURLResponse, ok := callTerraform(rw, ctx, payload)
	if !ok {
		return
	}
	writeNotificationResponse(downloadURLResponse, rw)
}

func callTerraform(rw http.ResponseWriter, ctx *Context, payload Payload) (*HostedStateDownloadURLResponse, bool) {
	// call run apply webhook to get state version
	runApplyResponse, ok := callRunApplyWebhook(ctx, payload.RunID, rw)
	if !ok {
		return nil, false
	}
	stateVersionData := runApplyResponse.Data.Relationships.StateVersions.Data
	if len(stateVersionData) == 0 {
		writeErrorResponse(rw)
		return nil, false
	}
	// call state version webhook to get created ip public
	stateVersionResponse, ok := callStateVersionWebhook(stateVersionData[0].ID, ctx, rw)
	if !ok {
		writeErrorResponse(rw)
		return nil, false
	}
	// call hosted state download url
	hostedStateDownloadURL := stateVersionResponse.Data.Attributes.HostedStateDownloadURL
	hostedStateDownloadURLResponse, ok := callHostedStateDownloadURL(hostedStateDownloadURL, ctx, rw)
	if !ok {
		writeErrorResponse(rw)
		return nil, false
	}
	return hostedStateDownloadURLResponse, true
}

func writeErrorResponse(rw http.ResponseWriter) {
	handler.WriteJsonResponse([]byte("{\"message\":\"error\"}"), rw)
}

func isRunStatusApplied(notifications []Notification) bool {
	return len(notifications) > 0 && notifications[0].RunStatus == RunStatusApplied
}

func writeDefaultResponse(v interface{}, rw http.ResponseWriter) {
	byteArr, err := json.Marshal(v)
	if err != nil {
		logrus.Error(fmt.Errorf("failed marshal to bytes: %w", err))
		handler.WriteJsonResponse([]byte("{\"message\":\"error\"}"), rw)
		return
	}
	handler.WriteJsonResponse(byteArr, rw)
	logrus.Info("writing object: \n", string(byteArr))
}

func writeResponse(byteArr []byte, rw http.ResponseWriter) {
	handler.WriteJsonResponse(byteArr, rw)
	logrus.Info("writing object: \n", string(byteArr))
}

func doRequest(request *http.Request, rw http.ResponseWriter) (*http.Response, bool) {
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		logrus.Error("error do http request: ", err)
		return nil, false
	}
	return response, true
}

func callRunApplyWebhook(ctx *Context, runID string, rw http.ResponseWriter) (*ApplyResponse, bool) {
	url := fmt.Sprintf(runApplyAPI, runID)
	// construct http request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Error(fmt.Errorf("error create new http request: %w", err))
		writeErrorResponse(rw)
		return nil, false
	}
	constructAuthHeader(ctx, request)
	httpResponse, ok := doRequest(request, rw)
	if !ok {
		writeErrorResponse(rw)
		return nil, false
	}
	var response ApplyResponse
	if ok := decodeJsonResponse(&response, httpResponse, rw); ok {
		return &response, true
	}
	return nil, false
}

func constructAuthHeader(ctx *Context, request *http.Request) {
	token := ctx.Config.Terraform.Token
	request.Header.Set("Authorization", token.Bearer)
}

func callStateVersionWebhook(stateVersion string, ctx *Context, rw http.ResponseWriter) (*StateVersionsResponse, bool) {
	url := fmt.Sprintf(stateVersionAPI, stateVersion)
	// construct http request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Error(fmt.Errorf("error create new http request: %w", err))
		writeErrorResponse(rw)
		return nil, false
	}
	constructAuthHeader(ctx, request)
	httpResponse, ok := doRequest(request, rw)
	if !ok {
		writeErrorResponse(rw)
		return nil, false
	}
	var response StateVersionsResponse
	if ok := decodeJsonResponse(&response, httpResponse, rw); ok {
		return &response, true
	}
	return nil, false
}

func callHostedStateDownloadURL(url string, ctx *Context, rw http.ResponseWriter) (*HostedStateDownloadURLResponse, bool) {
	// construct http request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Error(fmt.Errorf("error create new http request: %w", err))
		writeErrorResponse(rw)
		return nil, false
	}
	constructAuthHeader(ctx, request)
	httpResponse, ok := doRequest(request, rw)
	if !ok {
		writeErrorResponse(rw)
		return nil, false
	}
	var response HostedStateDownloadURLResponse
	if ok := decodeTextResponse(&response, httpResponse, rw); !ok {
		return nil, false
	}
	var finalResponse HostedStateDownloadURLResponse
	funcCast := func(v interface{}) bool {
		hostedStateDownloadURLResource, ok := v.(HostedStateDownloadURLResponse)
		if !ok {
			return false
		}
		finalResponse = hostedStateDownloadURLResource
		return true
	}
	if ok && funcCast(response) {
		return &finalResponse, true
	}
	return nil, false
}

func decodeJsonResponse(v interface{}, httpResponse *http.Response, rw http.ResponseWriter) bool {
	if err := json.Decode(v, httpResponse.Body); err != nil {
		logrus.Error("error decode response from downstream: %w", err)
		writeErrorResponse(rw)
		return false
	}
	return true
}

func decodeTextResponse(response interface{}, httpResponse *http.Response, rw http.ResponseWriter) bool {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(httpResponse.Body)
	if err != nil {
		writeErrorResponse(rw)
		return false
	}
	if err := json.Unmarshal(buf.String(), response); err != nil {
		logrus.Error("error decode response from downstream: %w", err)
		writeErrorResponse(rw)
		return false
	}
	return true
}

func writeNotificationResponse(downloadURLResponse *HostedStateDownloadURLResponse, rw http.ResponseWriter) {
	resources := downloadURLResponse.Resources
	if isValidResources(resources) {
		writeErrorResponse(rw)
	}
	hostedStateResource := resources[0]
	hostedStateInstances := hostedStateResource.Instances
	publicIP := []byte(fmt.Sprintf(`{"public_ip":"%v"}`, hostedStateInstances[0].Attributes.PublicIP))
	writeResponse(publicIP, rw)
}

func isValidResources(resources []HostedStateDownloadURLResource) bool {
	return len(resources) == 0 || len(resources[0].Instances) == 0
}

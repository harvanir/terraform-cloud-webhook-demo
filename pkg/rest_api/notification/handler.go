package notification

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/sirupsen/logrus"
	"harvanir/terraform-cloud-webhook-demo/pkg/json"
	"harvanir/terraform-cloud-webhook-demo/pkg/rest_api/common"
	"net/http"
)

const (
	runApplyAPI = "https://app.terraform.io/api/v2/runs/%v/apply"
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
	// tfe context
	tfeContext := context.Background()
	config := &tfe.Config{Token: ctx.Config.Terraform.Token.Value}
	client, err := tfe.NewClient(config)
	if err != nil {
		logrus.Error("error new tfe client: ", err)
		writeErrorResponse(rw)
		return nil, false
	}
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
	// call state version
	stateVersion, err := client.StateVersions.Read(tfeContext, stateVersionData[0].ID)
	if err != nil {
		logrus.Error("error stateVersion: ", err)
		writeErrorResponse(rw)
		return nil, false
	}
	// call download
	download, err := client.StateVersions.Download(tfeContext, stateVersion.DownloadURL)
	if err != nil {
		logrus.Error("error download: ", err)
		writeErrorResponse(rw)
		return nil, false
	}
	var response HostedStateDownloadURLResponse
	err = json.UnmarshalByte(download, &response)
	if err != nil {
		logrus.Error("error unmarshal: ", err)
		writeErrorResponse(rw)
		return nil, false
	}
	return &response, true
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
	request.Header.Set("Authorization", "Bearer "+token.Value)
}

func decodeJsonResponse(v interface{}, httpResponse *http.Response, rw http.ResponseWriter) bool {
	if err := json.Decode(v, httpResponse.Body); err != nil {
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
	bytes, err := json.Marshal(&hostedStateResource)
	if err != nil {
		logrus.Error("error marshal: ", err)
		writeErrorResponse(rw)
	}
	writeResponse(bytes, rw)
}

func isValidResources(resources []HostedStateDownloadURLResource) bool {
	return len(resources) == 0 || len(resources[0].Instances) == 0
}

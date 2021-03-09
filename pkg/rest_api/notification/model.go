package notification

import (
	"harvanir/terraform-cloud-webhook-demo/pkg/config"
)

const (
	RunStatusApplied = "applied"
)

type Context struct {
	Config *config.AppConfig
}

// Payload https://www.terraform.io/docs/cloud/api/notification-configurations.html#notification-payload
type Payload struct {
	PayloadVersion              int            `json:"payload_version,omitempty"`
	NotificationConfigurationId string         `json:"notification_configuration_id,omitempty"`
	RunURL                      string         `json:"run_url,omitempty"`
	RunID                       string         `json:"run_id,omitempty"`
	RunMessage                  string         `json:"run_message,omitempty"`
	RunCreatedAt                string         `json:"run_created_at,omitempty"`
	RunCreatedBy                string         `json:"run_created_by,omitempty"`
	WorkspaceID                 string         `json:"workspace_id,omitempty"`
	WorkspaceName               string         `json:"workspace_name,omitempty"`
	OrganizationName            string         `json:"organization_name,omitempty"`
	Notifications               []Notification `json:"notifications,omitempty"`
}

type Notification struct {
	Message      string `json:"message,omitempty"`
	Trigger      string `json:"trigger,omitempty"`
	RunStatus    string `json:"run_status,omitempty"`
	RunUpdatedAt string `json:"run_updated_at,omitempty"`
	RunUpdatedBy string `json:"run_updated_by,omitempty"`
}

type ApplyResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
}

type Attributes struct {
	Status            string `json:"status"`
	ResourceAdditions int    `json:"resource-additions"`
}
type Relationships struct {
	StateVersions StateVersions `json:"state-versions"`
}
type StateVersions struct {
	Data []StateVersionsData `json:"data"`
}
type StateVersionsData struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes StateVersionsAttribute `json:"attributes,omitempty"`
}

type StateVersionsResponse struct {
	Data StateVersionsData `json:"data"`
}
type StateVersionsAttribute struct {
	HostedStateDownloadURL string `json:"hosted-state-download-url"`
}

type HostedStateDownloadURLResponse struct {
	Version   int                              `json:"version"`
	Resources []HostedStateDownloadURLResource `json:"resources"`
}

type HostedStateDownloadURLResource struct {
	Name      string                           `json:"name"`
	Instances []HostedStateDownloadURLInstance `json:"instances"`
}
type HostedStateDownloadURLInstance struct {
	SchemaVersion int                             `json:"schema_version"`
	Attributes    HostedStateDownloadURLAttribute `json:"attributes"`
}
type HostedStateDownloadURLAttribute struct {
	PublicIP string `json:"public_ip"`
}

func NewNotificationCtx(config *config.AppConfig) Context {
	return Context{Config: config}
}

package notification

type Payload struct {
	PayloadVersion              int            `json:"payload_version,omitempty"`
	NotificationConfigurationId string         `json:"notification_configuration_id,omitempty"`
	RunUrl                      string         `json:"run_url,omitempty"`
	RunId                       string         `json:"run_id,omitempty"`
	RunMessage                  string         `json:"run_message,omitempty"`
	RunCreatedAt                string         `json:"run_created_at,omitempty"`
	RunCreatedBy                string         `json:"run_created_by,omitempty"`
	WorkspaceId                 string         `json:"workspace_id,omitempty"`
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

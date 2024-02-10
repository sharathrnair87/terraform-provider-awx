package awx

import (
	"fmt"
	"net/http"
)

// This variable is mandatory and to be populated for creating services API
var mandatoryFields []string

// AWX represents awx api endpoints with services, and using
// client to communicate with awx server.
type AWX struct {
	client *Client

	ApplicationService                              *ApplicationService
	ExecutionEnvironmentsService                    *ExecutionEnvironmentsService
	PingService                                     *PingService
	InventoriesService                              *InventoriesService
	JobService                                      *JobService
	JobTemplateService                              *JobTemplateService
	JobTemplateNotificationTemplatesService         *JobTemplateNotificationTemplatesService
	ProjectService                                  *ProjectService
	ProjectUpdatesService                           *ProjectUpdatesService
	UserService                                     *UserService
	GroupService                                    *GroupService
	HostService                                     *HostService
	CredentialsService                              *CredentialsService
	CredentialTypeService                           *CredentialTypeService
	CredentialInputSourceService                    *CredentialInputSourceService
	InventorySourcesService                         *InventorySourcesService
	InventorySourcesSchedulesServices               *InventorySourcesSchedulesService
	InventoryGroupService                           *InventoryGroupService
	InstanceGroupsService                           *InstanceGroupsService
	NotificationTemplatesService                    *NotificationTemplatesService
	OrganizationsService                            *OrganizationsService
	ScheduleService                                 *SchedulesService
	SettingService                                  *SettingService
	TeamService                                     *TeamService
	TokenService                                    *TokenService
	WorkflowJobTemplateScheduleService              *WorkflowJobTemplateScheduleService
	WorkflowJobTemplateService                      *WorkflowJobTemplateService
	WorkflowJobTemplateNodeService                  *WorkflowJobTemplateNodeService
	WorkflowJobTemplateNodeAlwaysService            *WorkflowJobTemplateNodeStepService
	WorkflowJobTemplateNodeFailureService           *WorkflowJobTemplateNodeStepService
	WorkflowJobTemplateNodeSuccessService           *WorkflowJobTemplateNodeStepService
	WorkflowJobTemplateNotificationTemplatesService *WorkflowJobTemplateNotificationTemplatesService
}

// Client implement http client.
type Client struct {
	BaseURL   string
	Requester *Requester
}

// CheckResponse do http response check, and return err if not in [200, 300).
func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("responsed with %d, resp: %v", resp.StatusCode, resp)
}

// ValidateParams is to validate the input to use the services.
func ValidateParams(data map[string]interface{}, mandatoryFields []string) (notfound []string, status bool) {
	status = true
	for _, key := range mandatoryFields {
		_, exists := data[key]

		if !exists {
			notfound = append(notfound, key)
			status = false
		}
	}
	return notfound, status
}

// NewAWX news an awx handler with basic auth support, you could customize the http
// transport by passing custom client.
func NewAWX(baseURL, userName, passwd string, client *http.Client) (*AWX, error) {
	r := &Requester{Base: baseURL, Authenticator: &BasicAuth{Username: userName, Password: passwd}, Client: client}
	if r.Client == nil {
		r.Client = http.DefaultClient
	}

	awxClient := &Client{
		BaseURL:   baseURL,
		Requester: r,
	}

	newAWX := newAWX(awxClient)

	// test the connection and return and error if there's an issue
	_, err := newAWX.PingService.Ping()
	if err != nil {
		return nil, err
	}

	return newAWX, nil
}

// NewAWXToken creates an AWX handler with token support.
func NewAWXToken(baseURL, token string, client *http.Client) (*AWX, error) {
	r := &Requester{Base: baseURL, Authenticator: &TokenAuth{Token: token}, Client: client}
	if r.Client == nil {
		r.Client = http.DefaultClient
	}

	awxClient := &Client{
		BaseURL:   baseURL,
		Requester: r,
	}

	newAWX := newAWX(awxClient)

	// test the connection and return and error if there's an issue
	_, err := newAWX.PingService.Ping()
	if err != nil {
		return nil, err
	}

	return newAWX, nil
}

func newAWX(c *Client) *AWX {
	return &AWX{
		client: c,

		ApplicationService: &ApplicationService{
			client: c,
		},
		ExecutionEnvironmentsService: &ExecutionEnvironmentsService{
			client: c,
		},
		PingService: &PingService{
			client: c,
		},
		InventoriesService: &InventoriesService{
			client: c,
		},
		JobService: &JobService{
			client: c,
		},
		JobTemplateService: &JobTemplateService{
			client: c,
		},
		JobTemplateNotificationTemplatesService: &JobTemplateNotificationTemplatesService{
			client: c,
		},
		ProjectService: &ProjectService{
			client: c,
		},
		ProjectUpdatesService: &ProjectUpdatesService{
			client: c,
		},
		UserService: &UserService{
			client: c,
		},
		GroupService: &GroupService{
			client: c,
		},
		HostService: &HostService{
			client: c,
		},
		CredentialsService: &CredentialsService{
			client: c,
		},
		CredentialTypeService: &CredentialTypeService{
			client: c,
		},
		CredentialInputSourceService: &CredentialInputSourceService{
			client: c,
		},
		InventorySourcesService: &InventorySourcesService{
			client: c,
		},
		InventorySourcesSchedulesServices: &InventorySourcesSchedulesService{
			client: c,
		},
		InventoryGroupService: &InventoryGroupService{
			client: c,
		},
		InstanceGroupsService: &InstanceGroupsService{
			client: c,
		},
		NotificationTemplatesService: &NotificationTemplatesService{
			client: c,
		},
		OrganizationsService: &OrganizationsService{
			client: c,
		},
		ScheduleService: &SchedulesService{
			client: c,
		},
		SettingService: &SettingService{
			client: c,
		},
		TeamService: &TeamService{
			client: c,
		},
		TokenService: &TokenService{
			client: c,
		},
		WorkflowJobTemplateScheduleService: &WorkflowJobTemplateScheduleService{
			client: c,
		},
		WorkflowJobTemplateService: &WorkflowJobTemplateService{
			client: c,
		},
		WorkflowJobTemplateNodeService: &WorkflowJobTemplateNodeService{
			client: c,
		},
		WorkflowJobTemplateNodeSuccessService: &WorkflowJobTemplateNodeStepService{
			endpoint: fmt.Sprintf("%s%s", workflowJobTemplateNodeAPIEndpoint, "%d/success_nodes/"),
			client:   c,
		},
		WorkflowJobTemplateNodeFailureService: &WorkflowJobTemplateNodeStepService{
			endpoint: fmt.Sprintf("%s%s", workflowJobTemplateNodeAPIEndpoint, "%d/failure_nodes/"),
			client:   c,
		},
		WorkflowJobTemplateNodeAlwaysService: &WorkflowJobTemplateNodeStepService{
			endpoint: fmt.Sprintf("%s%s", workflowJobTemplateNodeAPIEndpoint, "%d/always_nodes/"),
			client:   c,
		},
		WorkflowJobTemplateNotificationTemplatesService: &WorkflowJobTemplateNotificationTemplatesService{
			client: c,
		},
	}
}

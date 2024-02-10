package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type CredentialTypeService struct {
	client *Client
}

type ListCredentialTypeResponse struct {
	Pagination
	Results []*CredentialType `json:"results"`
}

const CredentialTypesAPIEndpoint = "/api/v2/credential_types/"

func (cs *CredentialTypeService) ListCredentialTypes(params map[string]string) ([]*CredentialType, error) {
	results, err := cs.getAllCredentialTypes(CredentialTypesAPIEndpoint, params)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (cs *CredentialTypeService) CreateCredentialType(data map[string]interface{}, params map[string]string) (*CredentialType, error) {
	result := new(CredentialType)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PostJSON(CredentialTypesAPIEndpoint, bytes.NewReader(payload), result, params)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *CredentialTypeService) GetCredentialTypeByID(id int, params map[string]string) (*CredentialType, error) {
	result := new(CredentialType)
	endpoint := fmt.Sprintf("%s%d", CredentialTypesAPIEndpoint, id)
	resp, err := cs.client.Requester.GetJSON(endpoint, result, params)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *CredentialTypeService) GetCredentialTypeByName(params map[string]string) ([]*CredentialType, error) {
	result, err := cs.ListCredentialTypes(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *CredentialTypeService) UpdateCredentialTypeByID(id int, data map[string]interface{}, params map[string]string) (*CredentialType, error) {
	result := new(CredentialType)
	endpoint := fmt.Sprintf("%s%d", CredentialTypesAPIEndpoint, id)

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PutJSON(endpoint, bytes.NewReader(payload), result, params)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *CredentialTypeService) DeleteCredentialTypeByID(id int, params map[string]string) error {
	endpoint := fmt.Sprintf("%s%d", CredentialTypesAPIEndpoint, id)
	resp, err := cs.client.Requester.Delete(endpoint, nil, params)
	if err != nil {
		return err
	}

	err = CheckResponse(resp)
	if err != nil {
		return err
	}

	return nil
}

// make generic function
func (cs *CredentialTypeService) getAllCredentialTypes(firstURL string, params map[string]string) ([]*CredentialType, error) {
	results := make([]*CredentialType, 0)
	nextURL := firstURL
	for {
		nextURLParsed, err := url.Parse(nextURL)
		if err != nil {
			return nil, err
		}

		nextURLQueryParams := make(map[string]string)
		for paramName, paramValues := range nextURLParsed.Query() {
			if len(paramValues) > 0 {
				nextURLQueryParams[paramName] = paramValues[0]
			}
		}

		for paramName, paramValue := range params {
			nextURLQueryParams[paramName] = paramValue
		}

		result := new(ListCredentialTypeResponse)
		resp, err := cs.client.Requester.GetJSON(nextURLParsed.Path, result, nextURLQueryParams)
		if err != nil {
			return nil, err
		}

		if err := CheckResponse(resp); err != nil {
			return nil, err
		}

		results = append(results, result.Results...)

		if result.Next == nil || result.Next.(string) == "" {
			break
		}
		nextURL = result.Next.(string)
	}
	return results, nil
}

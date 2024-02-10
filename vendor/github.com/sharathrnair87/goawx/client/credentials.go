package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type CredentialsService struct {
	client *Client
}

type ListCredentialsResponse struct {
	Pagination
	Results []*Credential `json:"results"`
}

const CredentialsAPIEndpoint = "/api/v2/credentials/"

func (cs *CredentialsService) ListCredentials(params map[string]string) ([]*Credential, error) {
	results, err := cs.getAllPages(CredentialsAPIEndpoint, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (cs *CredentialsService) getAllPages(firstURL string, params map[string]string) ([]*Credential, error) {
	results := make([]*Credential, 0)
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

		result := new(ListCredentialsResponse)
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

func (cs *CredentialsService) CreateCredentials(data map[string]interface{}, params map[string]string) (*Credential, error) {
	result := new(Credential)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PostJSON(CredentialsAPIEndpoint, bytes.NewReader(payload), result, params)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *CredentialsService) GetCredentialsByID(id int, params map[string]string) (*Credential, error) {
	result := new(Credential)
	endpoint := fmt.Sprintf("%s%d", CredentialsAPIEndpoint, id)
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

func (cs *CredentialsService) UpdateCredentialsByID(id int, data map[string]interface{},
	params map[string]string) (*Credential, error) {
	result := new(Credential)
	endpoint := fmt.Sprintf("%s%d", CredentialsAPIEndpoint, id)

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, params)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *CredentialsService) DeleteCredentialsByID(id int, params map[string]string) error {
	endpoint := fmt.Sprintf("%s%d", CredentialsAPIEndpoint, id)
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

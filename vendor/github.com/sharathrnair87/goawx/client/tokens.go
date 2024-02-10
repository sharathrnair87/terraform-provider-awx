package awx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type TokenService struct {
	client *Client
}

type ListTokensResponse struct {
	Pagination
	Results []*Token `json:"results"`
}

const TokensAPIEndpoint = "/api/v2/tokens/"

// GetTokenByID shows the details of a project.
func (p *TokenService) GetTokenByID(id int, params map[string]string) (*Token, error) {
	result := new(Token)
	endpoint := fmt.Sprintf("%s%d/", TokensAPIEndpoint, id)
	resp, err := p.client.Requester.GetJSON(endpoint, result, params)
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateToken creates an awx project.
func (p *TokenService) CreateToken(data map[string]interface{}, params map[string]string) (*Token, error) {
	mandatoryFields = []string{"scope"}
	validate, status := ValidateParams(data, mandatoryFields)

	if !status {
		err := fmt.Errorf("Mandatory input arguments are absent: %s", validate)
		return nil, err
	}

	result := new(Token)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Add check if project exists and return proper error

	resp, err := p.client.Requester.PostJSON(TokensAPIEndpoint, bytes.NewReader(payload), result, params)
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateToken update an awx Token.
func (p *TokenService) UpdateToken(id int, data map[string]interface{}, params map[string]string) (*Token, error) {
	result := new(Token)
	endpoint := fmt.Sprintf("%s%d", TokensAPIEndpoint, id)
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Requester.PatchJSON(endpoint, bytes.NewReader(payload), result, nil)
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteToken delete an awx Token.
func (p *TokenService) DeleteToken(id int) (*Token, error) {
	result := new(Token)
	endpoint := fmt.Sprintf("%s%d", TokensAPIEndpoint, id)

	resp, err := p.client.Requester.Delete(endpoint, result, nil)
	if err != nil {
		return nil, err
	}

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	return result, nil
}

package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Validator is the instance used to perform validation
type Validator struct {
	// Client is the http Client used to perform API requests
	Client *http.Client

	token string
	url   string
}

// New creates a new Validator
func New(token string) Validator {
	return Validator{
		token:  token,
		Client: &http.Client{},
		url:    "https://validator-api.semweb.yandex.ru",
	}
}

func (v Validator) call(ctx context.Context, method, path, body string, queryParams url.Values, resp interface{}) error {
	if queryParams == nil {
		queryParams = url.Values{}
	}
	queryParams.Set("api_key", v.token)
	url := fmt.Sprintf("%s%s?%s", v.url, path, queryParams.Encode())
	req, errReq := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	if errReq != nil {
		return errReq
	}
	httpResp, errResp := v.Client.Do(req)
	if errResp != nil {
		return errResp
	}
	defer httpResp.Body.Close()
	errJSON := json.NewDecoder(httpResp.Body).Decode(resp)
	if errJSON != nil {
		return errJSON
	}
	return nil
}

// CheckDocument performs a document check, according to
// https://tech.yandex.com/validator/doc/dg/concepts/html-validation-docpage/
func (v Validator) CheckDocument(ctx context.Context, document string) (StandardResponse, error) {
	resp := StandardResponse{}
	return resp, v.call(ctx, "POST", "/1.1/document_parser", document, nil, &resp)
}

// StandardResponse is the validation response from Yandex
type StandardResponse struct {
	ID   string `json:"id"`
	Data struct {
		Microdata   []map[string]interface{} `json:"microdata"`
		RDFA        []map[string]interface{} `json:"rdfa"`
		Microformat []map[string]interface{} `json:"microformat"`
		JSONLD      []map[string]interface{} `json:"json-ld"`
	} `json:"data"`
}

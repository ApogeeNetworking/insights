package insights

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GenerateRequest ...
func (s *Service) generateRequest(uri, method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, s.baseURL+uri, body)
	if err != nil {
		return nil, fmt.Errorf("unabled to create request: %v", err)
	}
	if s.apiToken != "" {
		req.Header.Set("Api-Token", s.apiToken)
	}
	if method != "GET" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// CreateReqBody ...
func (s *Service) createReqBody(v interface{}) (*bytes.Reader, error) {
	payload, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payload), nil
}

// MakeRequest ...
func (s *Service) makeRequest(req *http.Request) (*http.Response, error) {
	resp, err := s.http.Do(req)
	switch {
	case err != nil:
		return resp, err
	case resp.StatusCode == 401:
		return resp, errors.New(resp.Status)
	case resp.StatusCode >= 400 && resp.StatusCode <= 599:
		return resp, errors.New(resp.Status)
	}
	return resp, nil
}

package insights

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Service ..
type Service struct {
	http     *http.Client
	baseURL  string
	apiToken string
}

// NewService ..
func NewService(host string, apiToken string, insecureSSL bool) *Service {
	return &Service{
		baseURL: host,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureSSL,
				},
			},
			Timeout: 90 * time.Second,
		},
		apiToken: apiToken,
	}
}

type AuthParams struct {
	Username string `json:"username,omitempty"`
}

type AuthResp struct {
	ApiToken string `json:"data"`
}

func (s *Service) Auth() error {
	var authRes AuthResp
	sr, _ := s.createReqBody(&AuthParams{})
	req, err := s.generateRequest("/auth", "POST", sr)
	if err != nil {
		return err
	}
	res, err := s.makeRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&authRes); err != nil {
		return err
	}
	fmt.Println(authRes)
	s.apiToken = authRes.ApiToken
	return nil
}

type GetSchoolsResp struct {
	Data []School `json:"data"`
	Meta struct {
		NextURL string `json:"next_url"`
	} `json:"meta"`
}

type GetSchoolResp struct {
	Data School `json:"data"`
}

type School struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Activated     bool      `json:"activated"`
	ShortName     string    `json:"apogee_short_internal_name"`
	ActivatedDate time.Time `json:"activated_timestamp"`
}

func (s *Service) GetSchools() ([]School, error) {
	var resp GetSchoolsResp
	req, err := s.generateRequest("/schools", "GET", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("limit", "20")
	q.Add("activated_status", "active")
	req.URL.RawQuery = q.Encode()
	res, err := s.makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (s *Service) GetSchool(id string) (School, error) {
	var resp GetSchoolResp
	req, err := s.generateRequest(
		fmt.Sprintf("/schools/%s", id), "GET", nil,
	)
	if err != nil {
		return School{}, err
	}
	res, err := s.makeRequest(req)
	if err != nil {
		return School{}, err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return School{}, err
	}
	return resp.Data, nil
}

type AccessPoint struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_timestamp"`
	Serial    string    `json:"serial"`
	Switch    string    `json:"switch"`
	Floor     int       `json:"floor"`
	Building  string    `json:"building"`
	Room      string    `json:"room"`
	Location  struct {
		Latitude  int64 `json:"latitude"`
		Longitude int64 `json:"longitude"`
	} `json:"location"`
	LastHealthCheck time.Time `json:"last_health_check_timestamp"`
}

type GetAccessPointsResp struct {
	Data []AccessPoint `json:"data"`
	Meta struct {
		NextURL string `json:"next_url"`
	} `json:"meta"`
}

type SyncAp struct {
	Name    string `json:"name"`
	MacAddr string `json:"mac_address"`
	Serial  string `json:"serial"`
	Switch  string `json:"switch"`
}

// SyncAps synchronizes accesspoints to the DA DB
// Only allows pushing 100 APs at a time
func (s *Service) SyncAps(schoolID string, aps []SyncAp) (bool, error) {
	d, _ := json.Marshal(aps)
	data := strings.NewReader(string(d))
	req, err := s.generateRequest(
		fmt.Sprintf("/schools/%s/access_points/sync", schoolID),
		"POST",
		data,
	)
	if err != nil {
		return false, err
	}
	res, err := s.makeRequest(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	resp := struct {
		Data bool `json:"data"`
	}{}
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return false, err
	}
	return resp.Data, nil
}

// BulkSyncAps ...
func (s *Service) BulkSyncAps(schoolID string, aps []SyncAp) (bool, error) {
	var (
		results []bool
		syncAps []SyncAp
	)
	for _, ap := range aps {
		syncAps = append(syncAps, ap)
		if len(syncAps) == 100 {
			p, err := s.SyncAps(schoolID, syncAps)
			if err != nil {
				return false, err
			}
			results = append(results, p)
			syncAps = nil
		}
	}
	if len(syncAps) != 0 {
		p, err := s.SyncAps(schoolID, syncAps)
		if err != nil {
			return false, err
		}
		results = append(results, p)
	}
	success := true
	if len(results) == 0 {
		success = false
	}
	for _, sync := range results {
		if !sync {
			success = false
			break
		}
	}
	return success, nil
}

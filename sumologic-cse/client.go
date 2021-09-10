package sumologic_cse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

func NewClient(host, apiKey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    *host,
	}

	if apiKey == nil {
		return nil, errors.New("apiKey is not set")
	}

	if host == nil {
		return nil, errors.New("host is not set")
	}

	c.ApiKey = *apiKey
	return &c, nil
}

// ----------------------------------------------------------------------------
// Helper Functions
// ----------------------------------------------------------------------------

func (c *Client) translateToPermissionIds(pemissionNames []interface{}) ([]string, error) {
	pl := make([]string, 0, len(pemissionNames))

	permissions, err := c.ReadAll("permissions")
	if err != nil {
		return nil, err
	}

	for _, permission := range permissions.(PermissionListResponse).Data.Permissions {
		for _, name := range pemissionNames {
			if permission.Name == name {
				pl = append(pl, permission.Id)
				continue
			}
		}
	}

	return pl, nil
}

// ----------------------------------------------------------------------------
// CRUD Functions
// ----------------------------------------------------------------------------

func (c *Client) ReadAll(objectType string) (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?offset=0&limit=100", c.HostURL, objectType), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	switch objectType {
	case "custom-insights":
		var data CustomInsightListResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case "permissions":
		var data PermissionListResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case "users":
		var data UsersListResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	default:
		return nil, errors.New("type not expected")
	}
}

func (c *Client) Create(data interface{}) (string, error) {
	var objectType string
	var payload []byte
	var err error

	switch data.(type) {
	case CustomEntityTypeRequest:
		objectType = "custom-entity-types"
		payload, err = json.Marshal(data.(CustomEntityTypeRequest))
		if err != nil {
			return "", err
		}
	case CustomInsightRequest:
		objectType = "custom-insights"
		payload, err = json.Marshal(data.(CustomInsightRequest))
		if err != nil {
			return "", err
		}
	case NetworkBlockRequest:
		objectType = "network-blocks"
		payload, err = json.Marshal(data.(NetworkBlockRequest))
		if err != nil {
			return "", err
		}
	case RoleRequest:
		objectType = "roles"
		payload, err = json.Marshal(data.(RoleRequest))
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("type not expected")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, objectType), bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	switch data.(type) {
	case CustomEntityTypeRequest:
		return data.(CustomEntityTypeResponse).CustomEntityType.Id, nil
	case CustomInsightRequest:
		return data.(CustomInsightResponse).CustomInsight.Id, nil
	case NetworkBlockRequest:
		return data.(NetworkBlockResponse).NetworkBlock.Id, nil
	case RoleRequest:
		return data.(RoleResponse).Role.Id, nil
	}

	return "", errors.New("type not expected")
}

func (c *Client) Read(objectType string, id string) (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, objectType, id), nil)
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	switch objectType {
	case "custom-entity-types":
		var data CustomEntityTypeResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case "network-blocks":
		var data NetworkBlockResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case "roles":
		var data RoleResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	default:
		return nil, errors.New("type not expected")
	}
}

func (c *Client) Update(id string, data interface{}) error {
	var objectType string
	var payload []byte
	var err error

	switch data.(type) {
	case CustomEntityTypeRequest:
		objectType = "custom-entity-types"
		payload, err = json.Marshal(data.(CustomEntityTypeRequest))
		if err != nil {
			return err
		}
	case CustomInsightRequest:
		objectType = "custom-insights"
		payload, err = json.Marshal(data.(CustomInsightRequest))
		if err != nil {
			return err
		}
	case NetworkBlockRequest:
		objectType = "network-blocks"
		payload, err = json.Marshal(data.(NetworkBlockRequest))
		if err != nil {
			return err
		}
	case RoleRequest:
		objectType = "roles"
		payload, err = json.Marshal(data.(RoleRequest))
		if err != nil {
			return err
		}
	default:
		return errors.New("type not expected")
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", c.HostURL, objectType, id), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	switch data.(type) {
	case CustomEntityTypeRequest:
		var data CustomEntityTypeResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	case CustomInsightRequest:
		var data CustomInsightRequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	case NetworkBlockRequest:
		var data NetworkBlockResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	case RoleRequest:
		var data RoleResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("type not expected")
	}
}

func (c *Client) Delete(id, objectType string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, objectType, id), nil)

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-API-Key", c.ApiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

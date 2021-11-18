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

const (
	AggregationRule string  = "rules/aggregation"
	CustomEntityTypes       = "custom-entity-types"
	CustomInsights          = "custom-insights"
	LogMappings			 	= "log-mappings"
	MatchLists 			    = "match-lists"
	MatchListItems 		    = "match-list-items"
	MatchRule 				= "rules/templated"
	NetworkBlocks           = "network-blocks"
	Permissions			 	= "permissions"
	Roles  			        = "roles"
	Rules  			        = "rules"
	Users  			        = "users"
	ThresholdRule           = "rules/threshold"
)

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

func (c *Client) Enabled(id, objectType string, enable bool) error {
	payload := []byte(fmt.Sprintf(`{"enabled": %t}`, enable))
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/%s/enabled", c.HostURL, objectType, id), bytes.NewBuffer(payload))

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getListItemId(listId, itemValue string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/match-list-items", c.HostURL), nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("listIds", listId)
	q.Add("value", itemValue)
	req.URL.RawQuery = q.Encode()

	resp, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	result := MatchListItemGetResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return "", err
	}

	if result.Data.Total != 1 {
		return "", errors.New("expected exactly one item to be returned")
	}

	return result.Data.Objects[0].Id, nil
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
	case CustomInsights:
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
		return nil, errors.New("ReadAll: type not expected")
	}
}

func (c *Client) Create(data interface{}) (string, error) {
	var objectType string
	var payload []byte
	var err error

	switch data.(type) {
	case CustomEntityTypeRequest:
		objectType = CustomEntityTypes
		payload, err = json.Marshal(data.(CustomEntityTypeRequest))
		if err != nil {
			return "", err
		}
	case CustomInsightRequest:
		objectType = CustomInsights
		payload, err = json.Marshal(data.(CustomInsightRequest))
		if err != nil {
			return "", err
		}
	case NetworkBlockRequest:
		objectType = NetworkBlocks
		payload, err = json.Marshal(data.(NetworkBlockRequest))
		if err != nil {
			return "", err
		}
	case RoleRequest:
		objectType = Roles
		payload, err = json.Marshal(data.(RoleRequest))
		if err != nil {
			return "", err
		}
	case AggregationRuleRequest:
		objectType = AggregationRule
		payload, err = json.Marshal(data.(AggregationRuleRequest))
		if err != nil {
			return "", err
		}
	case TemplatedRuleRequest:
		objectType = MatchRule
		payload, err = json.Marshal(data.(TemplatedRuleRequest))
		if err != nil {
			return "", err
		}
	case ThresholdRuleRequest:
		objectType = ThresholdRule
		payload, err = json.Marshal(data.(ThresholdRuleRequest))
		if err != nil {
			return "", err
		}
	case MatchListCreateRequest:
		objectType = MatchLists
		payload, err = json.Marshal(data.(MatchListCreateRequest))
		if err != nil {
			return "", err
		}
	case MatchListItemCreateRequest:
		payload, err = json.Marshal(data.(MatchListItemCreateRequest).Payload)
		objectType = fmt.Sprintf("match-lists/%s/items", data.(MatchListItemCreateRequest).ListId)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("type not expected when preparing creation request")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.HostURL, objectType), bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	switch data.(type) {
	case CustomEntityTypeRequest:
		result := CustomEntityTypeResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.CustomEntityType.Id, nil
	case CustomInsightRequest:
		result := CustomInsightResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.CustomInsight.Id, nil
	case NetworkBlockRequest:
		result := NetworkBlockResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.NetworkBlock.Id, nil
	case RoleRequest:
		result := RoleResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.Role.Id, nil
	case AggregationRuleRequest:
		result := RuleResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.Rule.Id, nil
	case TemplatedRuleRequest:
		result := RuleResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.Rule.Id, nil
	case ThresholdRuleRequest:
		result := RuleResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.Rule.Id, nil
	case MatchListCreateRequest:
		result := MatchListResponse{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return "", err
		}
		return result.MatchList.Id, nil
	case MatchListItemCreateRequest:
		time.Sleep(5 * time.Second) // This will do for now
		return c.getListItemId(data.(MatchListItemCreateRequest).ListId, data.(MatchListItemCreateRequest).Payload.Items[0].Value)
	}

	return "", errors.New("type not expected when processing creation response")
}

func (c *Client) Read(objectType string, id string) (interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", c.HostURL, objectType, id), nil)

	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	switch objectType {
	case CustomEntityTypes:
		var data CustomEntityTypeResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case LogMappings:
		var data LogMappingResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case MatchLists:
		var data MatchListResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case MatchListItems:
		var data MatchListItemResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case NetworkBlocks:
		var data NetworkBlockResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case Roles:
		var data RoleResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	case Rules:
		var data RuleResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	default:
		return nil, errors.New(fmt.Sprintf("type not expected when reading for object %s", objectType))
	}
}

func (c *Client) Update(id string, data interface{}) error {
	var objectType string
	var payload []byte
	var err error

	switch data.(type) {
	case CustomEntityTypeRequest:
		objectType = CustomEntityTypes
		payload, err = json.Marshal(data.(CustomEntityTypeRequest))
		if err != nil {
			return err
		}
	case CustomInsightRequest:
		objectType = CustomInsights
		payload, err = json.Marshal(data.(CustomInsightRequest))
		if err != nil {
			return err
		}
	case LogMappingRequest:
		objectType = LogMappings
		payload, err = json.Marshal(data.(LogMappingRequest))
		if err != nil {
			return err
		}
	case MatchListCreateRequest:
		objectType = MatchLists
		payload, err = json.Marshal(data.(MatchListCreateRequest))
		if err != nil {
			return err
		}
	case MatchListItemCreateRequest:
		objectType = MatchListItems
		payload, err = json.Marshal(data.(MatchListItemCreateRequest))
		if err != nil {
			return err
		}
	case NetworkBlockRequest:
		objectType = NetworkBlocks
		payload, err = json.Marshal(data.(NetworkBlockRequest))
		if err != nil {
			return err
		}
	case RoleRequest:
		objectType = Roles
		payload, err = json.Marshal(data.(RoleRequest))
		if err != nil {
			return err
		}
	case AggregationRuleRequest:
		objectType = AggregationRule
		payload, err = json.Marshal(data.(AggregationRuleRequest))
		if err != nil {
			return err
		}
	case TemplatedRuleRequest:
		objectType = MatchRule
		payload, err = json.Marshal(data.(TemplatedRuleRequest))
		if err != nil {
			return err
		}
	case ThresholdRuleRequest:
		objectType = ThresholdRule
		payload, err = json.Marshal(data.(ThresholdRuleRequest))
		if err != nil {
			return err
		}
	default:
		return errors.New("type not expected when preparing update request")
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/%s", c.HostURL, objectType, id), bytes.NewBuffer(payload))
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
	case LogMappingRequest:
		var data LogMappingResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	case MatchListUpdateRequest:
		var data MatchListResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	case MatchListItemUpdateRequest:
		var data MatchListItemResponse
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
	case AggregationRuleRequest:
		var data RuleResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	case TemplatedRuleRequest:
		var data RuleResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	case ThresholdRuleRequest:
		var data RuleResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("type not expected when processing update response")
	}
}

func (c *Client) Delete(id, objectType string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", c.HostURL, objectType, id), nil)

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


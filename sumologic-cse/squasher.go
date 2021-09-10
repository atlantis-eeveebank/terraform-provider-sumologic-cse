package sumologic_cse

import (
	"errors"
)

func flattenData(data interface{}) ([]map[string]interface{}, error) {
	var flattenData []map[string]interface{}

	switch d := data.(type) {
	case []LogMappingField:
		for _, f := range d {
			field := make(map[string]interface{})
			field["name"] = f.Name
			field["value"] = f.Value
			field["format"] = f.Format
			field["value_type"] = f.ValueType
			flattenData = append(flattenData, field)
		}
	case LogMappingInput:
		input := make(map[string]interface{})
		input["event_id_pattern"] = d.EventIdPattern
		input["log_format"] = d.LogFormat
		input["product"] = d.Product
		input["vendor"] = d.Vendor
		input["alternatives"] = d.Alternatives
		flattenData = append(flattenData, input)
	case []LogMappingStructuredInput:
		for _, i := range d {
			input := make(map[string]interface{})
			input["event_id_pattern"] = i.EventIdPattern
			input["log_format"] = i.LogFormat
			input["product"] = i.Product
			input["vendor"] = i.Vendor
			flattenData = append(flattenData, input)
		}
	case []Permission:
		for _, p := range d {
			permission := make(map[string]interface{})
			permission["id"] = p.Id
			permission["name"] = p.Name
			flattenData = append(flattenData, permission)
		}
	case []User:
		for _, u := range d {
			user := make(map[string]interface{})
			user["display_name"] = u.DisplayName
			user["email"] = u.Email
			user["mfa_enabled"] = u.MfaEnabled
			user["permissions"] = u.Permissions
			user["role"] = u.Role
			user["teams"] = u.Teams
			user["username"] = u.Username
			user["id"] = u.Id
			flattenData = append(flattenData, user)
		}
	default:
		return nil, errors.New("type not expected by squasher")
	}

	return flattenData, nil
}

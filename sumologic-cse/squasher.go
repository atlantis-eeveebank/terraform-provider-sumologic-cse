package sumologic_cse

import (
	"errors"
)

func flattenData(data interface{}) ([]map[string]interface{}, error) {
	var flattenData []map[string]interface{}

	switch d := data.(type) {
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
		return nil, errors.New("type not expected")
	}

	return flattenData, nil
}

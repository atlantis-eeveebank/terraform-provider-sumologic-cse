package sumologic_cse

import (
	"errors"
)

func flattenData(data interface{}) ([]map[string]interface{}, error) {
	var flattenData []map[string]interface{}

	switch d := data.(type) {
	case []LogMappingField:
		for _, f := range d {
			var field = map[string]interface{}{
				"name":             f.Name,
				"value":            f.Value,
				"format":           f.Format,
				"value_type":       f.ValueType,
				"alternate_values": f.AlternateValues,
				"case_insensitive": f.CaseInsensitive,
				"default_value":    f.DefaultValue,
				"lookup":           f.Lookup,
			}

			flattenData = append(flattenData, field)
		}
	case []LogMappingStructuredInput:
		for _, i := range d {
			var input = map[string]interface{}{
				"event_id_pattern": i.EventIdPattern,
				"log_format":       i.LogFormat,
				"product":          i.Product,
				"vendor":           i.Vendor,
			}
			flattenData = append(flattenData, input)
		}
	case []Permission:
		for _, p := range d {
			var permission = map[string]interface{}{
				"id":   p.Id,
				"name": p.Name,
			}
			flattenData = append(flattenData, permission)
		}
	case []User:
		for _, u := range d {
			var user = map[string]interface{}{
				"display_name": u.DisplayName,
				"email":        u.Email,
				"mfa_enabled":  u.MfaEnabled,
				"permissions":  u.Permissions,
				"role":         u.Role,
				"teams":        u.Teams,
				"username":     u.Username,
				"id":           u.Id,
			}
			flattenData = append(flattenData, user)
		}
	default:
		return nil, errors.New("type not expected by squasher")
	}

	return flattenData, nil
}

package sumologic_cse

import (
	"errors"
)

func flattenData(data interface{}) ([]map[string]interface{}, error) {
	var flattenData []map[string]interface{}

	switch d := data.(type) {
	case []LogMappingField:
		for _, f := range d {
			var lookups []map[string]interface{}
			for _, l := range f.Lookup {
				lookup := map[string]interface{}{
					"key":   l.Key,
					"value": l.Value,
				}
				lookups = append(lookups, lookup)
			}

			var field = map[string]interface{}{
				"name":             f.Name,
				"value":            f.Value,
				"format":           f.Format,
				"value_type":       f.ValueType,
				"alternate_values": f.AlternateValues,
				"case_insensitive": f.CaseInsensitive,
				"default_value":    f.DefaultValue,
				"lookup":           lookups,
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
	case []RuleAggregationFunction:
		for _, i := range d {
			var input = map[string]interface{}{
				"arguments": i.Arguments,
				"function":  i.Function,
				"name":      i.Name,
			}
			flattenData = append(flattenData, input)
		}
	case []RuleEntitySelector:
		for _, i := range d {
			var input = map[string]interface{}{
				"entity_type": i.EntityType,
				"expression":  i.Expression,
			}
			flattenData = append(flattenData, input)
		}
	case RuleScoreMapping:
		var input = map[string]interface{}{
			"default": d.Default,
			"type":    d.Type,
		}
		flattenData = append(flattenData, input)
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
		return nil, errors.New("type not expected when attempting to flatten data")
	}

	return flattenData, nil
}

func toStringSlice(data interface{}) []string {
	datum := data.([]interface{})

	ss := make([]string, len(datum))
	for i, v := range datum {
		ss[i] = v.(string)
	}

	return ss
}

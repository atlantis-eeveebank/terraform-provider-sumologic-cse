package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_CSE_API_KEY", nil),
			},
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_CSE_HOST", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sumologiccse_log_mapping": dataSourceLogMapping(),
			"sumologiccse_permissions": dataSourcePermissions(),
			"sumologiccse_users":       dataSourceUsers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"sumologiccse_role":               resourceRole(),
			"sumologiccse_custom_entity_type": resourceCustomEntityType(),
			"sumologiccse_custom_insight":     resourceCustomInsight(),
			"sumologiccse_log_mapping":        resourceLogMapping(),
			"sumologiccse_network_block":      resourceNetworkBlock(),
			"sumologiccse_aggregation_rule":   resourceAggregationRule(),
			"sumologiccse_match_rule":         resourceTemplatedRule(),
			"sumologiccse_match_list":         resourceMatchList(),
			"sumologiccse_match_list_item":    resourceMatchListItem(),
			"sumologiccse_threshold_rule":     resourceThresholdRule(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiKey := d.Get("api_key").(string)
	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	if apiKey == "" {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sumologic client",
			Detail:   "API Key is empty.",
		})
	}

	c, err := NewClient(host, &apiKey)
	if err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Sumologic client",
			Detail:   err.Error(),
		})
	}

	return c, nil
}

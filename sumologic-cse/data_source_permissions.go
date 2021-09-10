package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

type PermissionListResponse struct {
	Data struct {
		Permissions []Permission `json:"objects"`
	} `json:"data"`
}

type Permission struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func dataSourcePermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePermissionRead,
		Schema: map[string]*schema.Schema{
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)
	result, err := c.ReadAll("permissions")
	if err != nil {
		return diag.FromErr(err)
	}

	data, err := flattenData(result.(PermissionListResponse).Data.Permissions)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("permissions", data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

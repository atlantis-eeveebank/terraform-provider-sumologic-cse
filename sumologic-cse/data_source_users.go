package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

type UsersListResponse struct {
	Data struct {
		Users []User `json:"objects"`
	} `json:"data"`
}

type User struct {
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email"`
	Id          string   `json:"id"`
	MfaEnabled  bool     `json:"mfaEnabled"`
	Permissions []string `json:"permissions"`
	Role        string   `json:"role"`
	Teams       []string `json:"teams"`
	Username    string   `json:"username"`
}

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"mfa_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"permissions": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schema.TypeString,
						},
						"role": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"teams": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schema.TypeString,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)
	result, err := c.ReadAll(Users)
	if err != nil {
		return diag.FromErr(err)
	}

	data, err := flattenData(result.(UsersListResponse).Data.Users)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("users", data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

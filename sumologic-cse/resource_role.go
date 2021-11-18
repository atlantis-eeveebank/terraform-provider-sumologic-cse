package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RoleResponse struct {
	Role Role `json:"data"`
}

type Role struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type RoleRequest struct {
	Fields RolePayload `json:"fields"`
}

type RolePayload struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	permissions, err := c.translateToPermissionIds(d.Get("permissions").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := c.Create(RoleRequest{
		Fields: RolePayload{
			Name:        d.Get("name").(string),
			Permissions: permissions,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceRoleRead(ctx, d, m)

	return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	role, err := c.Read(Roles, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	permissions := make([]string, 0, len(role.(RoleResponse).Role.Permissions))
	for _, p := range role.(RoleResponse).Role.Permissions {
		permissions = append(permissions, p.Name)
	}

	err = d.Set("permissions", permissions)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", role.(RoleResponse).Role.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("name", "permissions") {
		c := m.(*Client)

		permissions, err := c.translateToPermissionIds(d.Get("permissions").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		err = c.Update(d.Id(), RoleRequest{
			Fields: RolePayload{
				Name:        d.Get("name").(string),
				Permissions: permissions,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), Roles)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomEntityTypeResponse struct {
	CustomEntityType CustomEntityType `json:"data"`
}

type CustomEntityType struct {
	Id         string   `json:"id"`
	Identifier string   `json:"identifier"`
	Name       string   `json:"name"`
	Fields     []string `json:"fields"`
}

type CustomEntityTypeRequest struct {
	Fields CustomEntityTypePayload `json:"fields"`
}

type CustomEntityTypePayload struct {
	Identifier string   `json:"identifier"`
	Name       string   `json:"name"`
	Fields     []string `json:"fields"`
}

func resourceCustomEntityType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomEntityTypeCreate,
		ReadContext:   resourceCustomEntityTypeRead,
		UpdateContext: resourceCustomEntityTypeUpdate,
		DeleteContext: resourceCustomEntityTypeDelete,
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
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"fields": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCustomEntityTypeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	setFields := d.Get("Items").([]interface{})
	fields := make([]string, len(setFields))
	for _, field := range setFields {
		fields = append(fields, field.(string))
	}

	id, err := c.Create(CustomEntityTypeRequest{
		Fields: CustomEntityTypePayload{
			Identifier: d.Get("identifier").(string),
			Name:       d.Get("name").(string),
			Fields:     fields,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceRoleRead(ctx, d, m)

	return diags
}

func resourceCustomEntityTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	data, err := c.Read(CustomEntityTypes, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("fields", data.(CustomEntityTypeResponse).CustomEntityType.Fields)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("identifier", data.(CustomEntityTypeResponse).CustomEntityType.Identifier)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", data.(CustomEntityTypeResponse).CustomEntityType.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCustomEntityTypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("name", "fields") {
		c := m.(*Client)

		setFields := d.Get("Items").([]interface{})
		fields := make([]string, len(setFields))
		for _, field := range setFields {
			fields = append(fields, field.(string))
		}

		err := c.Update(d.Id(), CustomEntityTypeRequest{
			Fields: CustomEntityTypePayload{
				Name:   d.Get("name").(string),
				Fields: fields,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCustomEntityTypeRead(ctx, d, m)
}

func resourceCustomEntityTypeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), CustomEntityTypes)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

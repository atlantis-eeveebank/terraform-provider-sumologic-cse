package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

type MatchListResponse struct {
	MatchList MatchList `json:"data"`
}

type MatchList struct {
	Created       time.Time `json:"created"`
	CreatedBy     string    `json:"createdBy"`
	DefaultTtl    int       `json:"defaultTtl"`
	Description   string    `json:"description"`
	Id            string    `json:"id"`
	LastUpdated   time.Time `json:"lastUpdated"`
	LastUpdatedBy string    `json:"lastUpdatedBy"`
	Name          string    `json:"name"`
	TargetColumn  string    `json:"targetColumn"`
}

type MatchListCreateRequest struct {
	Fields MatchListCreatePayload `json:"fields"`
}

type MatchListUpdateRequest struct {
	Fields MatchListUpdatePayload `json:"fields"`
}

type MatchListCreatePayload struct {
	Active       bool   `json:"active"`
	DefaultTtl   int    `json:"defaultTtl"`
	Description  string `json:"description"`
	Name         string `json:"name"`
	TargetColumn string `json:"targetColumn"`
}

type MatchListUpdatePayload struct {
	Active       bool   `json:"active"`
	DefaultTtl   int    `json:"defaultTtl"`
	Description  string `json:"description"`
}

func resourceMatchList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchListCreate,
		ReadContext:   resourceMatchListRead,
		UpdateContext: resourceMatchListUpdate,
		DeleteContext: resourceMatchListDelete,
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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"default_ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"target_column": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMatchListCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(MatchListCreateRequest{
		Fields: MatchListCreatePayload{
			Active:       true,
			DefaultTtl:   d.Get("default_ttl").(int),
			Description:  d.Get("description").(string),
			Name:         d.Get("name").(string),
			TargetColumn: d.Get("target_column").(string),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceMatchListRead(ctx, d, m)

	return diags
}

func resourceMatchListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	matchList, err := c.Read(MatchLists, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", matchList.(MatchListResponse).MatchList.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("description", matchList.(MatchListResponse).MatchList.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("default_ttl", matchList.(MatchListResponse).MatchList.DefaultTtl)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("target_column", matchList.(MatchListResponse).MatchList.TargetColumn)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceMatchListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges( "description", "default_ttl") {
		c := m.(*Client)

		err := c.Update(d.Id(), MatchListUpdateRequest{
			Fields: MatchListUpdatePayload{
				Active:      true,
				DefaultTtl:  d.Get("default_ttl").(int),
				Description: d.Get("description").(string),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMatchListRead(ctx, d, m)
}

func resourceMatchListDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), MatchLists)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

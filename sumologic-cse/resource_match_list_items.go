package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MatchListItemGetResponse struct {
	Data struct {
		Objects []struct {
			Id string `json:"id"`
		} `json:"objects"`
		Total int `json:"total"`
	} `json:"data"`
}

type MatchListItemResponse struct {
	MatchListItem MatchListItem `json:"data"`
}

type MatchListItem struct {
	Active     bool   `json:"active"`
	Expiration string `json:"expiration"`
	Id         string `json:"id"`
	ListName   string `json:"listName"`
	Value      string `json:"value"`

	Meta struct {
		Description string `json:"description"`
	} `json:"meta"`
}

type MatchListItemCreateRequest struct {
	ListId  string `json:"omit"`
	Payload MatchListItemCreatePayload
}

type MatchListItemCreatePayload struct {
	Items []MatchListItemPayload `json:"items"`
}

type MatchListItemPayload struct {
	Active      bool   `json:"active"`
	Description string `json:"description"`
	Expiration  string `json:"expiration,omitempty"`
	Value       string `json:"value"`
}

type MatchListItemUpdateRequest struct {
	Fields MatchListItemUpdatePayload `json:"fields"`
}

type MatchListItemUpdatePayload struct {
	Active      bool   `json:"active"`
	Description string `json:"description"`
	Expiration  string `json:"expiration"`
}

func resourceMatchListItem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchListItemCreate,
		ReadContext:   resourceMatchListItemRead,
		UpdateContext: resourceMatchListItemUpdate,
		DeleteContext: resourceMatchListItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"match_list_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expiration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMatchListItemCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(MatchListItemCreateRequest{
		ListId: d.Get("match_list_id").(string),
		Payload: MatchListItemCreatePayload{
			Items: []MatchListItemPayload{
				{
					Active:      true,
					Description: d.Get("description").(string),
					Expiration:  d.Get("expiration").(string),
					Value:       d.Get("value").(string),
				},
			},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceMatchListItemRead(ctx, d, m)

	return diags
}

func resourceMatchListItemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	matchList, err := c.Read(MatchListItems, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("description", matchList.(MatchListItemResponse).MatchListItem.Meta.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("expiration", matchList.(MatchListItemResponse).MatchListItem.Expiration)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("value", matchList.(MatchListItemResponse).MatchListItem.Value)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceMatchListItemUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("description", "expiration") {
		c := m.(*Client)

		err := c.Update(d.Id(), MatchListItemUpdateRequest{
			Fields: MatchListItemUpdatePayload{
				Active:      true,
				Description: d.Get("description").(string),
				Expiration:  d.Get("expiration").(string),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMatchListItemRead(ctx, d, m)
}

func resourceMatchListItemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), "match-lists-items") // This actually doesn't exist
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

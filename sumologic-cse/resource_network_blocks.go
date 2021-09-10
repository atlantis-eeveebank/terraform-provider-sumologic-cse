package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NetworkBlockResponse struct {
	NetworkBlock NetworkBlock `json:"data"`
}

type NetworkBlock struct {
	Id           string `json:"id"`
	AddressBlock string `json:"addressBlock"`
	Internal     bool   `json:"internal"`
	Label        string `json:"label"`
}

type NetworkBlockRequest struct {
	Fields PostNetworkPayload `json:"fields"`
}

type PostNetworkPayload struct {
	AddressBlock      string `json:"addressBlock"`
	Internal          bool   `json:"internal"`
	Label             string `json:"label"`
	SuppressesSignals bool   `json:"suppressesSignals"`
}

func resourceNetworkBlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkBlockCreate,
		ReadContext:   resourceNetworkBlockRead,
		UpdateContext: resourceNetworkBlockUpdate,
		DeleteContext: resourceNetworkBlockDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_block": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"internal": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"suppresses_signals": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceNetworkBlockCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(NetworkBlockRequest{
		Fields: PostNetworkPayload{
			AddressBlock:      d.Get("address_block").(string),
			Internal:          d.Get("internal").(bool),
			Label:             d.Get("label").(string),
			SuppressesSignals: d.Get("suppresses_signals").(bool),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceNetworkBlockRead(ctx, d, m)

	return diags
}

func resourceNetworkBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	cetD, err := c.Read("network-blocks", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("address_block", cetD.(NetworkBlockResponse).NetworkBlock.AddressBlock)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("internal", cetD.(NetworkBlockResponse).NetworkBlock.Internal)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("label", cetD.(NetworkBlockResponse).NetworkBlock.Label)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceNetworkBlockUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("address_block", "internal", "label", "suppresses_signals") {
		c := m.(*Client)

		err := c.Update(d.Id(), NetworkBlockRequest{
			Fields: PostNetworkPayload{
				AddressBlock:      d.Get("address_block").(string),
				Internal:          d.Get("internal").(bool),
				Label:             d.Get("label").(string),
				SuppressesSignals: d.Get("suppresses_signals").(bool),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNetworkBlockRead(ctx, d, m)
}

func resourceNetworkBlockDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), "network-blocks")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

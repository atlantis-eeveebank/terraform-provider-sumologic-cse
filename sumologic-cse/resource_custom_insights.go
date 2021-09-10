package sumologic_cse

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

type CustomInsightListResponse struct {
	Data struct {
		CustomInsights []CustomInsight `json:"objects"`
	} `json:"data"`
}

type CustomInsightResponse struct {
	CustomInsight CustomInsight`json:"data"`
}

type CustomInsight struct {
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Id          string    `json:"id"`
	LastUpdated time.Time `json:"lastUpdated"`
	Name        string    `json:"name"`
	Ordered     bool      `json:"ordered"`
	RuleIds     []string  `json:"ruleIds"`
	Severity    string    `json:"severity"`
	Tags        []string  `json:"tags"`
}

type CustomInsightRequest struct {
	Fields CustomInsightPayload `json:"fields"`
}

type CustomInsightPayload struct {
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Name        string   `json:"name"`
	Ordered     bool     `json:"ordered"`
	RuleIds     []string `json:"ruleIds"`
	Severity    string   `json:"severity"`
	Tags        []string `json:"tags"`
}

func resourceCustomInsight() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomInsightCreate,
		ReadContext:   resourceCustomInsightRead,
		UpdateContext: resourceCustomInsightUpdate,
		DeleteContext: resourceCustomInsightDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ordered": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"rule_ids": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"severity": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCustomInsightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	setRuleIds := d.Get("rule_ids").([]interface{})
	ruleIds := make([]string, len(setRuleIds))
	for _, ruleId := range setRuleIds {
		ruleIds = append(ruleIds, ruleId.(string))
	}

	setTags := d.Get("tags").([]interface{})
	tags := make([]string, len(setTags))
	for _, tag := range setTags {
		tags = append(tags, tag.(string))
	}

	id, err := c.Create(CustomInsightRequest{
		Fields: CustomInsightPayload{
			Description: d.Get("description").(string),
			Enabled:     d.Get("enabled").(bool),
			Name:        d.Get("name").(string),
			Ordered:     d.Get("ordered").(bool),
			RuleIds:     ruleIds,
			Severity:    d.Get("severity").(string),
			Tags:        tags,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceCustomInsightRead(ctx, d, m)

	return diags
}

func resourceCustomInsightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	result, err := c.ReadAll("custom-insights")
	if err != nil {
		return diag.FromErr(err)
	}

	var customInsight CustomInsight
	found := false
	for _, ci := range 	result.(CustomInsightListResponse).Data.CustomInsights {
		if ci.Id == d.Id() {
			customInsight = ci
			found = true
			break
		}
	}
	if !found {
		return diag.FromErr(errors.New(fmt.Sprintf("could not find custom insight with id: %s", d.Id())))
	}

	err = d.Set("description", customInsight.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("enabled", customInsight.Enabled)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", customInsight.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("ordered", customInsight.Ordered)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("rule_ids", customInsight.RuleIds)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("severity", customInsight.Severity)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("tags", customInsight.Tags)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCustomInsightUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("ordered", "description", "enabled", "name", "severity", "rule_ids", "tags") {
		c := m.(*Client)

		setRuleIds := d.Get("rule_ids").([]interface{})
		ruleIds := make([]string, len(setRuleIds))
		for _, ruleId := range setRuleIds {
			ruleIds = append(ruleIds, ruleId.(string))
		}

		setTags := d.Get("tags").([]interface{})
		tags := make([]string, len(setTags))
		for _, tag := range setTags {
			tags = append(tags, tag.(string))
		}

		err := c.Update(d.Id(), CustomInsightRequest{
			Fields: CustomInsightPayload{
				Description: d.Get("description").(string),
				Enabled:     d.Get("enabled").(bool),
				Name:        d.Get("name").(string),
				Ordered:     d.Get("ordered").(bool),
				RuleIds:     ruleIds,
				Severity:    d.Get("severity").(string),
				Tags:        tags,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCustomInsightRead(ctx, d, m)
}

func resourceCustomInsightDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), "custom-insights")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

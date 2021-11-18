package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ThresholdRuleRequest struct {
	Fields ThresholdRulePayload `json:"fields"`
}

type ThresholdRulePayload struct {
	AssetField          string               `json:"assetField"`
	Category            string               `json:"category"`
	CountDistinct       bool                 `json:"countDistinct"`
	CountField          string               `json:"countField"`
	Description         string               `json:"description"`
	Enabled             bool                 `json:"enabled"`
	EntitySelectors     []RuleEntitySelector `json:"entitySelectors"`
	Expression          string               `json:"expression"`
	GroupByFields       []string             `json:"groupByFields"`
	IsPrototype         bool                 `json:"isPrototype"`
	Limit               int                  `json:"limit"`
	Name                string               `json:"name"`
	Score               int                  `json:"score"`
	Stream              string               `json:"stream"`
	SummaryExpression   string               `json:"summaryExpression"`
	Tags                []string             `json:"tags"`
	TuningExpressionIds []string             `json:"tuningExpressionIds"`
	Version             int                  `json:"version"`
	WindowSize          string               `json:"windowSize"`
}

func resourceThresholdRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThresholdRuleCreate,
		ReadContext:   resourceThresholdRuleRead,
		UpdateContext: resourceThresholdRuleUpdate,
		DeleteContext: resourceRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"count_distinct": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_prototype": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"count_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"severity": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"summary_expression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"window_size": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_by_fields": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"entity_selectors": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"expression": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func thresholdRuleHasChanges(d resourceDiffer) bool {
	return d.HasChange("count_distinct") ||
		d.HasChange("count_field") ||
		d.HasChange("description") ||
		d.HasChange("entity_selectors") ||
		d.HasChange("expression") ||
		d.HasChange("group_by_fields") ||
		d.HasChange("limit") ||
		d.HasChange("name") ||
		d.HasChange("severity") ||
		d.HasChange("summary_expression") ||
		d.HasChange("tags") ||
		d.HasChange("version") ||
		d.HasChange("window_size") ||
		d.HasChange("is_prototype")
}

func resourceThresholdRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(ThresholdRuleRequest{
		Fields: ThresholdRulePayload{
			CountDistinct:     d.Get("count_distinct").(bool),
			CountField:        d.Get("count_field").(string),
			Description:       d.Get("description").(string),
			Enabled:           d.Get("enabled").(bool),
			EntitySelectors:   toEntitySelectorSlice(d.Get("entity_selectors")),
			Expression:        d.Get("expression").(string),
			GroupByFields:     toStringSlice(d.Get("group_by_fields")),
			IsPrototype:       d.Get("is_prototype").(bool),
			Limit:             d.Get("limit").(int),
			Name:              d.Get("name").(string),
			Score:             d.Get("severity").(int),
			Stream:            "record",
			SummaryExpression: d.Get("summary_expression").(string),
			Tags:              toStringSlice(d.Get("tags")),
			Version:           0,
			WindowSize:        d.Get("window_size").(string),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceThresholdRuleRead(ctx, d, m)

	return diags
}

func resourceThresholdRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	result, err := c.Read(Rules, d.Id()+"?expand=tuningExpressions")
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("enabled", result.(RuleResponse).Rule.Enabled)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("is_prototype", result.(RuleResponse).Rule.IsPrototype)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("count_distinct", result.(RuleResponse).Rule.CountDistinct)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("count_field", result.(RuleResponse).Rule.CountField)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("description", result.(RuleResponse).Rule.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("expression", result.(RuleResponse).Rule.Expression)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("limit", result.(RuleResponse).Rule.Limit)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", result.(RuleResponse).Rule.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("severity", result.(RuleResponse).Rule.Score)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("summary_expression", result.(RuleResponse).Rule.SummaryExpression)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("window_size", result.(RuleResponse).Rule.WindowSize)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("group_by_fields", result.(RuleResponse).Rule.GroupByFields)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("tags", result.(RuleResponse).Rule.Tags)
	if err != nil {
		return diag.FromErr(err)
	}

	es, err := flattenData(result.(RuleResponse).Rule.EntitySelectors)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("entity_selectors", es)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceThresholdRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if thresholdRuleHasChanges(d) {
		err := c.Update(d.Id(), ThresholdRuleRequest{
			Fields: ThresholdRulePayload{
				CountDistinct:     d.Get("count_distinct").(bool),
				CountField:        d.Get("count_field").(string),
				Description:       d.Get("description").(string),
				EntitySelectors:   toEntitySelectorSlice(d.Get("entity_selectors")),
				Expression:        d.Get("expression").(string),
				GroupByFields:     toStringSlice(d.Get("group_by_fields")),
				IsPrototype:       d.Get("is_prototype").(bool),
				Limit:             d.Get("limit").(int),
				Name:              d.Get("name").(string),
				Score:             d.Get("severity").(int),
				SummaryExpression: d.Get("summary_expression").(string),
				Tags:              toStringSlice(d.Get("tags")),
				WindowSize:        d.Get("window_size").(string),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("enabled") {
		err := c.Enabled(d.Id(), Rules, d.Get("enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceThresholdRuleRead(ctx, d, m)
}

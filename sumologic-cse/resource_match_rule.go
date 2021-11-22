package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TemplatedRuleRequest struct {
	Fields TemplatedRulePayload `json:"fields"`
}

type TemplatedRulePayload struct {
	AssetField            string               `json:"assetField"`
	Category              string               `json:"category"`
	DescriptionExpression string               `json:"descriptionExpression"`
	Enabled               bool                 `json:"enabled"`
	EntitySelectors       []RuleEntitySelector `json:"entitySelectors"`
	Expression            string               `json:"expression"`
	IsPrototype           bool                 `json:"isPrototype"`
	Name                  string               `json:"name"`
	NameExpression        string               `json:"nameExpression"`
	ScoreMapping          RuleScoreMapping     `json:"scoreMapping"`
	Stream                string               `json:"stream"`
	SummaryExpression     string               `json:"summaryExpression"`
	Tags                  []string             `json:"tags"`
	TuningExpressionIds   []string             `json:"tuningExpressionIds"`
}

func resourceTemplatedRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplatedRuleCreate,
		ReadContext:   resourceTemplatedRuleRead,
		UpdateContext: resourceTemplatedRuleUpdate,
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
			"is_prototype": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description_expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name_expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"summary_expression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"severity_mapping": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func templatedRuleHasChanges(d resourceDiffer) bool {
	return d.HasChange("description_expression") ||
		d.HasChange("entity_selectors") ||
		d.HasChange("expression") ||
		d.HasChange("is_prototype") ||
		d.HasChange("name") ||
		d.HasChange("name_expression") ||
		d.HasChange("severity_mapping") ||
		d.HasChange("summary_expression") ||
		d.HasChange("tags")
}

func templatedRulePayload(d *schema.ResourceData) TemplatedRuleRequest {
	return TemplatedRuleRequest{
		Fields: TemplatedRulePayload{
			DescriptionExpression: d.Get("description_expression").(string),
			Enabled:               d.Get("enabled").(bool),
			EntitySelectors:       toEntitySelectorSlice(d.Get("entity_selectors")),
			Expression:            d.Get("expression").(string),
			IsPrototype:           d.Get("is_prototype").(bool),
			Name:                  d.Get("name").(string),
			NameExpression:        d.Get("name_expression").(string),
			ScoreMapping:          toStructRuleScoreMapping(d.Get("severity_mapping")),
			Stream:                "record",
			SummaryExpression:     d.Get("summary_expression").(string),
			Tags:                  toStringSlice(d.Get("tags")),
		},
	}
}

func resourceTemplatedRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(thresholdRulePayload(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceTemplatedRuleRead(ctx, d, m)

	return diags
}

func resourceTemplatedRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	err = d.Set("description_expression", result.(RuleResponse).Rule.DescriptionExpression)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("expression", result.(RuleResponse).Rule.Expression)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", result.(RuleResponse).Rule.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name_expression", result.(RuleResponse).Rule.NameExpression)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("summary_expression", result.(RuleResponse).Rule.SummaryExpression)
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

	sm, err := flattenData(result.(RuleResponse).Rule.ScoreMapping)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("severity_mapping", sm)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceTemplatedRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if templatedRuleHasChanges(d) {
		err := c.Update(d.Id(), templatedRulePayload(d))
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

	return resourceTemplatedRuleRead(ctx, d, m)
}

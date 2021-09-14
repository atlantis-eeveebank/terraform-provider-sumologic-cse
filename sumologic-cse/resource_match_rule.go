package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMatchRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchRuleCreate,
		ReadContext:   resourceMatchRuleRead,
		UpdateContext: resourceMatchRuleUpdate,
		DeleteContext: resourceRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"asset_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_prototype": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"count_distinct": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"category": &schema.Schema{
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_jask_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"score": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"stream": &schema.Schema{
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
			"tuning_expression_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"entity_selector": &schema.Schema{
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

func matchRuleHasChanges(d resourceDiffer) bool {
	return d.HasChange("is_prototype") ||
		d.HasChange("category") ||
		d.HasChange("description") ||
		d.HasChange("asset_field") ||
		d.HasChange("expression") ||
		d.HasChange("name") ||
		d.HasChange("parent_jask_id") ||
		d.HasChange("count_distinct") ||
		d.HasChange("stream") ||
		d.HasChange("summary_expression") ||
		d.HasChange("tags") ||
		d.HasChange("tuning_expression_ids") ||
		d.HasChange("entity_selector") ||
		d.HasChange("score_mapping")
}

func resourceMatchRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(RuleRequest{
		Fields: RulePayload{
			AssetField: d.Get("asset_field").(string),
			Category: d.Get("category").(string),
			CountDistinct: d.Get("count_distinct").(bool),
			Description: d.Get("description").(string),
			Enabled: d.Get("enabled").(bool),
			EntitySelectors: d.Get("entity_selector").([]RuleEntitySelector),
			IsPrototype: d.Get("is_prototype").(bool),
			Expression: d.Get("expression").(string),
			Name: d.Get("name").(string),
			ParentJaskId: d.Get("parent_jask_id").(string),
			Score: d.Get("score").(int),
			Stream: d.Get("stream").(string),
			SummaryExpression: d.Get("summary_expression").(string),
			Tags: d.Get("tags").([]string),
			TuningExpressionIds: d.Get("tags").([]string),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceMatchRuleRead(ctx, d, m)

	return diags
}

func resourceMatchRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	result, err := c.Read("rules", d.Id()+"?expand=tuningExpressions")
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("asset_field", result.(RuleResponse).Rule.AssetField)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("count_distinct", result.(RuleResponse).Rule.CountDistinct)
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

	err = d.Set("category", result.(RuleResponse).Rule.Category)
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

	err = d.Set("name", result.(RuleResponse).Rule.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("parent_jask_id", result.(RuleResponse).Rule.ParentJaskId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("score", result.(RuleResponse).Rule.Score)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("stream", result.(RuleResponse).Rule.Stream)
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

	err = d.Set("tuning_expression_ids", result.(RuleResponse).Rule.TuningExpressionIds)
	if err != nil {
		return diag.FromErr(err)
	}

	es, err := flattenData(result.(RuleResponse).Rule.EntitySelectors)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("entity_selector", es)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceMatchRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if matchRuleHasChanges(d) {
		err := c.Update(d.Id(), RuleRequest{
			Fields: RulePayload{
				AssetField: d.Get("asset_field").(string),
				Category: d.Get("category").(string),
				CountDistinct: d.Get("count_distinct").(bool),
				Description: d.Get("description").(string),
				EntitySelectors: d.Get("entity_selector").([]RuleEntitySelector),
				IsPrototype: d.Get("is_prototype").(bool),
				Expression: d.Get("expression").(string),
				Name: d.Get("name").(string),
				ParentJaskId: d.Get("parent_jask_id").(string),
				Score: d.Get("score_mapping").(int),
				Stream: d.Get("stream").(string),
				SummaryExpression: d.Get("summary_expression").(string),
				Tags: d.Get("tags").([]string),
				TriggerExpression: d.Get("tuning_expression_ids").(string),
				TuningExpressionIds: d.Get("tags").([]string),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("enabled") {
		err := c.Enabled(d.Id(), "rules", d.Get("enabled").(bool))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMatchRuleRead(ctx, d, m)
}

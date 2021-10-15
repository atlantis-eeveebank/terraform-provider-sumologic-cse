package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AggregationRuleRequest struct {
	Fields AggregationRulePayload `json:"fields"`
}

type AggregationRulePayload struct {
	AggregationFunctions  []RuleAggregationFunction `json:"aggregationFunctions"`
	AssetField            string                    `json:"assetField"`
	Category              string                    `json:"category"`
	DescriptionExpression string                    `json:"descriptionExpression"`
	Enabled               bool                      `json:"enabled"`
	EntitySelectors       []RuleEntitySelector      `json:"entitySelectors"`
	GroupByAsset          bool                      `json:"groupByAsset"`
	GroupByFields         []string                  `json:"groupByFields"`
	IsPrototype           bool                      `json:"isPrototype"`
	MatchExpression       string                    `json:"matchExpression"`
	Name                  string                    `json:"name"`
	NameExpression        string                    `json:"nameExpression"`
	ScoreMapping          RuleScoreMapping          `json:"scoreMapping"`
	Stream                string                    `json:"stream"`
	SummaryExpression     string                    `json:"summaryExpression"`
	Tags                  []string                  `json:"tags"`
	TriggerExpression     string                    `json:"triggerExpression"`
	TuningExpressionIds   []string                  `json:"tuningExpressionIds"`
	WindowSize            string                    `json:"windowSize"`
}

type RuleAggregationFunction struct {
	Arguments []string `json:"arguments"`
	Function  string   `json:"function"`
	Name      string   `json:"name"`
}

func resourceAggregationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAggregationRuleCreate,
		ReadContext:   resourceAggregationRuleRead,
		UpdateContext: resourceAggregationRuleUpdate,
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
			"group_by_asset": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_prototype": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description_expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"match_expression": &schema.Schema{
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
			"parent_jask_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"stream": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"summary_expression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"trigger_expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"window_size": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_by_fields": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tuning_expression_ids": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"aggregation_function": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"arguments": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
			"score_mapping": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"field": &schema.Schema{
							Type:     schema.TypeString,
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

func aggregationRuleHasChanges(d resourceDiffer) bool {
	return d.HasChange("aggregation_function") ||
		d.HasChange("asset_field") ||
		d.HasChange("category") ||
		d.HasChange("description_expression") ||
		d.HasChange("entity_selector") ||
		d.HasChange("group_by_fields") ||
		d.HasChange("is_prototype") ||
		d.HasChange("match_expression") ||
		d.HasChange("name") ||
		d.HasChange("name_expression") ||
		d.HasChange("parent_jask_id") ||
		d.HasChange("score_mapping") ||
		d.HasChange("stream") ||
		d.HasChange("summary_expression") ||
		d.HasChange("tags") ||
		d.HasChange("trigger_expression") ||
		d.HasChange("tuning_expression_ids") ||
		d.HasChange("window_size") ||
		d.HasChange("group_by_asset")
}

func resourceAggregationRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(AggregationRuleRequest{
		Fields: AggregationRulePayload{
			AggregationFunctions:  toAggregationFunctionSlice(d.Get("aggregation_function")),
			AssetField:            d.Get("asset_field").(string),
			Category:              d.Get("category").(string),
			DescriptionExpression: d.Get("description_expression").(string),
			Enabled:               d.Get("enabled").(bool),
			EntitySelectors:       toEntitySelectorSlice(d.Get("entity_selector")),
			GroupByAsset:          d.Get("group_by_asset").(bool),
			GroupByFields:         toStringSlice(d.Get("group_by_fields")),
			IsPrototype:           d.Get("is_prototype").(bool),
			MatchExpression:       d.Get("match_expression").(string),
			Name:                  d.Get("name").(string),
			NameExpression:        d.Get("name_expression").(string),
			ScoreMapping:          toStructRuleScoreMapping(d.Get("score_mapping")),
			Stream:                d.Get("stream").(string),
			SummaryExpression:     d.Get("summary_expression").(string),
			Tags:                  toStringSlice(d.Get("tags")),
			TriggerExpression:     d.Get("trigger_expression").(string),
			TuningExpressionIds:   toStringSlice(d.Get("tuning_expression_ids")),
			WindowSize:            d.Get("window_size").(string),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceAggregationRuleRead(ctx, d, m)

	return diags
}

func resourceAggregationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	err = d.Set("enabled", result.(RuleResponse).Rule.Enabled)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("group_by_asset", result.(RuleResponse).Rule.GroupByAsset)
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

	err = d.Set("description_expression", result.(RuleResponse).Rule.DescriptionExpression)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("match_expression", result.(RuleResponse).Rule.MatchExpression)
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

	err = d.Set("parent_jask_id", result.(RuleResponse).Rule.ParentJaskId)
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

	err = d.Set("trigger_expression", result.(RuleResponse).Rule.TriggerExpression)
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

	err = d.Set("tuning_expression_ids", result.(RuleResponse).Rule.TuningExpressionIds)
	if err != nil {
		return diag.FromErr(err)
	}

	af, err := flattenData(result.(RuleResponse).Rule.AggregationFunctions)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("aggregation_function", af)
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

	sm, err := flattenData(result.(RuleResponse).Rule.ScoreMapping)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("score_mapping", sm)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAggregationRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if aggregationRuleHasChanges(d) {
		err := c.Update(d.Id(), AggregationRuleRequest{
			Fields: AggregationRulePayload{
				AggregationFunctions:  toAggregationFunctionSlice(d.Get("aggregation_function")),
				AssetField:            d.Get("asset_field").(string),
				Category:              d.Get("category").(string),
				DescriptionExpression: d.Get("description_expression").(string),
				EntitySelectors:       toEntitySelectorSlice(d.Get("entity_selector")),
				GroupByAsset:          d.Get("group_by_asset").(bool),
				GroupByFields:         toStringSlice(d.Get("group_by_fields")),
				IsPrototype:           d.Get("is_prototype").(bool),
				MatchExpression:       d.Get("match_expression").(string),
				Name:                  d.Get("name").(string),
				NameExpression:        d.Get("name_expression").(string),
				ScoreMapping:          toStructRuleScoreMapping(d.Get("score_mapping")),
				Stream:                d.Get("stream").(string),
				SummaryExpression:     d.Get("summary_expression").(string),
				Tags:                  toStringSlice(d.Get("tags")),
				TriggerExpression:     d.Get("trigger_expression").(string),
				TuningExpressionIds:   toStringSlice(d.Get("tuning_expression_ids")),
				WindowSize:            d.Get("window_size").(string),
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

	return resourceAggregationRuleRead(ctx, d, m)
}

package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RuleResponse struct {
	Rule Rule `json:"data"`
}

type Rule struct {
	AggregationFunctions  []RuleAggregationFunction `json:"aggregationFunctions"`
	AssetField            string                    `json:"assetField"`
	Category              string                    `json:"category"`
	ContentType           string                    `json:"contentType"`
	CountDistinct         bool                      `json:"countDistinct"`
	CountField            string                    `json:"countField"`
	CreatedBy             string                    `json:"createdBy"`
	Deleted               bool                      `json:"deleted"`
	Description           string                    `json:"description"`
	DescriptionExpression string                    `json:"descriptionExpression"`
	Enabled               bool                      `json:"enabled"`
	EntitySelectors       []RuleEntitySelector      `json:"entitySelectors"`
	Expression            string                    `json:"expression"`
	GroupByAsset          bool                      `json:"groupByAsset"`
	GroupByFields         []string                  `json:"groupByFields"`
	HasOverride           bool                      `json:"hasOverride"`
	Id                    string                    `json:"id"`
	IsPrototype           bool                      `json:"isPrototype"`
	LastUpdatedBy         string                    `json:"lastUpdatedBy"`
	Limit                 int                       `json:"limit"`
	MatchExpression       string                    `json:"matchExpression"`
	Name                  string                    `json:"name"`
	NameExpression        string                    `json:"nameExpression"`
	ParentJaskId          string                    `json:"parentJaskId"`
	RuleId                int                       `json:"ruleId"`
	RuleSource            string                    `json:"ruleSource"`
	RuleType              string                    `json:"ruleType"`
	Score                 int                       `json:"score"`
	ScoreMapping          RuleScoreMapping          `json:"scoreMapping"`
	SignalCount07D        int                       `json:"signalCount07d"`
	SignalCount24H        int                       `json:"signalCount24h"`
	Status                RuleStatus                `json:"status"`
	Stream                string                    `json:"stream"`
	SummaryExpression     string                    `json:"summaryExpression"`
	Tags                  []string                  `json:"tags"`
	TriggerExpression     string                    `json:"triggerExpression"`
	TuningExpressionIds   []string                  `json:"tuningExpressionIds"`
	Version               int                       `json:"version"`
	WindowSize            string                    `json:"windowSizeName"`
}

type RuleEntitySelector struct {
	EntityType string `json:"entityType"`
	Expression string `json:"expression"`
}

type RuleScoreMapping struct {
	Default int    `json:"default"`
	Type    string `json:"type"`
}

type RuleStatus struct {
	Message interface{} `json:"message"`
	Status  string      `json:"status"`
}

func toAggregationFunctionSlice(data interface{}) []RuleAggregationFunction {
	datum := data.([]interface{})

	aggregationFunctions := make([]RuleAggregationFunction, len(datum))
	for i, v := range datum {
		af := v.(map[string]interface{})

		aggregationFunctions[i] = RuleAggregationFunction{
			Arguments: toStringSlice(af["arguments"]),
			Function:  af["function"].(string),
			Name:      af["name"].(string),
		}
	}

	return aggregationFunctions
}

func toEntitySelectorSlice(data interface{}) []RuleEntitySelector {
	datum := data.([]interface{})

	entitySelectors := make([]RuleEntitySelector, len(datum))
	for i, v := range datum {
		es := v.(map[string]interface{})

		entitySelectors[i] = RuleEntitySelector{
			EntityType: es["entity_type"].(string),
			Expression: es["expression"].(string),
		}
	}

	return entitySelectors
}

func toStructRuleScoreMapping(data interface{}) RuleScoreMapping {
	datum := data.([]interface{})[0].(map[string]interface{})

	return RuleScoreMapping{
		Default: datum["default"].(int),
		Type:    datum["type"].(string),
	}
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), Rules)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

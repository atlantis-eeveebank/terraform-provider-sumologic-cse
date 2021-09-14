package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RuleResponse struct {
	Rule Rule `json:"data"`
}

type RuleRequest struct {
	Fields RulePayload `json:"fields"`
}

type Rule struct {
	Deleted               bool     `json:"deleted"`
	Enabled               bool     `json:"enabled"`
	HasOverride           bool     `json:"hasOverride"`
	IsPrototype           bool     `json:"isPrototype"`
	GroupByAsset          bool     `json:"groupByAsset"`
	CountDistinct         bool     `json:"countDistinct"`
	AssetField            string   `json:"assetField"`
	Category              string   `json:"category"`
	ContentType           string   `json:"contentType"`
	CountField            string   `json:"countField"`
	CreatedBy             string   `json:"createdBy"`
	Description           string   `json:"description"`
	DescriptionExpression string   `json:"descriptionExpression"`
	Expression            string   `json:"expression"`
	Id                    string   `json:"id"`
	LastUpdatedBy         string   `json:"lastUpdatedBy"`
	MatchExpression       string   `json:"matchExpression"`
	Name                  string   `json:"name"`
	NameExpression        string   `json:"nameExpression"`
	ParentJaskId          string   `json:"parentJaskId"`
	RuleSource            string   `json:"ruleSource"`
	RuleType              string   `json:"ruleType"`
	Stream                string   `json:"stream"`
	SummaryExpression     string   `json:"summaryExpression"`
	TriggerExpression     string   `json:"triggerExpression"`
	RuleId                int      `json:"ruleId"`
	Limit                 int      `json:"limit"`
	Score                 int      `json:"score"`
	SignalCount07D        int      `json:"signalCount07d"`
	SignalCount24H        int      `json:"signalCount24h"`
	WindowSize            int      `json:"windowSize"`
	Version               int      `json:"version"`
	GroupByFields         []string `json:"groupByFields"`
	Tags                  []string `json:"tags"`
	TuningExpressionIds   []string `json:"tuningExpressionIds"`

	AggregationFunctions []RuleAggregationFunction `json:"aggregationFunctions"`
	EntitySelectors      []RuleEntitySelector      `json:"entitySelectors"`
	ScoreMapping         RuleScoreMapping          `json:"scoreMapping"`
	Status               RuleStatus                `json:"status"`
}

type RulePayload struct {
	Enabled               bool     `json:"enabled"`
	GroupByAsset          bool     `json:"groupByAsset"`
	CountDistinct         bool     `json:"countDistinct"`
	IsPrototype           bool     `json:"isPrototype"`
	AssetField            string   `json:"assetField"`
	Category              string   `json:"category"`
	CountField            string   `json:"countField"`
	Description           string   `json:"description"`
	DescriptionExpression string   `json:"descriptionExpression"`
	MatchExpression       string   `json:"matchExpression"`
	Expression            string   `json:"expression"`
	Name                  string   `json:"name"`
	NameExpression        string   `json:"nameExpression"`
	ParentJaskId          string   `json:"parentJaskId"`
	Stream                string   `json:"stream"`
	SummaryExpression     string   `json:"summaryExpression"`
	TriggerExpression     string   `json:"triggerExpression"`
	WindowSize            string   `json:"windowSize"`
	Limit                 int      `json:"limit"`
	Score                 int      `json:"score"`
	Version               int      `json:"version"`
	GroupByFields         []string `json:"groupByFields"`
	Tags                  []string `json:"tags"`
	TuningExpressionIds   []string `json:"tuningExpressionIds"`

	AggregationFunctions []RuleAggregationFunction `json:"aggregationFunctions"`
	EntitySelectors      []RuleEntitySelector      `json:"entitySelectors"`
	ScoreMapping         RuleScoreMapping          `json:"scoreMapping"`
}

type RuleAggregationFunction struct {
	Arguments []string `json:"arguments"`
	Function  string   `json:"function"`
	Name      string   `json:"name"`
}

type RuleEntitySelector struct {
	EntityType string `json:"entityType"`
	Expression string `json:"expression"`
}

type RuleScoreMapping struct {
	Default int    `json:"default"`
	Field   string `json:"field"`
	Type    string `json:"type"`
}

type RuleStatus struct {
	Message interface{} `json:"message"`
	Status  string      `json:"status"`
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), "rules")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type resourceDiffer interface {
	HasChange(string) bool
}

type LogMappingResponse struct {
	LogMapping LogMapping `json:"data"`
}

type LogMapping struct {
	Id                 string                      `json:"id"`
	Enabled            bool                        `json:"enabled"`
	Name               string                      `json:"name"`
	ProductGuid        string                      `json:"productGuid"`
	RecordType         string                      `json:"recordType"`
	RelatesEntities    bool                        `json:"relatesEntities"`
	Source             string                      `json:"source"`
	UnstructuredFields string                      `json:"unstructuredFields"`
	SkippedValues      []string                    `json:"skippedValues"`
	Fields             []LogMappingField           `json:"fields"`
	StructuredInputs   []LogMappingStructuredInput `json:"structuredInputs"`
	Input              LogMappingInput             `json:"input"`
}

type LogMappingField struct {
	AlternateValues  string `json:"alternateValues"`
	CaseInsensitive  string `json:"caseInsensitive"`
	DefaultValue     string `json:"defaultValue"`
	FieldJoin        string `json:"fieldJoin"`
	Format           string `json:"format"`
	FormatParameters string `json:"formatParameters"`
	JoinDelimiter    string `json:"joinDelimiter"`
	Lookup           string `json:"lookup"`
	Name             string `json:"name"`
	SkippedValues    string `json:"skippedValues"`
	SplitDelimiter   string `json:"splitDelimiter"`
	SplitIndex       string `json:"splitIndex"`
	TimeZone         string `json:"timeZone"`
	Value            string `json:"value"`
	ValueType        string `json:"valueType"`
}

type LogMappingInput struct {
	Alternatives   []string `json:"alternatives"`
	EventIdPattern string   `json:"event_id_pattern"`
	LogFormat      string   `json:"log_format"`
	Product        string   `json:"product"`
	Vendor         string   `json:"vendor"`
}

type LogMappingStructuredInput struct {
	EventIdPattern string `json:"eventIdPattern"`
	LogFormat      string `json:"logFormat"`
	Product        string `json:"product"`
	Vendor         string `json:"vendor"`
}

type LogMappingRequest struct {
	Fields PostLogMappingPayload `json:"fields"`
}

type PostLogMappingPayload struct {
	Id                 string                      `json:"id"`
	Enabled            bool                        `json:"enabled"`
	Name               string                      `json:"name"`
	ProductGuid        string                      `json:"productGuid"`
	RecordType         string                      `json:"recordType"`
	RelatesEntities    bool                        `json:"relatesEntities"`
	Source             string                      `json:"source"`
	UnstructuredFields string                      `json:"unstructuredFields"`
	SkippedValues      []string                    `json:"skippedValues"`
	Fields             []LogMappingField           `json:"fields"`
	StructuredInputs   []LogMappingStructuredInput `json:"structuredInputs"`
	Input              LogMappingInput             `json:"input"`
}

func resourceLogMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogMappingCreate,
		ReadContext:   resourceLogMappingRead,
		UpdateContext: resourceLogMappingUpdate,
		DeleteContext: resourceLogMappingDelete,
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"product_guid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"record_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"relates_entities": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"skipped_values": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"input": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alternatives": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"event_id_pattern": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"log_format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"product": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"vendor": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"structured_input": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id_pattern": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"log_format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"product": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"vendor": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"field": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func hasConfigChanges(d resourceDiffer) bool {
	return d.HasChange("enabled") ||
		d.HasChange("name") ||
		d.HasChange("product_guid") ||
		d.HasChange("record_type") ||
		d.HasChange("relates_entities") ||
		d.HasChange("source") ||
		d.HasChange("unstructured_field") ||
		d.HasChange("skipped_values") ||
		d.HasChange("structured_input") ||
		d.HasChange("input")
}

func resourceLogMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(LogMappingRequest{
		Fields: PostLogMappingPayload{
			Enabled:          d.Get("enabled").(bool),
			Name:             d.Get("name").(string),
			ProductGuid:      d.Get("product_guid").(string),
			RecordType:       d.Get("record_type").(string),
			RelatesEntities:  d.Get("relates_entities").(bool),
			Source:           d.Get("source").(string),
			SkippedValues:    d.Get("skipped_values").([]string),
			Input:            d.Get("input").(LogMappingInput),
			StructuredInputs: d.Get("structured_input").([]LogMappingStructuredInput),
			Fields:           d.Get("field").([]LogMappingField),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	resourceLogMappingRead(ctx, d, m)

	return diags
}

func resourceLogMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	lmD, err := c.Read("log-mappings", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("enabled", lmD.(LogMappingResponse).LogMapping.Enabled)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", lmD.(LogMappingResponse).LogMapping.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("product_guid", lmD.(LogMappingResponse).LogMapping.ProductGuid)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("record_type", lmD.(LogMappingResponse).LogMapping.RecordType)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("relates_entities", lmD.(LogMappingResponse).LogMapping.RelatesEntities)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("source", lmD.(LogMappingResponse).LogMapping.Source)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("skipped_values", lmD.(LogMappingResponse).LogMapping.SkippedValues)
	if err != nil {
		return diag.FromErr(err)
	}

	input, err := flattenData(lmD.(LogMappingResponse).LogMapping.Input)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("input", input)
	if err != nil {
		return diag.FromErr(err)
	}

	structuredInputs, err := flattenData(lmD.(LogMappingResponse).LogMapping.StructuredInputs)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("structured_input", structuredInputs)
	if err != nil {
		return diag.FromErr(err)
	}

	fields, err := flattenData(lmD.(LogMappingResponse).LogMapping.Fields)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("field", fields)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceLogMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if hasConfigChanges(d) {
		c := m.(*Client)

		err := c.Update(d.Id(), LogMappingRequest{
			Fields: PostLogMappingPayload{
				Enabled:          d.Get("enabled").(bool),
				Name:             d.Get("name").(string),
				ProductGuid:      d.Get("product_guid").(string),
				RecordType:       d.Get("record_type").(string),
				RelatesEntities:  d.Get("relates_entities").(bool),
				Source:           d.Get("source").(string),
				SkippedValues:    d.Get("skipped_values").([]string),
				Input:            d.Get("input").(LogMappingInput),
				StructuredInputs: d.Get("structured_input").([]LogMappingStructuredInput),
				Fields:           d.Get("field").([]LogMappingField),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLogMappingRead(ctx, d, m)
}

func resourceLogMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.Delete(d.Id(), "log-mappings")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

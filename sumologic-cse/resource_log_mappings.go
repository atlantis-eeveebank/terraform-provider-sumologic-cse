package sumologic_cse

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type resourceDiffer interface {
	HasChange(string) bool
}

type LogMappingSearchResponse struct {
	Data struct {
		Objects []struct {
			Enabled          bool                        `json:"enabled"`
			Fields           []LogMappingField           `json:"fields"`
			Id               string                      `json:"id"`
			Name             string                      `json:"name"`
			ProductGuid      string                      `json:"productGuid"`
			RecordType       string                      `json:"recordType"`
			RelatesEntities  bool                        `json:"relatesEntities"`
			SkippedValues    []string                    `json:"skippedValues"`
			Source           string                      `json:"source"`
			StructuredInputs []LogMappingStructuredInput `json:"structuredInputs"`
		} `json:"objects"`
		Total int `json:"total"`
	} `json:"data"`
}

type LogMappingResponse struct {
	LogMapping LogMapping `json:"data"`
}

type LogMapping struct {
	Enabled            bool                        `json:"enabled"`
	Id                 string                      `json:"id"`
	Fields             []LogMappingField           `json:"fields"`
	Name               string                      `json:"name"`
	ProductGuid        string                      `json:"productGuid"`
	RecordType         string                      `json:"recordType"`
	RelatesEntities    bool                        `json:"relatesEntities"`
	SkippedValues      []string                    `json:"skippedValues"`
	Source             string                      `json:"source"`
	StructuredInputs   []LogMappingStructuredInput `json:"structuredInputs"`
	UnstructuredFields LogMappingUnstructuredInput `json:"unstructuredFields"`
}

type LogMappingField struct {
	AlternateValues  []string                `json:"alternateValues,omitempty"`
	CaseInsensitive  bool                    `json:"caseInsensitive,omitempty"`
	DefaultValue     string                  `json:"defaultValue,omitempty"`
	FieldJoin        []string                `json:"fieldJoin,omitempty"`
	Format           string                  `json:"format,omitempty"`
	FormatParameters []string                `json:"formatParameters,omitempty"`
	JoinDelimiter    string                  `json:"joinDelimiter,omitempty"`
	Lookup           []LogMappingFieldLookup `json:"lookup,omitempty"`
	Name             string                  `json:"name"`
	SkippedValues    []string                `json:"skippedValues,omitempty"`
	SplitDelimiter   string                  `json:"splitDelimiter,omitempty"`
	SplitIndex       string                  `json:"splitIndex,omitempty"`
	TimeZone         string                  `json:"timeZone,omitempty"`
	Value            string                  `json:"value,omitempty"`
	ValueType        string                  `json:"valueType,omitempty"`
}

type LogMappingFieldLookup struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LogMappingStructuredInput struct {
	EventIdPattern string `json:"eventIdPattern"`
	LogFormat      string `json:"logFormat"`
	Product        string `json:"product"`
	Vendor         string `json:"vendor"`
}

type LogMappingUnstructuredInput struct {
	PatternNames []string `json:"patternNames"`
}

type LogMappingRequest struct {
	Fields PostLogMappingPayload `json:"fields"`
}

type PostLogMappingPayload struct {
	Fields             []LogMappingField           `json:"fields"`
	Name               string                      `json:"name"`
	Enabled            bool                        `json:"enabled,omitempty"`
	ParentId           string                      `json:"parentId,omitempty"`
	ProductGuid        string                      `json:"productGuid,omitempty"`
	RecordType         string                      `json:"recordType"`
	RelatesEntities    bool                        `json:"relatesEntities,omitempty"`
	SkippedValues      []string                    `json:"skippedValues,omitempty"`
	StructuredInputs   []LogMappingStructuredInput `json:"structuredInputs,omitempty"`
	UnstructuredFields LogMappingUnstructuredInput `json:"unstructuredFields,omitempty"`
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
			"parent_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"skipped_values": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"fields": &schema.Schema{
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
						"alternate_values": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"case_insensitive": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"default_value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"lookup": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func logMappingHasChanges(d resourceDiffer) bool {
	return d.HasChange("enabled") ||
		d.HasChange("name") ||
		d.HasChange("product_guid") ||
		d.HasChange("record_type") ||
		d.HasChange("relates_entities") ||
		d.HasChange("fields") ||
		d.HasChange("skipped_values") ||
		d.HasChange("structured_input")
}

func resourceLogMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	id, err := c.Create(LogMappingRequest{
		Fields: PostLogMappingPayload{
			Enabled:          d.Get("enabled").(bool),
			Fields:           expandLogMappingField(d),
			Name:             d.Get("name").(string),
			ParentId:         d.Get("parent_id").(string),
			ProductGuid:      d.Get("product_guid").(string),
			RecordType:       d.Get("record_type").(string),
			RelatesEntities:  d.Get("relates_entities").(bool),
			SkippedValues:    toStringSlice(d.Get("skipped_values")),
			StructuredInputs: expandStructuredInputs(d),
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

	lmD, err := c.Read(LogMappings, d.Id())
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

	err = d.Set("skipped_values", lmD.(LogMappingResponse).LogMapping.SkippedValues)
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
	err = d.Set("fields", fields)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceLogMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if logMappingHasChanges(d) {
		c := m.(*Client)

		err := c.Update(d.Id(), LogMappingRequest{
			Fields: PostLogMappingPayload{
				Enabled:          d.Get("enabled").(bool),
				Fields:           expandLogMappingField(d),
				Name:             d.Get("name").(string),
				ParentId:         d.Get("parent_id").(string),
				ProductGuid:      d.Get("product_guid").(string),
				RecordType:       d.Get("record_type").(string),
				RelatesEntities:  d.Get("relates_entities").(bool),
				SkippedValues:    toStringSlice(d.Get("skipped_values")),
				StructuredInputs: expandStructuredInputs(d),
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

	err := c.Delete(d.Id(), LogMappings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func expandLogMappingField(d *schema.ResourceData) []LogMappingField {
	f := d.Get("fields").([]interface{})
	var fields []LogMappingField

	for _, field := range f {
		fieldMap := field.(map[string]interface{})

		lus := fieldMap["lookup"].([]interface{})
		var lookups []LogMappingFieldLookup
		for _, lu := range lus {
			lookupMap := lu.(map[string]interface{})
			lookups = append(lookups, LogMappingFieldLookup{
				Key:   lookupMap["key"].(string),
				Value: lookupMap["value"].(string),
			})
		}

		fields = append(fields, LogMappingField{
			AlternateValues: toStringSlice(fieldMap["alternate_values"]),
			CaseInsensitive: fieldMap["case_insensitive"].(bool),
			DefaultValue:    fieldMap["default_value"].(string),
			Format:          fieldMap["format"].(string),
			Lookup:          lookups,
			Name:            fieldMap["name"].(string),
			Value:           fieldMap["value"].(string),
			ValueType:       fieldMap["value_type"].(string),
		})
	}

	return fields
}

func expandStructuredInputs(d *schema.ResourceData) []LogMappingStructuredInput {
	si := d.Get("structured_input").([]interface{})
	var structuredInputs []LogMappingStructuredInput

	for _, structuredInput := range si {
		structuredInputMap := structuredInput.(map[string]interface{})

		structuredInputs = append(structuredInputs, LogMappingStructuredInput{
			EventIdPattern: structuredInputMap["event_id_pattern"].(string),
			LogFormat:      structuredInputMap["log_format"].(string),
			Product:        structuredInputMap["product"].(string),
			Vendor:         structuredInputMap["vendor"].(string),
		})
	}

	return structuredInputs
}

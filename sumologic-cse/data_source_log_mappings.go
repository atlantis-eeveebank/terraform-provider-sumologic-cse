package sumologic_cse

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLogMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogMappingRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"product_guid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
				Computed: true,
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
				Computed: true,
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

func dataSourceLogMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*Client)

	lmD, err := c.Search(LogMappings, fmt.Sprintf("`name:\"%s\"`", d.Get("name")))
	if err != nil {
		return diag.FromErr(err)
	}

	if lmD.(LogMappingSearchResponse).Data.Total != 1 {
		diag.FromErr(errors.New("expected only 1 object to return when searching log mappings by name"))
	}

	logMapper := lmD.(LogMappingSearchResponse).Data.Objects[0]
	
	d.SetId(logMapper.Id)

	err = d.Set("enabled", logMapper.Enabled)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", logMapper.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("product_guid", logMapper.ProductGuid)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("record_type", logMapper.RecordType)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("relates_entities", logMapper.RelatesEntities)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("skipped_values", logMapper.SkippedValues)
	if err != nil {
		return diag.FromErr(err)
	}

	structuredInputs, err := flattenData(logMapper.StructuredInputs)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("structured_input", structuredInputs)
	if err != nil {
		return diag.FromErr(err)
	}

	fields, err := flattenData(logMapper.Fields)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("fields", fields)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

package main

import (
	"context"
	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/samber/lo"
	"strconv"
)

func resourceRedashVisualization() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashVisualizationRead,
		CreateContext: resourceRedashVisualizationCreate,
		UpdateContext: resourceRedashVisualizationUpdate,
		DeleteContext: resourceRedashVisualizationDelete,
		Schema: map[string]*schema.Schema{
			// Base Data
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			// References
			"query_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// Options (By Type)
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"items_per_page": {
							Type:     schema.TypeString,
							Required: true,
						},
						"columns": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// General
									"visible": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"title": {
										Type:     schema.TypeString,
										Required: true,
									},
									// Type
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"display_as": {
										Type:     schema.TypeString,
										Required: true,
									},
									"align_content": {
										Type:     schema.TypeString,
										Required: true,
									},
									"allow_search": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"order": {
										Type:     schema.TypeInt,
										Required: true,
									},
									// Text
									"allow_html": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"highlight_links": {
										Type:     schema.TypeBool,
										Required: true,
									},
									// Number
									"number_format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// Date/Time
									"date_time_format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									// Boolean
									"boolean_values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									// Link
									"link_url_template": {
										Type:     schema.TypeString,
										Required: true,
									},
									"link_text_template": {
										Type:     schema.TypeString,
										Required: true,
									},
									"link_open_in_new_tab": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"link_title_template": {
										Type:     schema.TypeString,
										Required: true,
									},
									// Image
									"image_url_template": {
										Type:     schema.TypeString,
										Required: true,
									},
									"image_width": {
										Type:     schema.TypeString,
										Required: true,
									},
									"image_height": {
										Type:     schema.TypeString,
										Required: true,
									},
									"image_title_template": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"chart_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// General
						"global_series_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"column_mapping": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"column": {
										Type:     schema.TypeString,
										Required: true,
									},
									"axis": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"legend": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"series": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"stacking": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"missing_values_as_zero": {
							Type:     schema.TypeBool,
							Required: true,
						},
						// X-Axis
						"x_axis": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"labels": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"sort_x": {
							Type:     schema.TypeBool,
							Required: true,
						},
						// Y-Axis
						"y_axis": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"opposite": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						// Series
						"series_options": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"z_index": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"index": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"y_axis": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						// Data Labels
						"show_data_labels": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"number_format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"percent_format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"date_time_format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"text_format": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceRedashVisualizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	queryID := d.Get("query_id").(int)
	visualizationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	visualization, err := c.GetVisualization(queryID, visualizationID)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", visualization.Name)
	_ = d.Set("type", visualization.Type)
	_ = d.Set("description", visualization.Description)

	return diags
}

func resourceRedashVisualizationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	var vOptions redash.VisualizationOptions
	vType := d.Get("type").(string)
	switch vType {
	case "CHART":
		chartOptions := d.Get("chart_options").([]interface{})[0].(map[string]interface{})
		chartLegend := chartOptions["legend"].(map[string]interface{})
		chartSeries := chartOptions["series"].(map[string]interface{})
		chartXAxis := chartOptions["x_axis"].([]interface{})[0].(map[string]interface{})
		chartYAxis := chartOptions["y_axis"].([]interface{})
		chartSeriesOptions := chartOptions["series_options"].([]interface{})
		vOptions = redash.VisualizationOptions{
			// General
			GlobalSeriesType: chartOptions["globalSeriesType"].(string),
			ColumnMapping:    chartOptions["columnMapping"].(map[string]string),
			Legend: redash.VisualizationLegendOptions{
				Enabled:   chartLegend["enabled"].(bool),
				Placement: chartLegend["placement"].(string),
			},
			Series: redash.Series{
				Stacking: chartSeries["stacking"].(string),
			},
			MissingValuesAsZero: chartOptions["missingValuesAsZero"].(bool),
			// CHART TYPE - X-Axis
			XAxis: redash.VisualizationAxisOptions{
				Type:     chartXAxis["type"].(string),
				Opposite: chartXAxis["opposite"].(bool),
				Labels: redash.VisualizationLabelOptions{
					Enabled: chartXAxis["labels"].(map[string]interface{})["enabled"].(bool),
				},
			},
			SortX: chartOptions["sortX"].(bool),
			// CHART TYPE - Y-Axis
			YAxis: lo.Map(chartYAxis, func(item interface{}, _ int) redash.VisualizationAxisOptions {
				yAxis := item.(map[string]interface{})

				return redash.VisualizationAxisOptions{
					Type:     yAxis["type"].(string),
					Opposite: yAxis["opposite"].(bool),
					Labels: redash.VisualizationLabelOptions{
						Enabled: yAxis["labels"].(map[string]interface{})["enabled"].(bool),
					},
				}
			}),
			// CHART TYPE - Series
			SeriesOptions: lo.Associate(chartSeriesOptions, func(value interface{}) (string, redash.SeriesOptions) {
				seriesOption := value.(map[string]interface{})

				return seriesOption["name"].(string), redash.SeriesOptions{
					ZIndex: seriesOption["z_index"].(int),
					Index:  seriesOption["index"].(int),
					Type:   seriesOption["type"].(string),
					YAxis:  seriesOption["y_axis"].(int),
				}
			}),
			// CHART TYPE - Colors
			// CHART TYPE - Data Labels
			ShowDataLabels: chartOptions["showDataLabels"].(bool),
			NumberFormat:   chartOptions["numberFormat"].(string),
			PercentFormat:  chartOptions["percentFormat"].(string),
			DateTimeFormat: chartOptions["dateTimeFormat"].(string),
			TextFormat:     chartOptions["textFormat"].(string),
		}
		break
	case "TABLE":
		tableOptions := d.Get("table_options").([]interface{})[0].(map[string]interface{})
		tableColumns := tableOptions["columns"].([]interface{})
		vOptions = redash.VisualizationOptions{
			ItemsPerPage: tableOptions["items_per_page"].(int),
			Columns: lo.Map(tableColumns, func(item interface{}, _ int) redash.VisualizationColumnOptions {
				column := item.(map[string]interface{})
				return redash.VisualizationColumnOptions{
					// Shared
					Visible:      column["visible"].(bool),
					Name:         column["name"].(string),
					Title:        column["title"].(string),
					AlignContent: column["alignContent"].(string),
					AllowSearch:  column["allowSearch"].(bool),
					Type:         column["type"].(string),
					DisplayAs:    column["displayAs"].(string),
					Order:        column["order"].(int),
					// Text
					AllowHTML:      column["allowHTML"].(bool),
					HighlightLinks: column["highlightLinks"].(bool),
					// Number
					NumberFormat: column["numberFormat"].(string),
					// Date/Time
					DateTimeFormat: column["dateTimeFormat"].(string),
					// Boolean
					BooleanValues: column["booleanValues"].([]string),
					// Link
					LinkUrlTemplate:   column["linkUrlTemplate"].(string),
					LinkTextTemplate:  column["linkTextTemplate"].(string),
					LinkOpenInNewTab:  column["linkOpenInNewTab"].(bool),
					LinkTitleTemplate: column["linkTitleTemplate"].(string),
					// Image
					ImageUrlTemplate:   column["imageUrlTemplate"].(string),
					ImageWidth:         column["imageWidth"].(string),
					ImageHeight:        column["imageHeight"].(string),
					ImageTitleTemplate: column["imageTitleTemplate"].(string),
				}
			}),
		}
		break
	default:
		return diag.Errorf("Invalid visualization type: %s", vType)
	}

	payload := redash.VisualizationCreatePayload{
		QueryId:     d.Get("query_id").(int),
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Options:     vOptions,
	}
	visualization, err := c.CreateVisualization(&payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(visualization.ID))

	return diags
}

func resourceRedashVisualizationUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	payload := redash.VisualizationUpdatePayload{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Options:     redash.VisualizationOptions{},
	}
	visualization, err := c.UpdateVisualization(id, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", visualization.Name)

	return diags
}

func resourceRedashVisualizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteVisualization(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

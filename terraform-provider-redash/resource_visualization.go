package main

import (
	"context"
	"strconv"

	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/samber/lo"
)

func resourceRedashVisualization() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashVisualizationRead,
		CreateContext: resourceRedashVisualizationCreate,
		UpdateContext: resourceRedashVisualizationUpdate,
		DeleteContext: resourceRedashVisualizationDelete,
		Schema: map[string]*schema.Schema{
			// Base Data
			"visualization_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
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
				ForceNew: true,
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
							Type:     schema.TypeInt,
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
									"image_title_template": {
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
						"error_y": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"visible": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"type": {
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
										Optional: true,
									},
									"error_y": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"visible": {
													Type:     schema.TypeBool,
													Required: true,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
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
										Optional: true,
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

	// TODO: fix map/array with key transformation bugs which prevent updates
	// Base Data
	_ = d.Set("visualization_id", visualization.ID)
	_ = d.Set("name", visualization.Name)
	_ = d.Set("description", visualization.Description)
	// Options
	_ = d.Set("type", visualization.Type)
	if visualization.Type == "TABLE" {
		_ = d.Set("table_options", visualization.Options)
	} else if visualization.Type == "CHART" {
		_ = d.Set("chart_options", visualization.Options)
	}

	return diags
}

func resourceRedashVisualizationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	var vOptions interface{}
	vType := d.Get("type").(string)
	switch vType {
	case "TABLE":
		tableOptions := d.Get("table_options").([]interface{})[0].(map[string]interface{})
		tableColumns := tableOptions["columns"].([]interface{})
		vOptions = redash.TableOptions{
			ItemsPerPage: tableOptions["items_per_page"].(int),
			Columns: lo.Map(tableColumns, func(item interface{}, _ int) redash.TableColumn {
				column := item.(map[string]interface{})
				return redash.TableColumn{
					// Shared
					Visible: column["visible"].(bool),
					Name:    column["name"].(string),
					Title:   column["title"].(string),
					// Type
					Type:         column["type"].(string),
					DisplayAs:    column["display_as"].(string),
					AlignContent: column["align_content"].(string),
					AllowSearch:  column["allow_search"].(bool),
					Order:        column["order"].(int),
					// Text
					AllowHTML:      column["allow_html"].(bool),
					HighlightLinks: column["highlight_links"].(bool),
					// Number
					NumberFormat: column["number_format"].(string),
					// Date/Time
					DateTimeFormat: column["date_time_format"].(string),
					// Boolean
					BooleanValues: lo.Map(column["boolean_values"].([]interface{}), func(item interface{}, _ int) string {
						return item.(string)
					}),
					// Link
					LinkUrlTemplate:   column["link_url_template"].(string),
					LinkTextTemplate:  column["link_text_template"].(string),
					LinkOpenInNewTab:  column["link_open_in_new_tab"].(bool),
					LinkTitleTemplate: column["link_title_template"].(string),
					// Image
					ImageUrlTemplate:   column["image_url_template"].(string),
					ImageTitleTemplate: column["image_title_template"].(string),
					ImageWidth:         column["image_width"].(string),
					ImageHeight:        column["image_height"].(string),
				}
			}),
		}
		break
	case "CHART":
		chartOptions := d.Get("chart_options").([]interface{})[0].(map[string]interface{})

		var chartLegend map[string]interface{}
		dChartLegend := chartOptions["legend"].([]interface{})[0]
		if dChartLegend != nil {
			chartLegend = dChartLegend.(map[string]interface{})
		}

		var chartSeries map[string]interface{}
		dChartSeries := chartOptions["series"].([]interface{})[0]
		if dChartSeries != nil {
			chartSeries = dChartSeries.(map[string]interface{})
		}

		var chartXAxis map[string]interface{}
		dChartXAxis := chartOptions["x_axis"].([]interface{})[0]
		if dChartXAxis != nil {
			chartXAxis = dChartXAxis.(map[string]interface{})
		}

		chartYAxis := chartOptions["y_axis"].([]interface{})
		chartSeriesOptions := chartOptions["series_options"].([]interface{})
		vOptions = redash.ChartOptions{
			// General
			GlobalSeriesType: chartOptions["global_series_type"].(string),
			ColumnMapping: lo.Associate(
				chartOptions["column_mapping"].([]interface{}),
				func(item interface{}) (string, string) {
					axis := item.(map[string]interface{})["axis"].(string)
					column := item.(map[string]interface{})["column"].(string)
					return axis, column
				},
			),
			Legend: redash.ChartLegend{
				Enabled: chartLegend["enabled"].(bool),
				// Placement: chartLegend["placement"].(string),
			},
			Series: redash.ChartSeries{
				Stacking: lo.TernaryF(
					chartSeries["stacking"] != nil,
					func() *string { return lo.EmptyableToPtr(chartSeries["stacking"].(string)) },
					func() *string { return nil },
				),
			},
			MissingValuesAsZero: chartOptions["missing_values_as_zero"].(bool),
			// CHART TYPE - X-Axis
			XAxis: redash.ChartXAxis{
				Type: chartXAxis["type"].(string),
				Labels: struct {
					Enabled bool `json:"enabled"`
				}{
					Enabled: chartXAxis["labels"].([]interface{})[0].(map[string]interface{})["enabled"].(bool),
				},
			},
			SortX: chartOptions["sort_x"].(bool),
			// CHART TYPE - Y-Axis
			YAxis: lo.Map(chartYAxis, func(item interface{}, _ int) redash.ChartYAxis {
				yAxis := item.(map[string]interface{})

				return redash.ChartYAxis{
					Type:     yAxis["type"].(string),
					Opposite: yAxis["opposite"].(bool),
				}
			}),
			// CHART TYPE - Series
			SeriesOptions: lo.Associate(chartSeriesOptions, func(value interface{}) (string, redash.ChartSeriesOption) {
				seriesOption := value.(map[string]interface{})

				return seriesOption["name"].(string), redash.ChartSeriesOption{
					ZIndex: seriesOption["z_index"].(int),
					Index:  seriesOption["index"].(int),
					Type:   seriesOption["type"].(string),
					YAxis:  seriesOption["y_axis"].(int),
				}
			}),
			// CHART TYPE - Colors
			// CHART TYPE - Data Labels
			ShowDataLabels: chartOptions["show_data_labels"].(bool),
			NumberFormat:   chartOptions["number_format"].(string),
			PercentFormat:  chartOptions["percent_format"].(string),
			DateTimeFormat: chartOptions["date_time_format"].(string),
			TextFormat:     chartOptions["text_format"].(string),
		}
		break
	default:
		return diag.Errorf("Invalid visualization type: %s", vType)
	}

	payload := redash.VisualizationCreatePayload{
		// Base Data
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// Options
		Type:    d.Get("type").(string),
		Options: vOptions,
		// References
		QueryId: d.Get("query_id").(int),
	}
	visualization, err := c.CreateVisualization(&payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(visualization.ID))
	_ = d.Set("visualization_id", visualization.ID)

	return diags
}

func resourceRedashVisualizationUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	var vOptions interface{}
	vType := d.Get("type").(string)
	switch vType {
	case "TABLE":
		tableOptions := d.Get("table_options").([]interface{})[0].(map[string]interface{})
		tableColumns := tableOptions["columns"].([]interface{})
		vOptions = redash.TableOptions{
			ItemsPerPage: tableOptions["items_per_page"].(int),
			Columns: lo.Map(tableColumns, func(item interface{}, _ int) redash.TableColumn {
				column := item.(map[string]interface{})
				return redash.TableColumn{
					// Shared
					Visible: column["visible"].(bool),
					Name:    column["name"].(string),
					Title:   column["title"].(string),
					// Type
					Type:         column["type"].(string),
					DisplayAs:    column["display_as"].(string),
					AlignContent: column["align_content"].(string),
					AllowSearch:  column["allow_search"].(bool),
					Order:        column["order"].(int),
					// Text
					AllowHTML:      column["allow_html"].(bool),
					HighlightLinks: column["highlight_links"].(bool),
					// Number
					NumberFormat: column["number_format"].(string),
					// Date/Time
					DateTimeFormat: column["date_time_format"].(string),
					// Boolean
					BooleanValues: lo.Map(column["boolean_values"].([]interface{}), func(item interface{}, _ int) string {
						return item.(string)
					}),
					// Link
					LinkUrlTemplate:   column["link_url_template"].(string),
					LinkTextTemplate:  column["link_text_template"].(string),
					LinkOpenInNewTab:  column["link_open_in_new_tab"].(bool),
					LinkTitleTemplate: column["link_title_template"].(string),
					// Image
					ImageUrlTemplate:   column["image_url_template"].(string),
					ImageTitleTemplate: column["image_title_template"].(string),
					ImageWidth:         column["image_width"].(string),
					ImageHeight:        column["image_height"].(string),
				}
			}),
		}
		break
	case "CHART":
		chartOptions := d.Get("chart_options").([]interface{})[0].(map[string]interface{})

		var chartLegend map[string]interface{}
		dChartLegend := chartOptions["legend"].([]interface{})[0]
		if dChartLegend != nil {
			chartLegend = dChartLegend.(map[string]interface{})
		}

		var chartSeries map[string]interface{}
		dChartSeries := chartOptions["series"].([]interface{})[0]
		if dChartSeries != nil {
			chartSeries = dChartSeries.(map[string]interface{})
		}

		var chartXAxis map[string]interface{}
		dChartXAxis := chartOptions["x_axis"].([]interface{})[0]
		if dChartXAxis != nil {
			chartXAxis = dChartXAxis.(map[string]interface{})
		}

		chartYAxis := chartOptions["y_axis"].([]interface{})
		chartSeriesOptions := chartOptions["series_options"].([]interface{})
		vOptions = redash.ChartOptions{
			// General
			GlobalSeriesType: chartOptions["global_series_type"].(string),
			ColumnMapping: lo.Associate(
				chartOptions["column_mapping"].([]interface{}),
				func(item interface{}) (string, string) {
					axis := item.(map[string]interface{})["axis"].(string)
					column := item.(map[string]interface{})["column"].(string)
					return axis, column
				},
			),
			Legend: redash.ChartLegend{
				Enabled: chartLegend["enabled"].(bool),
				// Placement: chartLegend["placement"].(string),
			},
			Series: redash.ChartSeries{
				Stacking: lo.TernaryF(
					chartSeries["stacking"] != nil,
					func() *string { return lo.EmptyableToPtr(chartSeries["stacking"].(string)) },
					func() *string { return nil },
				),
			},
			MissingValuesAsZero: chartOptions["missing_values_as_zero"].(bool),
			// CHART TYPE - X-Axis
			XAxis: redash.ChartXAxis{
				Type: chartXAxis["type"].(string),
				Labels: struct {
					Enabled bool `json:"enabled"`
				}{
					Enabled: chartXAxis["labels"].([]interface{})[0].(map[string]interface{})["enabled"].(bool),
				},
			},
			SortX: chartOptions["sort_x"].(bool),
			// CHART TYPE - Y-Axis
			YAxis: lo.Map(chartYAxis, func(item interface{}, _ int) redash.ChartYAxis {
				yAxis := item.(map[string]interface{})

				return redash.ChartYAxis{
					Type:     yAxis["type"].(string),
					Opposite: yAxis["opposite"].(bool),
				}
			}),
			// CHART TYPE - Series
			SeriesOptions: lo.Associate(chartSeriesOptions, func(value interface{}) (string, redash.ChartSeriesOption) {
				seriesOption := value.(map[string]interface{})

				return seriesOption["name"].(string), redash.ChartSeriesOption{
					ZIndex: seriesOption["z_index"].(int),
					Index:  seriesOption["index"].(int),
					Type:   seriesOption["type"].(string),
					YAxis:  seriesOption["y_axis"].(int),
				}
			}),
			// CHART TYPE - Colors
			// CHART TYPE - Data Labels
			ShowDataLabels: chartOptions["show_data_labels"].(bool),
			NumberFormat:   chartOptions["number_format"].(string),
			PercentFormat:  chartOptions["percent_format"].(string),
			DateTimeFormat: chartOptions["date_time_format"].(string),
			TextFormat:     chartOptions["text_format"].(string),
		}
		break
	default:
		return diag.Errorf("Invalid visualization type: %s", vType)
	}

	payload := redash.VisualizationUpdatePayload{
		// Base Data
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// Options
		Type:    d.Get("type").(string),
		Options: vOptions,
	}
	_, err = c.UpdateVisualization(id, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

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

	// d.SetId("")

	return diags
}

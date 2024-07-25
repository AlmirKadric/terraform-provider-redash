package main

import (
	"context"
	"strconv"

	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRedashWidget() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashWidgetRead,
		CreateContext: resourceRedashWidgetCreate,
		UpdateContext: resourceRedashWidgetUpdate,
		DeleteContext: resourceRedashWidgetDelete,
		Schema: map[string]*schema.Schema{
			// Base Data
			"widget_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dashboard_slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dashboard_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			//
			"text": {
				Type:     schema.TypeString,
				Required: true,
				// Default:  "",
			},
			"width": {
				Type:     schema.TypeInt,
				Required: true,
				// Default:  6,
			},
			// References
			"visualization_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			// Options
			"options": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_hidden": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"parameter_mappings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"map_to": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"title": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"position": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_height": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"size_x": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"size_y": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"max_size_y": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"max_size_x": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"min_size_y": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"min_size_x": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"col": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"row": {
										Type:     schema.TypeInt,
										Required: true,
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

func resourceRedashWidgetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	widget, err := c.GetWidget(d.Get("dashboard_slug").(string), id)
	if err != nil {
		return diag.FromErr(err)
	}

	// Base Data
	_ = d.Set("widget_id", widget.ID)
	_ = d.Set("dashboard_id", widget.DashboardID)
	//
	_ = d.Set("text", widget.Text)
	_ = d.Set("width", widget.Width)
	// References
	_ = d.Set("visualization_id", widget.Visualization.ID)
	// Options
	_ = d.Set("options", widget.Options)

	return diags
}

func resourceRedashWidgetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	dashboard, err := c.GetDashboard(d.Get("dashboard_slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	dOptions := d.Get("options").([]interface{})[0].(map[string]interface{})
	dPosition := dOptions["position"].([]interface{})[0].(map[string]interface{})
	// dParameterMappings := dOptions["parameter_mappings"].([]interface{})

	options := redash.WidgetOptions{
		IsHidden: dOptions["is_hidden"].(bool),
		Position: redash.WidgetPosition{
			AutoHeight: dPosition["auto_height"].(bool),
			SizeX:      dPosition["size_x"].(int),
			SizeY:      dPosition["size_y"].(int),
			MaxSizeY:   dPosition["max_size_y"].(int),
			MaxSizeX:   dPosition["max_size_x"].(int),
			MinSizeY:   dPosition["min_size_y"].(int),
			MinSizeX:   dPosition["min_size_x"].(int),
			Col:        dPosition["col"].(int),
			Row:        dPosition["row"].(int),
		},
		ParameterMappings: nil,
	}

	dVisualizationID := d.Get("visualization_id").(int)

	var visualizationID *int = nil
	if dVisualizationID != 0 {
		visualizationID = &dVisualizationID
	}

	widget, err := c.CreateWidget(&redash.WidgetCreatePayload{
		// Base Data
		DashboardID: dashboard.ID,
		//
		Text:  d.Get("text").(string),
		Width: d.Get("width").(int),
		// References
		VisualizationID: visualizationID,
		// Options
		Options: options,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(widget.ID))
	_ = d.Set("widget_id", widget.ID)
	_ = d.Set("dashboard_id", dashboard.ID)

	return diags
}

func resourceRedashWidgetUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dOptions := d.Get("options").([]interface{})[0].(map[string]interface{})
	dPosition := dOptions["position"].([]interface{})[0].(map[string]interface{})
	// dParameterMappings := dOptions["parameter_mappings"].([]interface{})

	options := redash.WidgetOptions{
		IsHidden: dOptions["is_hidden"].(bool),
		Position: redash.WidgetPosition{
			AutoHeight: dPosition["auto_height"].(bool),
			SizeX:      dPosition["size_x"].(int),
			SizeY:      dPosition["size_y"].(int),
			MaxSizeY:   dPosition["max_size_y"].(int),
			MaxSizeX:   dPosition["max_size_x"].(int),
			MinSizeY:   dPosition["min_size_y"].(int),
			MinSizeX:   dPosition["min_size_x"].(int),
			Col:        dPosition["col"].(int),
			Row:        dPosition["row"].(int),
		},
		ParameterMappings: nil,
	}

	dVisualizationID := d.Get("visualization_id").(int)

	var visualizationID *int = nil
	if dVisualizationID != 0 {
		visualizationID = &dVisualizationID
	}

	_, err = c.UpdateWidget(id, &redash.WidgetUpdatePayload{
		//
		Text:  d.Get("text").(string),
		Width: d.Get("width").(int),
		// References
		VisualizationID: visualizationID,
		// Options
		Options: options,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRedashWidgetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteWidget(id)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("")

	return diags
}

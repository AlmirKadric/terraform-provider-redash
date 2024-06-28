package main

import (
	"context"
	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
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
							Required: true,
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
										Required: true,
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
									"height": {
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

	_, err = c.GetWidget(d.Get("dashboard_slug").(string), id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRedashWidgetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	dashboard, err := c.GetDashboard(d.Get("dashboard_slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	widget, err := c.CreateWidget(&redash.WidgetCreatePayload{
		DashboardID:     dashboard.ID,
		VisualizationID: d.Get("visualization_id").(int),
		Options: redash.WidgetOptions{
			IsHidden: d.Get("is_hidden").(bool),
			Position: redash.WidgetPosition{
				AutoHeight: d.Get("auto_height").(bool),
				SizeX:      d.Get("width").(int),
				SizeY:      d.Get("height").(int),
				MaxSizeY:   1000,
				MaxSizeX:   6,
				MinSizeY:   1,
				MinSizeX:   2,
				Col:        d.Get("column").(int),
				Row:        d.Get("row").(int),
			},
			ParameterMappings: nil,
		},
		Text:  d.Get("text").(string),
		Width: 1,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(widget.ID))
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

	_, err = c.UpdateWidget(id, &redash.WidgetUpdatePayload{
		Options: redash.WidgetOptions{
			IsHidden: d.Get("is_hidden").(bool),
			Position: redash.WidgetPosition{
				AutoHeight: d.Get("auto_height").(bool),
				SizeX:      d.Get("width").(int),
				SizeY:      d.Get("height").(int),
				MaxSizeY:   1000,
				MaxSizeX:   6,
				MinSizeY:   1,
				MinSizeX:   2,
				Col:        d.Get("column").(int),
				Row:        d.Get("row").(int),
			},
			ParameterMappings: nil,
		},
		Text:  d.Get("text").(string),
		Width: 1,
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

	d.SetId("")

	return diags
}

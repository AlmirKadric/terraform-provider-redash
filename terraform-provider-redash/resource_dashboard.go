package main

import (
	"context"
	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceRedashDashboard() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashDashboardRead,
		CreateContext: resourceRedashDashboardCreate,
		UpdateContext: resourceRedashDashboardUpdate,
		DeleteContext: resourceRedashDashboardArchive,
		Schema: map[string]*schema.Schema{
			// Base Data
			"dashboard_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Options
			// "layout"
			// State
			"is_favorite": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_archived": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_draft": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"dashboard_filters_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// Metadata
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Dashboard Specific
			"public_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"can_edit": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRedashDashboardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	dashboard, err := c.GetDashboard(d.Get("slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", dashboard.Name)
	_ = d.Set("slug", dashboard.Slug)

	return diags
}

func resourceRedashDashboardCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	createPayload := redash.DashboardCreatePayload{
		Name: d.Get("name").(string),
	}
	dashboard, err := c.CreateDashboard(&createPayload)
	if err != nil {
		return nil
	}

	d.SetId(strconv.Itoa(dashboard.ID))
	_ = d.Set("name", dashboard.Name)
	_ = d.Set("slug", dashboard.Slug)

	return diags
}

func resourceRedashDashboardUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	updatePayload := redash.DashboardUpdatePayload{
		Name: d.Get("name").(string),
	}
	dashboard, err := c.UpdateDashboard(id, &updatePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", dashboard.Name)
	_ = d.Set("slug", dashboard.Slug)

	return diags
}

func resourceRedashDashboardArchive(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	err := c.ArchiveDashboard(d.Get("slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

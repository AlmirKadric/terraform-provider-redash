package main

import (
	"context"
	"strconv"

	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/samber/lo"
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
				Computed: true,
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
				Computed: true,
			},
			"can_edit": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
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

	// Base Data
	_ = d.Set("dashboard_id", dashboard.ID)
	_ = d.Set("name", dashboard.Name)
	_ = d.Set("slug", dashboard.Slug)
	// Options
	// "layout"
	// State
	_ = d.Set("is_favorite", dashboard.IsFavorite)
	_ = d.Set("is_archived", dashboard.IsArchived)
	_ = d.Set("is_draft", dashboard.IsDraft)
	_ = d.Set("dashboard_filters_enabled", dashboard.DashboardFiltersEnabled)
	_ = d.Set("version", dashboard.Version)
	// Metadata
	_ = d.Set("tags", dashboard.Tags)
	// Dashboard Specific
	_ = d.Set("public_url", dashboard.PublicUrl)
	_ = d.Set("can_edit", dashboard.CanEdit)
	_ = d.Set("api_key", dashboard.APIKey)

	return diags
}

func resourceRedashDashboardCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	createPayload := redash.DashboardCreatePayload{
		// Base Data
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
		// Options
		// Layout                  []interface{}     `json:"layout"`
		// State
		IsFavorite:              d.Get("is_favorite").(bool),
		IsArchived:              d.Get("is_archived").(bool),
		IsDraft:                 d.Get("is_draft").(bool),
		DashboardFiltersEnabled: d.Get("dashboard_filters_enabled").(bool),
		// Metadata
		Tags: lo.Map(d.Get("tags").([]interface{}), func(item interface{}, _ int) string {
			return item.(string)
		}),
	}
	dashboard, err := c.CreateDashboard(&createPayload)
	if err != nil {
		return nil
	}

	d.SetId(strconv.Itoa(dashboard.ID))
	_ = d.Set("dashboard_id", dashboard.ID)
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
		// Base Data
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
		// Options
		// Layout                  []interface{}     `json:"layout"`
		// State
		IsFavorite:              d.Get("is_favorite").(bool),
		IsArchived:              d.Get("is_archived").(bool),
		IsDraft:                 d.Get("is_draft").(bool),
		DashboardFiltersEnabled: d.Get("dashboard_filters_enabled").(bool),
		// Metadata
		Tags: lo.Map(d.Get("tags").([]interface{}), func(item interface{}, _ int) string {
			return item.(string)
		}),
	}
	_, err = c.UpdateDashboard(id, &updatePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRedashDashboardArchive(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	err := c.ArchiveDashboard(d.Get("slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("")

	return diags
}

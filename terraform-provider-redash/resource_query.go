package main

import (
	"context"
	"strconv"

	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRedashQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedashQueryCreate,
		ReadContext:   resourceRedashQueryRead,
		UpdateContext: resourceRedashQueryUpdate,
		DeleteContext: resourceRedashQueryArchive,
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
			// Query
			"data_source_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Options
			"parameters": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"locals": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enum_options": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			// State
			"is_draft": {
				Type:     schema.TypeBool,
				// Optional: true,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				// Optional: true,
				Required: true,
			},
		},
	}
}

func resourceRedashQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	parameters := d.Get("parameters").([]interface{})

	options := redash.QueryOptions{
		Parameters: make([]redash.QueryOptionsParameter, len(parameters)),
	}

	for i, p := range parameters {
		parameter := p.(map[string]interface{})
		options.Parameters[i] = redash.QueryOptionsParameter{
			Title:       parameter["title"].(string),
			Name:        parameter["name"].(string),
			Type:        parameter["type"].(string),
			EnumOptions: parameter["enum_options"].(string),
			Locals:      parameter["locals"].([]interface{}),
			Value:       parameter["value"].(string),
		}
	}

	createPayload := redash.QueryCreatePayload{
		Name:         d.Get("name").(string),
		Query:        d.Get("query").(string),
		DataSourceID: d.Get("data_source_id").(int),
		Description:  d.Get("description").(string),
		IsDraft:      d.Get("is_draft").(bool),
		Options:      options,
		Version:      d.Get("version").(int),
	}

	query, err := c.CreateQuery(&createPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(query.ID))
	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

	return diags
}

func resourceRedashQueryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	query, err := c.GetQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", query.Name)
	_ = d.Set("query", query.Query)
	_ = d.Set("data_source_id", query.DataSourceID)
	_ = d.Set("description", query.Description)
	_ = d.Set("is_draft", query.IsDraft)
	_ = d.Set("parameters", query.Options.Parameters)
	_ = d.Set("version", query.Version)

	return diags
}

func resourceRedashQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	parameters := d.Get("parameters").([]interface{})

	options := redash.QueryOptions{
		Parameters: make([]redash.QueryOptionsParameter, len(parameters)),
	}

	for i, p := range parameters {
		parameter := p.(map[string]interface{})
		options.Parameters[i] = redash.QueryOptionsParameter{
			Title:       parameter["title"].(string),
			Name:        parameter["name"].(string),
			Type:        parameter["type"].(string),
			EnumOptions: parameter["enum_options"].(string),
			Locals:      parameter["locals"].([]interface{}),
			Value:       parameter["value"].(string),
		}
	}

	updatePayload := redash.QueryUpdatePayload{
		Name:         d.Get("name").(string),
		Query:        d.Get("query").(string),
		DataSourceID: d.Get("data_source_id").(int),
		Description:  d.Get("description").(string),
		IsDraft:      d.Get("is_draft").(bool),
		Options:      options,
		Version:      d.Get("version").(int),
	}

	_, err = c.UpdateQuery(id, &updatePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

	return diags
}

func resourceRedashQueryArchive(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.ArchiveQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

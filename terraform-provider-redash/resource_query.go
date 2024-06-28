package main

import (
	"context"
	"fmt"
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
			"query_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"query_hash": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Options
			"options": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
									"parent_query_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									// "locals": {
									// 	Type:     schema.TypeList,
									// 	Required: true,
									// 	Elem: &schema.Schema{
									// 		Type: schema.TypeString,
									// 	},
									// },
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"string": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"range": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start": {
																Type:     schema.TypeString,
																Required: true,
															},
															"end": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
											},
										},
									},
									"enum_options": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"global": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			// State
			"is_draft": {
				Type: schema.TypeBool,
				// Optional: true,
				Required: true,
			},
			"is_archived": {
				Type: schema.TypeBool,
				// Optional: true,
				Required: true,
			},
			"is_safe": {
				Type: schema.TypeBool,
				// Optional: true,
				Required: true,
			},
			"version": {
				Type: schema.TypeInt,
				// Optional: true,
				Required: true,
			},
			// Metadata
			"api_key": {
				Type: schema.TypeString,
				// Optional: true,
				Required: true,
			},
			"tags": {
				Type: schema.TypeList,
				// Optional: true,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"latest_query_data_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"day_of_week": {
							Type:     schema.TypeString,
							Required: true,
						},
						// "until": {
						// 	Type:     schema.TypeString,
						// 	Required: true,
						// },
					},
				},
			},
			// Query Specific
			"is_favorite": {
				Type: schema.TypeBool,
				// Optional: true,
				Required: true,
			},
			"can_edit": {
				Type: schema.TypeBool,
				// Optional: true,
				Required: true,
			},
		},
	}
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

	// Base Data
	_ = d.Set("query_id", query.ID)
	_ = d.Set("name", query.Name)
	_ = d.Set("description", query.Description)
	// Query
	_ = d.Set("data_source_id", query.DataSourceID)
	_ = d.Set("query", query.Query)
	_ = d.Set("query_hash", query.QueryHash)
	// Options
	_ = d.Set("options", query.Options)
	// State
	_ = d.Set("is_draft", query.IsDraft)
	_ = d.Set("is_archived", query.IsArchived)
	_ = d.Set("is_safe", query.IsSafe)
	_ = d.Set("version", query.Version)
	// Metadata
	_ = d.Set("api_key", query.APIKey)
	_ = d.Set("tags", query.Tags)
	_ = d.Set("latest_query_data_id", query.LatestQueryDataID)
	_ = d.Set("schedule", query.Schedule)
	// Query Specific
	_ = d.Set("is_favorite", query.IsFavorite)
	_ = d.Set("can_edit", query.CanEdit)

	return diags
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

		pType := parameter["type"].(string)
		var pValue interface{}
		switch pType {
		case "string":
			pValue = parameter["value"].([]map[string]interface{})[0]["string"].([]interface{})
			break
		case "range":
			pValue = parameter["value"].([]map[string]interface{})[0]["range"].([]interface{})
			break
		default:
			return diag.FromErr(fmt.Errorf("Invalid parameter type: %s", pType))
		}

		options.Parameters[i] = redash.QueryOptionsParameter{
			Name:  parameter["name"].(string),
			Title: parameter["title"].(string),

			ParentQueryId: parameter["parent_query_id"].(int),

			// Locals: parameter["locals"].([]interface{}),

			Type:        pType,
			Value:       pValue,
			EnumOptions: parameter["enum_options"].(string),

			Global: parameter["global"].(bool),
		}
	}

	schedule := d.Get("schedule").([]interface{})[0].(map[string]interface{})

	createPayload := redash.QueryCreatePayload{
		// Base Data
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// Query
		DataSourceID: d.Get("data_source_id").(int),
		Query:        d.Get("query").(string),
		QueryHash:    d.Get("query_hash").(string),
		// Options
		Options: options,
		// State
		IsDraft:    d.Get("is_draft").(bool),
		IsArchived: d.Get("is_archived").(bool),
		IsSafe:     d.Get("is_safe").(bool),
		Version:    d.Get("version").(int),
		// Metadata
		APIKey:            d.Get("api_key").(string),
		Tags:              d.Get("tags").([]string),
		LatestQueryDataID: d.Get("latest_query_data_id").(int),
		Schedule: redash.QuerySchedule{
			Interval:  schedule["interval"].(int),
			Time:      schedule["time"].(string),
			DayOfWeek: schedule["day_of_week"].(string),
			// Until:     schedule["until"].(interface{}),
		},
		// Query Specific
		IsFavorite: d.Get("is_favorite").(bool),
		CanEdit:    d.Get("can_edit").(bool),
	}

	query, err := c.CreateQuery(&createPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("query_id", query.ID)
	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

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

		pType := parameter["type"].(string)
		var pValue interface{}
		switch pType {
		case "string":
			pValue = parameter["value"].([]map[string]interface{})[0]["string"].([]interface{})
			break
		case "range":
			pValue = parameter["value"].([]map[string]interface{})[0]["range"].([]interface{})
			break
		default:
			return diag.FromErr(fmt.Errorf("Invalid parameter type: %s", pType))
		}

		options.Parameters[i] = redash.QueryOptionsParameter{
			Name:  parameter["name"].(string),
			Title: parameter["title"].(string),

			ParentQueryId: parameter["parent_query_id"].(int),

			// Locals: parameter["locals"].([]interface{}),

			Type:        pType,
			Value:       pValue,
			EnumOptions: parameter["enum_options"].(string),

			Global: parameter["global"].(bool),
		}
	}

	schedule := d.Get("schedule").([]interface{})[0].(map[string]interface{})

	updatePayload := redash.QueryUpdatePayload{
		// Base Data
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// Query
		DataSourceID: d.Get("data_source_id").(int),
		Query:        d.Get("query").(string),
		QueryHash:    d.Get("query_hash").(string),
		// Options
		Options: options,
		// State
		IsDraft:    d.Get("is_draft").(bool),
		IsArchived: d.Get("is_archived").(bool),
		IsSafe:     d.Get("is_safe").(bool),
		Version:    d.Get("version").(int),
		// Metadata
		APIKey:            d.Get("api_key").(string),
		Tags:              d.Get("tags").([]string),
		LatestQueryDataID: d.Get("latest_query_data_id").(int),
		Schedule: redash.QuerySchedule{
			Interval:  schedule["interval"].(int),
			Time:      schedule["time"].(string),
			DayOfWeek: schedule["day_of_week"].(string),
			// Until:     schedule["until"].(interface{}),
		},
		// Query Specific
		IsFavorite: d.Get("is_favorite").(bool),
		CanEdit:    d.Get("can_edit").(bool),
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

	// d.SetId("")

	return diags
}

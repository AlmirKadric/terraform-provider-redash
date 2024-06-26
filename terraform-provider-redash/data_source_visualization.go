package main

import (
	"context"
	"fmt"
	"github.com/AlmirKadric/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRedashVisualization() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"query_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"visualization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: dataSourceRedashVisualizationRead,
	}
}

func dataSourceRedashVisualizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	queryID := d.Get("query_id").(int)
	visualizationID := d.Get("visualization_id").(int)

	visualization, err := c.GetVisualization(queryID, visualizationID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(visualization.ID))
	_ = d.Set("name", visualization.Name)

	return diags
}

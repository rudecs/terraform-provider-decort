/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Terraform DECORT provider - manage resources provided by DECORT (Digital Energy Cloud
Orchestration Technology) with Terraform by Hashicorp.

Source code: https://github.com/rudecs/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://github.com/rudecs/terraform-provider-decort/wiki
*/

package grid

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenGrid(d *schema.ResourceData, grid *Grid) {
	d.Set("name", grid.Name)
	d.Set("flag", grid.Flag)
	d.Set("gid", grid.Gid)
	d.Set("guid", grid.Guid)
	d.Set("location_code", grid.LocationCode)
	d.Set("id", grid.Id)
}

func dataSourceGridRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	grid, err := utilityGridCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(grid.Id))
	flattenGrid(d, grid)

	return nil
}

func dataSourceGetGridSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"grid_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"flag": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"location_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceGrid() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceGridRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceGetGridSchemaMake(),
	}
}

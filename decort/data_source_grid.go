/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>, <svs1370@gmail.com>

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
This file is part of Terraform (by Hashicorp) provider for Digital Energy Cloud Orchestration
Technology platfom.

Visit https://github.com/rudecs/terraform-provider-decort for full source code package and updates.
*/

package decort

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func flattenGrid(d *schema.ResourceData, grid *Grid) {
	d.Set("name", grid.Name)
	d.Set("flag", grid.Flag)
	d.Set("gid", grid.Gid)
	d.Set("guid", grid.Guid)
	d.Set("location_code", grid.LocationCode)
	d.Set("id", grid.Id)
	return
}

func dataSourceGridRead(d *schema.ResourceData, m interface{}) error {
	grid, err := utilityGridCheckPresence(d, m)
	if err != nil {
		return err
	}
	d.SetId("1234")
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

func dataSourceGrid() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceGridRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceGetGridSchemaMake(),
	}
}

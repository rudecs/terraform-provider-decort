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

package locations

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/flattens"
)

func flattenLocationsList(ll LocationsList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, l := range ll {
		temp := map[string]interface{}{
			"ckey":          l.CKey,
			"meta":          flattens.FlattenMeta(l.Meta),
			"flag":          l.Flag,
			"gid":           l.GridID,
			"guid":          l.Guid,
			"id":            l.Id,
			"location_code": l.LocationCode,
			"name":          l.Name,
		}
		res = append(res, temp)
	}
	return res

}

func dataSourceLocationsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	locations, err := utilityLocationsListCheckPresence(ctx, d, m)
	if err != nil {

		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())

	d.Set("items", flattenLocationsList(locations))

	return nil
}

func dataSourceLocationsListSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"page": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "page number",
		},
		"size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "page size",
		},
		"items": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Locations list",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ckey": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"meta": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"flag": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"gid": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Grid id",
					},
					"guid": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "location id",
					},
					"id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "location id",
					},
					"location_code": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Location code",
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Location name",
					},
				},
			},
		},
	}
}

func DataSourceLocationsList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceLocationsListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceLocationsListSchemaMake(),
	}
}

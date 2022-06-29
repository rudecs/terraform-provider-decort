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

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenGridList(gl GridList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range gl {
		temp := map[string]interface{}{
			"name":          item.Name,
			"flag":          item.Flag,
			"gid":           item.Gid,
			"guid":          item.Guid,
			"location_code": item.LocationCode,
			"id":            item.Id,
		}

		res = append(res, temp)
	}
	return res
}

func dataSourceGridListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	gridList, err := utilityGridListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenGridList(gridList))

	return nil
}

func dataSourceGridListSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
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
			Description: "grid list",
			Elem: &schema.Resource{
				Schema: dataSourceGridSchemaMake(),
			},
		},
	}

	return rets
}

func dataSourceGridSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

func DataSourceGridList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceGridListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceGridListSchemaMake(),
	}
}

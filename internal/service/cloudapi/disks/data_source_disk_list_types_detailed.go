/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

package disks

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenDiskListTypesDetailed(tld TypesDetailedList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, typeListDetailed := range tld {
		temp := map[string]interface{}{
			"pools":  flattenListTypesDetailedPools(typeListDetailed.Pools),
			"sep_id": typeListDetailed.SepID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenListTypesDetailedPools(pools PoolList) []interface{} {
	res := make([]interface{}, 0)
	for _, pool := range pools {
		temp := map[string]interface{}{
			"name":  pool.Name,
			"types": pool.Types,
		}
		res = append(res, temp)
	}

	return res
}

func dataSourceDiskListTypesDetailedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	listTypesDetailed, err := utilityDiskListTypesDetailedCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenDiskListTypesDetailed(listTypesDetailed))
	return nil
}

func dataSourceDiskListTypesDetailedSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"pools": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Pool name",
								},
								"types": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
									Description: "The types of disk in terms of its role in compute: 'B=Boot, D=Data, T=Temp'",
								},
							},
						},
					},
					"sep_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Storage endpoint provider ID to create disk",
					},
				},
			},
		},
	}
	return res
}

func DataSourceDiskListTypesDetailed() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		ReadContext:   dataSourceDiskListTypesDetailedRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceDiskListTypesDetailedSchemaMake(),
	}
}

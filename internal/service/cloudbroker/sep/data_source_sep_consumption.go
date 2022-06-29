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

package sep

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceSepConsumptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sepCons, err := utilitySepConsumptionCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())

	d.Set("type", sepCons.Type)
	d.Set("total", flattenSepConsumption(sepCons.Total))
	d.Set("by_pool", flattenSepConsumptionPools(sepCons.ByPool))

	return nil
}

func flattenSepConsumptionPools(mp map[string]SepConsumptionInd) []map[string]interface{} {
	sh := make([]map[string]interface{}, 0)
	for k, v := range mp {
		temp := map[string]interface{}{
			"name":           k,
			"disk_count":     v.DiskCount,
			"disk_usage":     v.DiskUsage,
			"snapshot_count": v.SnapshotCount,
			"snapshot_usage": v.SnapshotUsage,
			"usage":          v.Usage,
			"usage_limit":    v.UsageLimit,
		}
		sh = append(sh, temp)
	}
	return sh
}

func flattenSepConsumption(sc SepConsumptionTotal) []map[string]interface{} {
	sh := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"capacity_limit": sc.CapacityLimit,
		"disk_count":     sc.DiskCount,
		"disk_usage":     sc.DiskUsage,
		"snapshot_count": sc.SnapshotCount,
		"snapshot_usage": sc.SnapshotUsage,
		"usage":          sc.Usage,
		"usage_limit":    sc.UsageLimit,
	}
	sh = append(sh, temp)
	return sh
}

func dataSourceSepConsumptionSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sep_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "sep id",
		},
		"by_pool": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "pool name",
					},
					"disk_count": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "number of disks",
					},
					"disk_usage": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "disk usage",
					},
					"snapshot_count": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "number of snapshots",
					},
					"snapshot_usage": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "snapshot usage",
					},
					"usage": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "usage",
					},
					"usage_limit": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "usage limit",
					},
				},
			},
			Description: "consumption divided by pool",
		},
		"total": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"capacity_limit": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"disk_count": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "number of disks",
					},
					"disk_usage": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "disk usage",
					},
					"snapshot_count": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "number of snapshots",
					},
					"snapshot_usage": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "snapshot usage",
					},
					"usage": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "usage",
					},
					"usage_limit": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "usage limit",
					},
				},
			},
			Description: "total consumption",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "sep type",
		},
	}
}

func DataSourceSepConsumption() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceSepConsumptionRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceSepConsumptionSchemaMake(),
	}
}

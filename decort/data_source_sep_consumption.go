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
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSepConsumptionRead(d *schema.ResourceData, m interface{}) error {
	sepCons, err := utilitySepConsumptionCheckPresence(d, m)
	if err != nil {
		return err
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

func dataSourceSepConsumption() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepConsumptionRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepConsumptionSchemaMake(),
	}
}

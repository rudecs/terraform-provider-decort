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

func dataSourceSepHitachiPoolRead(d *schema.ResourceData, m interface{}) error {
	hitachiSepPool, err := utilitySepHitachiPoolCheckPresence(d, m)
	if err != nil {
		return err
	}
	id := uuid.New()
	d.SetId(id.String())

	d.Set("pool", flattenSepHitachiPool(hitachiSepPool))

	return nil
}

func flattenSepHitachiPool(pool *HitachiConfigPool) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	t := map[string]interface{}{
		"clone_technology": pool.CloneTechnology,
		"id":               pool.Id,
		"max_l_dev_id":     pool.MaxLdevId,
		"min_l_dev_id":     pool.MinLdevId,
		"name":             pool.Name,
		"snapshot_pool_id": pool.SnapshotPoolId,
		"snapshotable":     pool.Snapshotable,
		"types":            pool.Types,
		"usage_limit":      pool.UsageLimit,
	}
	temp = append(temp, t)
	return temp
}

func dataSourceSepHitachiPoolSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"clone_technology": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"max_l_dev_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"min_l_dev_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"snapshot_pool_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"snapshotable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"types": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"usage_limit": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func dataSourceSepHitachiPool() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepHitachiPoolRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepHitachiPoolSchemaMake(),
	}
}

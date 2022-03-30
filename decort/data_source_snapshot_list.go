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

func flattenSnapshotList(gl SnapshotList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range gl {
		temp := map[string]interface{}{
			"label":     item.Label,
			"guid":      item.Guid,
			"disks":     item.Disks,
			"timestamp": item.Timestamp,
		}

		res = append(res, temp)
	}
	return res
}

func dataSourceSnapshotListRead(d *schema.ResourceData, m interface{}) error {
	snapshotList, err := utilitySnapshotListCheckPresence(d, m)
	if err != nil {
		return err
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenSnapshotList(snapshotList))

	return nil
}

func dataSourceSnapshotListSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"compute_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the compute instance to create snapshot for.",
		},
		"items": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "snapshot list",
			Elem: &schema.Resource{
				Schema: dataSourceSnapshotSchemaMake(),
			},
		},
	}

	return rets
}

func dataSourceSnapshotSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "text label for snapshot. Must be unique among this compute snapshots.",
		},
		"disks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "guid of the snapshot",
		},
		"timestamp": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "timestamp",
		},
	}
}

func dataSourceSnapshotList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSnapshotListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSnapshotListSchemaMake(),
	}
}

/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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

func flattenAccountDisksList(adl AccountDisksList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, ad := range adl {
		temp := map[string]interface{}{
			"disk_id":   ad.ID,
			"disk_name": ad.Name,
			"pool":      ad.Pool,
			"sep_id":    ad.SepId,
			"size_max":  ad.SizeMax,
			"type":      ad.Type,
		}
		res = append(res, temp)
	}
	return res

}

func dataSourceAccountDisksListRead(d *schema.ResourceData, m interface{}) error {
	accountDisksList, err := utilityAccountDisksListCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountDisksList(accountDisksList))

	return nil
}

func dataSourceAccountDisksListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the account",
		},
		"items": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Search Result",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"disk_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pool": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sep_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"size_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func dataSourceAccountDisksList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceAccountDisksListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceAccountDisksListSchemaMake(),
	}
}

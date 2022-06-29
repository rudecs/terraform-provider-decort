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
package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
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

func dataSourceAccountDisksListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountDisksList, err := utilityAccountDisksListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
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

func DataSourceAccountDisksList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceAccountDisksListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceAccountDisksListSchemaMake(),
	}
}

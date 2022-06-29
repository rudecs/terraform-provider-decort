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

func flattenAccountVinsList(avl AccountVinsList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, av := range avl {
		temp := map[string]interface{}{
			"account_id":     av.AccountId,
			"account_name":   av.AccountName,
			"computes":       av.Computes,
			"created_by":     av.CreatedBy,
			"created_time":   av.CreatedTime,
			"deleted_by":     av.DeletedBy,
			"deleted_time":   av.DeletedTime,
			"external_ip":    av.ExternalIP,
			"vin_id":         av.ID,
			"vin_name":       av.Name,
			"network":        av.Network,
			"pri_vnf_dev_id": av.PriVnfDevId,
			"rg_id":          av.RgId,
			"rg_name":        av.RgName,
			"status":         av.Status,
			"updated_by":     av.UpdatedBy,
			"updated_time":   av.UpdatedTime,
		}
		res = append(res, temp)
	}
	return res

}

func dataSourceAccountVinsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountVinsList, err := utilityAccountVinsListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountVinsList(accountVinsList))

	return nil
}

func dataSourceAccountVinsListSchemaMake() map[string]*schema.Schema {
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
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"computes": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"created_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"deleted_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"deleted_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"external_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vin_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vin_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pri_vnf_dev_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"rg_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"rg_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func DataSourceAccountVinsList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceAccountVinsListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceAccountVinsListSchemaMake(),
	}
}

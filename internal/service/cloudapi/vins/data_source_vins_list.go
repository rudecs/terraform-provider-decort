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

package vins

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenVinsList(vl VinsList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, v := range vl {
		temp := map[string]interface{}{
			"account_id":   v.AccountId,
			"account_name": v.AccountName,
			"created_by":   v.CreatedBy,
			"created_time": v.CreatedTime,
			"deleted_by":   v.DeletedBy,
			"deleted_time": v.DeletedTime,
			"external_ip":  v.ExternalIP,
			"vins_id":      v.ID,
			"vins_name":    v.Name,
			"network":      v.Network,
			"rg_id":        v.RGID,
			"rg_name":      v.RGName,
			"status":       v.Status,
			"updated_by":   v.UpdatedBy,
			"updated_time": v.UpdatedTime,
			"vxlan_id":     v.VXLanID,
		}
		res = append(res, temp)
	}
	return res
}

func dataSourceVinsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vinsList, err := utilityVinsListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenVinsList(vinsList))

	return nil
}

func dataSourceVinsListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"include_deleted": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "include deleted computes",
		},
		"page": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Page number",
		},
		"size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Page size",
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
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
					"vins_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vins_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network": {
						Type:     schema.TypeString,
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
					"vxlan_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func DataSourceVinsList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceVinsListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceVinsListSchemaMake(),
	}
}

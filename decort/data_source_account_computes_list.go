/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

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

func flattenAccountComputesList(acl AccountComputesList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, acc := range acl {
		temp := map[string]interface{}{
			"account_id":       acc.AccountId,
			"account_name":     acc.AccountName,
			"cpus":             acc.CPUs,
			"created_by":       acc.CreatedBy,
			"created_time":     acc.CreatedTime,
			"deleted_by":       acc.DeletedBy,
			"deleted_time":     acc.DeletedTime,
			"compute_id":       acc.ComputeId,
			"compute_name":     acc.ComputeName,
			"ram":              acc.RAM,
			"registered":       acc.Registered,
			"rg_id":            acc.RgId,
			"rg_name":          acc.RgName,
			"status":           acc.Status,
			"tech_status":      acc.TechStatus,
			"total_disks_size": acc.TotalDisksSize,
			"updated_by":       acc.UpdatedBy,
			"updated_time":     acc.UpdatedTime,
			"user_managed":     acc.UserManaged,
			"vins_connected":   acc.VinsConnected,
		}
		res = append(res, temp)
	}
	return res

}

func dataSourceAccountComputesListRead(d *schema.ResourceData, m interface{}) error {
	accountComputesList, err := utilityAccountComputesListCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountComputesList(accountComputesList))

	return nil
}

func dataSourceAccountComputesListSchemaMake() map[string]*schema.Schema {
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
					"cpus": {
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
					"compute_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"compute_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ram": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"registered": {
						Type:     schema.TypeBool,
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
					"tech_status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"total_disks_size": {
						Type:     schema.TypeInt,
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
					"user_managed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"vins_connected": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func dataSourceAccountComputesList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceAccountComputesListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceAccountComputesListSchemaMake(),
	}
}

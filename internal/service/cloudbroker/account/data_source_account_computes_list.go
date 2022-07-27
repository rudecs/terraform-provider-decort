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

func dataSourceAccountComputesListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountComputesList, err := utilityAccountComputesListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
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
				Schema: dataSourceAccountComputeSchemaMake(),
			},
		},
	}
	return res
}

func DataSourceAccountComputesList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceAccountComputesListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceAccountComputesListSchemaMake(),
	}
}

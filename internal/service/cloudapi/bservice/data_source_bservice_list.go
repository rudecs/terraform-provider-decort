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

package bservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenBasicServiceList(bsl BasicServiceList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bs := range bsl {
		temp := map[string]interface{}{
			"account_id":      bs.AccountId,
			"account_name":    bs.AccountName,
			"base_domain":     bs.BaseDomain,
			"created_by":      bs.CreatedBy,
			"created_time":    bs.CreatedTime,
			"deleted_by":      bs.DeletedBy,
			"deleted_time":    bs.DeletedTime,
			"gid":             bs.GID,
			"groups":          bs.Groups,
			"guid":            bs.GUID,
			"service_id":      bs.ID,
			"service_name":    bs.Name,
			"parent_srv_id":   bs.ParentSrvId,
			"parent_srv_type": bs.ParentSrvType,
			"rg_id":           bs.RGID,
			"rg_name":         bs.RGName,
			"ssh_user":        bs.SSHUser,
			"status":          bs.Status,
			"tech_status":     bs.TechStatus,
			"updated_by":      bs.UpdatedBy,
			"updated_time":    bs.UpdatedTime,
			"user_managed":    bs.UserManaged,
		}
		res = append(res, temp)
	}
	return res
}

func dataSourceBasicServiceListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	basicServiceList, err := utilityBasicServiceListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenBasicServiceList(basicServiceList))

	return nil
}

func dataSourceBasicServiceListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID of the account to query for BasicService instances",
		},
		"rg_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID of the resource group to query for BasicService instances",
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
					"base_domain": {
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
					"gid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"groups": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"guid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"service_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"service_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"parent_srv_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"parent_srv_type": {
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
					"ssh_user": {
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
				},
			},
		},
	}
	return res
}

func DataSourceBasicServiceList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceBasicServiceListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceBasicServiceListSchemaMake(),
	}
}

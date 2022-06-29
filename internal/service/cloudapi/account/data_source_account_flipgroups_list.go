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

func flattenAccountFlipGroupsList(afgl AccountFlipGroupsList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, afg := range afgl {
		temp := map[string]interface{}{
			"account_id":   afg.AccountId,
			"client_type":  afg.ClientType,
			"conn_type":    afg.ConnType,
			"created_by":   afg.CreatedBy,
			"created_time": afg.CreatedTime,
			"default_gw":   afg.DefaultGW,
			"deleted_by":   afg.DeletedBy,
			"deleted_time": afg.DeletedTime,
			"desc":         afg.Desc,
			"gid":          afg.GID,
			"guid":         afg.GUID,
			"fg_id":        afg.ID,
			"ip":           afg.IP,
			"milestones":   afg.Milestones,
			"fg_name":      afg.Name,
			"net_id":       afg.NetID,
			"net_type":     afg.NetType,
			"netmask":      afg.NetMask,
			"status":       afg.Status,
			"updated_by":   afg.UpdatedBy,
			"updated_time": afg.UpdatedTime,
		}
		res = append(res, temp)
	}
	return res

}

func dataSourceAccountFlipGroupsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountFlipGroupsList, err := utilityAccountFlipGroupsListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountFlipGroupsList(accountFlipGroupsList))

	return nil
}

func dataSourceAccountFlipGroupsListSchemaMake() map[string]*schema.Schema {
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
					"client_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"conn_type": {
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
					"default_gw": {
						Type:     schema.TypeString,
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
					"desc": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"gid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"fg_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"milestones": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"fg_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"net_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"net_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"netmask": {
						Type:     schema.TypeInt,
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

func DataSourceAccountFlipGroupsList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceAccountFlipGroupsListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceAccountFlipGroupsListSchemaMake(),
	}
}

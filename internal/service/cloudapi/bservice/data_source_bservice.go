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

func dataSourceBasicServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bs, err := utilityBasicServiceCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("account_id", bs.AccountId)
	d.Set("account_name", bs.AccountName)
	d.Set("base_domain", bs.BaseDomain)
	d.Set("computes", flattenBasicServiceComputes(bs.Computes))
	d.Set("cpu_total", bs.CPUTotal)
	d.Set("created_by", bs.CreatedBy)
	d.Set("created_time", bs.CreatedTime)
	d.Set("deleted_by", bs.DeletedBy)
	d.Set("deleted_time", bs.DeletedTime)
	d.Set("disk_total", bs.DiskTotal)
	d.Set("gid", bs.GID)
	d.Set("groups", bs.Groups)
	d.Set("groups_name", bs.GroupsName)
	d.Set("guid", bs.GUID)
	d.Set("milestones", bs.Milestones)
	d.Set("service_name", bs.Name)
	d.Set("parent_srv_id", bs.ParentSrvId)
	d.Set("parent_srv_type", bs.ParentSrvType)
	d.Set("ram_total", bs.RamTotal)
	d.Set("rg_id", bs.RGID)
	d.Set("rg_name", bs.RGName)
	d.Set("snapshots", flattenBasicServiceSnapshots(bs.Snapshots))
	d.Set("ssh_key", bs.SSHKey)
	d.Set("ssh_user", bs.SSHUser)
	d.Set("status", bs.Status)
	d.Set("tech_status", bs.TechStatus)
	d.Set("updated_by", bs.UpdatedBy)
	d.Set("updated_time", bs.UpdatedTime)
	d.Set("user_managed", bs.UserManaged)
	return nil
}

func flattenBasicServiceComputes(bscs BasicServiceComputes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bsc := range bscs {
		temp := map[string]interface{}{
			"compgroup_id":   bsc.CompGroupId,
			"compgroup_name": bsc.CompGroupName,
			"compgroup_role": bsc.CompGroupRole,
			"id":             bsc.ID,
			"name":           bsc.Name,
		}
		res = append(res, temp)
	}

	return res
}

func flattenBasicServiceSnapshots(bsrvss BasicServiceSnapshots) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bsrvs := range bsrvss {
		temp := map[string]interface{}{
			"guid":      bsrvs.GUID,
			"label":     bsrvs.Label,
			"timestamp": bsrvs.Timestamp,
			"valid":     bsrvs.Valid,
		}
		res = append(res, temp)
	}
	return res
}

func dataSourceBasicServiceSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"service_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
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
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"compgroup_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"compgroup_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"compgroup_role": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"cpu_total": {
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
		"disk_total": {
			Type:     schema.TypeString,
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
		"groups_name": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
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
		"ram_total": {
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
		"snapshots": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"timestamp": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"valid": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"ssh_key": {
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
	}
	return res
}

func DataSourceBasicService() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceBasicServiceRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceBasicServiceSchemaMake(),
	}
}

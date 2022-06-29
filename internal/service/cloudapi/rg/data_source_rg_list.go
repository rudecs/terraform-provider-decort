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

package rg

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenRgList(rgl ResgroupListResp) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rg := range rgl {
		temp := map[string]interface{}{
			"account_id":        rg.AccountID,
			"account_name":      rg.AccountName,
			"acl":               flattenRgAcl(rg.ACLs),
			"created_by":        rg.CreatedBy,
			"created_time":      rg.CreatedTime,
			"def_net_id":        rg.DefaultNetID,
			"def_net_type":      rg.DefaultNetType,
			"deleted_by":        rg.DeletedBy,
			"deleted_time":      rg.DeletedTime,
			"desc":              rg.Decsription,
			"gid":               rg.GridID,
			"guid":              rg.GUID,
			"rg_id":             rg.ID,
			"lock_status":       rg.LockStatus,
			"milestones":        rg.Milestones,
			"name":              rg.Name,
			"register_computes": rg.RegisterComputes,
			"resource_limits":   flattenRgResourceLimits(rg.ResourceLimits),
			"secret":            rg.Secret,
			"status":            rg.Status,
			"updated_by":        rg.UpdatedBy,
			"updated_time":      rg.UpdatedTime,
			"vins":              rg.Vins,
			"vms":               rg.Computes,
		}
		res = append(res, temp)
	}
	return res

}

func flattenRgAcl(rgAcls []AccountAclRecord) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rgAcl := range rgAcls {
		temp := map[string]interface{}{
			"explicit":      rgAcl.IsExplicit,
			"guid":          rgAcl.Guid,
			"right":         rgAcl.Rights,
			"status":        rgAcl.Status,
			"type":          rgAcl.Type,
			"user_group_id": rgAcl.UgroupID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenRgResourceLimits(rl ResourceLimits) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cu_c":      rl.CUC,
		"cu_d":      rl.CUD,
		"cu_i":      rl.CUI,
		"cu_m":      rl.CUM,
		"cu_np":     rl.CUNP,
		"gpu_units": rl.GpuUnits,
	}
	res = append(res, temp)

	return res

}

func dataSourceRgListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rgList, err := utilityRgListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenRgList(rgList))

	return nil
}

func dataSourceRgListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"includedeleted": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "included deleted resource groups",
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
					"acl": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"explicit": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"guid": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"right": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"status": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"user_group_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"created_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"def_net_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"def_net_type": {
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
					"rg_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"lock_status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"milestones": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"register_computes": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"resource_limits": {
						Type:     schema.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cu_c": {
									Type:     schema.TypeFloat,
									Computed: true,
								},
								"cu_d": {
									Type:     schema.TypeFloat,
									Computed: true,
								},
								"cu_i": {
									Type:     schema.TypeFloat,
									Computed: true,
								},
								"cu_m": {
									Type:     schema.TypeFloat,
									Computed: true,
								},
								"cu_np": {
									Type:     schema.TypeFloat,
									Computed: true,
								},
								"gpu_units": {
									Type:     schema.TypeFloat,
									Computed: true,
								},
							},
						},
					},
					"secret": {
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
					"vins": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"vms": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
				},
			},
		},
	}
	return res
}

func DataSourceRgList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgListSchemaMake(),
	}
}

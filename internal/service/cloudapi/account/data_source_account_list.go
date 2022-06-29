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

func flattenAccountList(al AccountCloudApiList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, acc := range al {
		temp := map[string]interface{}{
			"acl":          flattenRgAcl(acc.Acl),
			"created_time": acc.CreatedTime,
			"deleted_time": acc.DeletedTime,
			"account_id":   acc.ID,
			"account_name": acc.Name,
			"status":       acc.Status,
			"updated_time": acc.UpdatedTime,
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

/*uncomment for cloudbroker
func flattenAccountList(al AccountList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, acc := range al {
		temp := map[string]interface{}{
			"dc_location":        acc.DCLocation,
			"ckey":               acc.CKey,
			"meta":               flattenMeta(acc.Meta),

			"acl": flattenRgAcl(acc.Acl),

			"company":            acc.Company,
			"companyurl":         acc.CompanyUrl,
			"created_by":         acc.CreatedBy,

			"created_time": acc.CreatedTime,

			"deactivation_time":  acc.DeactiovationTime,
			"deleted_by":         acc.DeletedBy,

			"deleted_time": acc.DeletedTime,

			"displayname":        acc.DisplayName,
			"guid":               acc.GUID,

			"account_id":   acc.ID,
			"account_name": acc.Name,

			"resource_limits":    flattenRgResourceLimits(acc.ResourceLimits),
			"send_access_emails": acc.SendAccessEmails,
			"service_account":    acc.ServiceAccount,

			"status":       acc.Status,
			"updated_time": acc.UpdatedTime,

			"version":            acc.Version,
			"vins":               acc.Vins,

		}
		res = append(res, temp)
	}
	return res
}
*/

func dataSourceAccountListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountList, err := utilityAccountListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountList(accountList))

	return nil
}

func dataSourceAccountListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
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
					/*uncomment for cloudbroker
					"dc_location": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ckey": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"meta": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},*/
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
					/*uncomment for cloudbroker
					"company": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"companyurl": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					*/
					"created_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					/*uncomment for cloudbroker
					"deactivation_time": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"deleted_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					*/
					"deleted_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					/*uncomment for cloudbroker
					"displayname": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					*/
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					/*uncomment for cloudbroker
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
					"send_access_emails": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"service_account": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					*/
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					/*uncomment for cloudbroker
					"version": {
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
					*/
				},
			},
		},
	}
	return res
}

func DataSourceAccountList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceAccountListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceAccountListSchemaMake(),
	}
}

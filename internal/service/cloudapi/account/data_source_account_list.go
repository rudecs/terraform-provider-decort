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

func dataSourceAccountListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountList, err := utilityAccountListCheckPresence(ctx, d, m)
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
					"created_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"deleted_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
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

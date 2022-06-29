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

func flattenAccountAuditsList(aal AccountAuditsList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, aa := range aal {
		temp := map[string]interface{}{
			"call":         aa.Call,
			"responsetime": aa.ResponseTime,
			"statuscode":   aa.StatusCode,
			"timestamp":    aa.Timestamp,
			"user":         aa.User,
		}
		res = append(res, temp)
	}
	return res

}

func dataSourceAccountAuditsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountAuditsList, err := utilityAccountAuditsListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountAuditsList(accountAuditsList))

	return nil
}

func dataSourceAccountAuditsListSchemaMake() map[string]*schema.Schema {
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
					"call": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"responsetime": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"statuscode": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"timestamp": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"user": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func DataSourceAccountAuditsList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceAccountAuditsListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceAccountAuditsListSchemaMake(),
	}
}

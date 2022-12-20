/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

func dataSourceVinsAuditsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	audits, err := utilityVinsAuditsCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenVinsAudits(audits))
	return nil
}

func DataSourceVinsAuditsSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"vins_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Unique ID of the ViNS. If ViNS ID is specified, then ViNS name, rg_id and account_id are ignored.",
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"call": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"response_time": {
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
	return rets
}

func DataSourceVinsAudits() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceVinsAuditsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},
		Schema: DataSourceVinsAuditsSchemaMake(),
	}
}

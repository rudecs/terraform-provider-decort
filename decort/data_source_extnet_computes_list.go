/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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

func flattenExtnetsComputes(ecs ExtnetExtendList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, ec := range ecs {
		temp := map[string]interface{}{
			"net_id": ec.ID,
			"ipaddr": ec.IPAddr,
			"ipcidr": ec.IPCidr,
			"name":   ec.Name,
		}
		res = append(res, temp)
	}
	return res
}

func flattenExtnetComputesList(ecl ExtnetComputesList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, ec := range ecl {
		temp := map[string]interface{}{
			"account_id":   ec.AccountId,
			"account_name": ec.AccountName,
			"extnets":      flattenExtnetsComputes(ec.Extnets),
			"id":           ec.ID,
			"name":         ec.Name,
			"rg_id":        ec.RGID,
			"rg_name":      ec.RGName,
		}
		res = append(res, temp)
	}
	return res
}

func dataSourceExtnetComputesListRead(d *schema.ResourceData, m interface{}) error {
	extnetComputesList, err := utilityExtnetComputesListCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenExtnetComputesList(extnetComputesList))

	return nil
}

func dataSourceExtnetComputesListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "filter by account ID",
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
					"extnets": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"net_id": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"ipaddr": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"ipcidr": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
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
				},
			},
		},
	}
	return res
}

func dataSourceExtnetComputesList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceExtnetComputesListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceExtnetComputesListSchemaMake(),
	}
}

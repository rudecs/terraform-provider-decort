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

func flattenAccountTemplatesList(atl AccountTemplatesList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, at := range atl {
		temp := map[string]interface{}{
			"unc_path":      at.UNCPath,
			"account_id":    at.AccountId,
			"desc":          at.Desc,
			"template_id":   at.ID,
			"template_name": at.Name,
			"public":        at.Public,
			"size":          at.Size,
			"status":        at.Status,
			"type":          at.Type,
			"username":      at.Username,
		}
		res = append(res, temp)
	}
	return res

}

func dataSourceAccountTemplatesListRead(d *schema.ResourceData, m interface{}) error {
	accountTemplatesList, err := utilityAccountTemplatesListCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountTemplatesList(accountTemplatesList))

	return nil
}

func dataSourceAccountTemplatesListSchemaMake() map[string]*schema.Schema {
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
					"unc_path": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"desc": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"template_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"template_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"public": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"size": {
						Type:     schema.TypeInt,
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
					"username": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func dataSourceAccountTemplatessList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceAccountTemplatesListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceAccountTemplatesListSchemaMake(),
	}
}

/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>, <svs1370@gmail.com>

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

func flattenSepList(sl SepList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range sl {
		temp := map[string]interface{}{
			"ckey":        item.Ckey,
			"meta":        flattenMeta(item.Meta),
			"consumed_by": item.ConsumedBy,
			"desc":        item.Desc,
			"gid":         item.Gid,
			"guid":        item.Guid,
			"sep_id":      item.Id,
			"milestones":  item.Milestones,
			"name":        item.Name,
			"obj_status":  item.ObjStatus,
			"provided_by": item.ProvidedBy,
			"tech_status": item.TechStatus,
			"type":        item.Type,
		}

		res = append(res, temp)
	}
	return res
}

func dataSourceSepListRead(d *schema.ResourceData, m interface{}) error {
	sepList, err := utilitySepListCheckPresence(d, m)
	if err != nil {
		return err
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenSepList(sepList))

	return nil
}

func dataSourceSepDesSchemaMake(sh map[string]*schema.Schema) map[string]*schema.Schema {
	sh["config"] = &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"api_ips": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"protocol": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"decs3o_app_secret": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"format": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"edgeuser_password": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"edgeuser_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"decs3o_app_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"transport": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}

	return sh
}

func dataSourceSepCommonSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sep_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "sep type des id",
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
		},
		"consumed_by": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
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
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"obj_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provided_by": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceSepDes() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepDesRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepDesSchemaMake(),
	}
}

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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func flattenGridList(gl GridList) []map[string]interface{} {
	res := make([]map[string]interface{}, len(gl), len(gl))
	for _, item := range gl {
		temp := map[string]interface{}{}
		temp["name"] = item.Name
		temp["flag"] = item.Flag
		temp["gid"] = item.Gid
		temp["guid"] = item.Guid
		temp["location_code"] = item.LocationCode
		temp["id"] = item.Id
		res = append(res, temp)
	}
	return res
}

func dataSourceGridListRead(d *schema.ResourceData, m interface{}) error {
	gridList, err := utilityGridListCheckPresence(d, m)
	if err != nil {
		return err
	}
	d.SetId("1234")
	d.Set("items", flattenGridList(gridList))

	return nil
}

func dataSourceGridListSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"page": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "page number",
		},
		"size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "page size",
		},
		"items": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "grid list",
			Elem: &schema.Resource{
				Schema: dataSourceGridSchemaMake(),
			},
		},
	}

	return rets
}

func dataSourceGridSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"flag": {
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
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"location_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceGridList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceGridListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceGridListSchemaMake(),
	}
}

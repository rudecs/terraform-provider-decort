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

package sep

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/flattens"
)

func flattenSepList(sl SepList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range sl {
		data, _ := json.Marshal(item.Config)
		temp := map[string]interface{}{
			"ckey":        item.Ckey,
			"meta":        flattens.FlattenMeta(item.Meta),
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
			"config":      string(data),
		}

		res = append(res, temp)
	}
	return res
}

func dataSourceSepListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sepList, err := utilitySepListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenSepList(sepList))

	return nil
}

func dataSourceSepListSchemaMake() map[string]*schema.Schema {
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
			Description: "sep list",
			Elem: &schema.Resource{
				Schema: dataSourceSepShortSchemaMake(),
			},
		},
	}

	return rets
}

func dataSourceSepShortSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"sep_id": {
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
		"config": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceSepList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceSepListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceSepListSchemaMake(),
	}
}

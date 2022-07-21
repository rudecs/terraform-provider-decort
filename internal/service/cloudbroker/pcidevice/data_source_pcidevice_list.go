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

package pcidevice

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/flattens"
)

func flattenPcideviceList(pl PcideviceList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range pl {
		temp := map[string]interface{}{
			"ckey":        item.CKey,
			"meta":        flattens.FlattenMeta(item.Meta),
			"compute_id":  item.Computeid,
			"description": item.Description,
			"guid":        item.Guid,
			"hw_path":     item.HwPath,
			"device_id":   item.ID,
			"rg_id":       item.RgID,
			"name":        item.Name,
			"stack_id":    item.StackID,
			"status":      item.Status,
			"system_name": item.SystemName,
		}
		res = append(res, temp)
	}
	return res
}

func dataSourcePcideviceListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pcideviceList, err := utilityPcideviceListCheckPresence(ctx, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("items", flattenPcideviceList(pcideviceList))

	id := uuid.New()
	d.SetId(id.String())

	return nil
}

func dataSourcePcideviceItem() map[string]*schema.Schema {
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
		"compute_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"hw_path": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"device_id": {
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
		"stack_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourcePcideviceListSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"items": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "pcidevice list",
			Elem: &schema.Resource{
				Schema: dataSourcePcideviceItem(),
			},
		},
	}

	return rets
}

func DataSourcePcideviceList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourcePcideviceListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourcePcideviceListSchemaMake(),
	}
}

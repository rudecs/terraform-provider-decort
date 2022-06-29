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
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/flattens"
)

func dataSourcePcideviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pcidevice, err := utilityPcideviceCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("ckey", pcidevice.CKey)
	d.Set("meta", flattens.FlattenMeta(pcidevice.Meta))
	d.Set("compute_id", pcidevice.Computeid)
	d.Set("description", pcidevice.Description)
	d.Set("guid", pcidevice.Guid)
	d.Set("hw_path", pcidevice.HwPath)
	d.Set("rg_id", pcidevice.RgID)
	d.Set("name", pcidevice.Name)
	d.Set("stack_id", pcidevice.StackID)
	d.Set("status", pcidevice.Status)
	d.Set("system_name", pcidevice.SystemName)

	d.SetId(strconv.Itoa(d.Get("device_id").(int)))

	return nil
}

func dataSourcePcideviceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"device_id": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
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

	return rets
}

func DataSourcePcidevice() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourcePcideviceRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourcePcideviceSchemaMake(),
	}
}

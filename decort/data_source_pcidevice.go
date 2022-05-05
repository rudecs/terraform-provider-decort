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
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePcideviceRead(d *schema.ResourceData, m interface{}) error {
	pcidevice, err := utilityPcideviceCheckPresence(d, m)
	if err != nil {
		return err
	}

	d.Set("ckey", pcidevice.CKey)
	d.Set("meta", flattenMeta(pcidevice.Meta))
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

func dataSourcePcidevice() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourcePcideviceRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourcePcideviceSchemaMake(),
	}
}

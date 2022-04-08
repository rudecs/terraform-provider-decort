/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Petr Krutov, <petr.krutov@digitalenergy.online>

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

func dataSourceVGPURead(d *schema.ResourceData, m interface{}) error {
	vgpu, err := utilityVGPUCheckPresence(d, m)
	if vgpu == nil {
		d.SetId("")
		return err
	}

	d.SetId(strconv.Itoa(vgpu.ID))
	d.Set("vgpu_id", vgpu.ID)
	d.Set("account_id", vgpu.AccountID)
	d.Set("mode", vgpu.Mode)
	d.Set("pgpu", vgpu.PgpuID)
	d.Set("profile_id", vgpu.ProfileID)
	d.Set("ram", vgpu.RAM)
	d.Set("status", vgpu.Status)
	d.Set("type", vgpu.Type)
	d.Set("vm_id", vgpu.VmID)

	return nil
}

func dataSourceVGPUSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vgpu_id": {
			Type:     schema.TypeInt,
			Required: true,
		},

		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"mode": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"pgpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"profile_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"ram": {
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

		"vm_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceVGPU() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceVGPURead,

		Schema: dataSourceVGPUSchemaMake(),
	}
}

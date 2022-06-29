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

package vgpu

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVGPURead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vgpu, err := utilityVGPUCheckPresence(d, m)
	if vgpu == nil {
		d.SetId("")
		return diag.FromErr(err)
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

func DataSourceVGPU() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceVGPURead,

		Schema: dataSourceVGPUSchemaMake(),
	}
}

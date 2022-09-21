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

package lb

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dsLBSchemaMake() map[string]*schema.Schema {
	sch := createLBSchema()
	sch["lb_id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	return sch
}

func dsLBListDeletedSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"page": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"size": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dsLBItemSchemaMake(),
			},
		},
	}
}

func dsLBListSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"includedeleted": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"page": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"size": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dsLBItemSchemaMake(),
			},
		},
	}
}

func dsLBItemSchemaMake() map[string]*schema.Schema {
	sch := createLBSchema()
	sch["dp_api_password"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	return sch
}

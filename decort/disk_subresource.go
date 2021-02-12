/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

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

package decort

import (
	// log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)


// ID, type,  name, size, account ID, SEP ID, SEP type, pool, status, tech status, compute ID, image ID
func diskSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of this disk.",
		},

		"size": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Size of the disk in GB.",
		},

		"account_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the account this disk belongs to.",
		},

		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Type of this disk.",
		},

		"sep_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			Description: "ID of the storage end-point provider serving this disk.",
		},

		"sep_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			Description: "Type of the storage provider serving this disk.",
		},

		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			Description: "Pool on the storage where this disk is located.",
		},

		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the binary Image this disk resource is cloned from (if any).",
		},
	}

	return rets
}

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

package image

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceImageSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Owner account id",
		},
		"architecture": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image architecture",
		},
		"boot_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Boot image type",
		},
		"bootable": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Flag, true if image is bootable, otherwise - false",
		},
		"cdrom": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Flag, true if image is cdrom image, otherwise - false",
		},
		"desc": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image description",
		},
		"drivers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Image drivers",
		},
		"hot_resize": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Flag, true if image supports hot resize, else if not",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Image id",
		},
		"link_to": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "For virtual images, id image, which current image linked",
		},
		"image_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image name",
		},
		"pool_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image pool",
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Image storage endpoint id",
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Image size",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image status",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image type",
		},
		"username": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "username",
		},
		"virtual": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "True if image is virtula, otherwise - else",
		},
	}
}

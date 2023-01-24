/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceImageSchemaMake(sch map[string]*schema.Schema) map[string]*schema.Schema {
	delete(sch, "show_all")
	sch["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the rescue disk",
	}

	sch["url"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "URL where to download media from",
	}

	sch["gid"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "grid (platform) ID where this template should be create in",
	}

	sch["image_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "image id",
	}

	sch["boot_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"bios", "uefi"}, true),
		Description:  "Boot type of image bios or uefi",
	}

	sch["type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"linux", "windows", "other"}, true),
		Description:  "Image type linux, windows or other",
	}

	sch["hot_resize"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Does this machine supports hot resize",
	}

	sch["username"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Optional username for the image",
	}

	sch["password"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Optional password for the image",
	}

	sch["account_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "AccountId to make the image exclusive",
	}

	sch["username_dl"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "username for upload binary media",
	}

	sch["password_dl"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "password for upload binary media",
	}

	sch["pool_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "pool for image create",
	}

	sch["sep_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "storage endpoint provider ID",
	}

	sch["architecture"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.StringInSlice([]string{"X86_64", "PPC64_LE"}, true),
		Description:  "binary architecture of this image, one of X86_64 of PPC64_LE",
	}

	sch["drivers"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	sch["permanently"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "whether to completely delete the image",
	}

	return sch
}

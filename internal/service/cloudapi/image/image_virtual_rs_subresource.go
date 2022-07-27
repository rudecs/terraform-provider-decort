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

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceImageVirtualSchemaMake(sch map[string]*schema.Schema) map[string]*schema.Schema {
	delete(sch, "show_all")
	sch["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the rescue disk",
	}

	sch["link_to"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "ID of real image to link this virtual image to upon creation",
	}

	sch["permanently"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "whether to completely delete the image",
	}

	sch["image_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Image id",
	}

	return sch
}

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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func flattenImageList(il ImageList) []map[string]interface{} {
	res := make([]map[string]interface{}, len(il), len(il))
	for _, item := range il {
		temp := map[string]interface{}{}
		temp["name"] = item.Name
		temp["url"] = item.Url
		temp["gid"] = item.Gid
		temp["drivers"] = item.Drivers
		temp["image_id"] = item.ImageId
		temp["boot_type"] = item.Boottype
		temp["image_type"] = item.Imagetype
		res = append(res, temp)
	}
	return res
}

func dataSourceImageListRead(d *schema.ResourceData, m interface{}) error {
	imageList, err := utilityImageListCheckPresence(d, m)
	if err != nil {
		return err
	}
	d.SetId("1234")
	d.Set("items", flattenImageList(imageList))

	return nil
}

func dataSourceImageListSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "filter images by storage endpoint provider ID",
		},
		"shared_with": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "filter images by account ID availability",
		},
		"page": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "page number",
		},
		"size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "page size",
		},
		"items": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "image list",
			Elem: &schema.Resource{
				Schema: resourceImageSchemaMake(),
			},
		},
	}

	return rets
}

func dataSourceImageList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceImageListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceImageListSchemaMake(),
	}
}

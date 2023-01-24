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
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenImageList(il ImageList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, img := range il {
		temp := map[string]interface{}{
			"account_id":   img.AccountId,
			"architecture": img.Architecture,
			"boot_type":    img.BootType,
			"bootable":     img.Bootable,
			"cdrom":        img.CDROM,
			"desc":         img.Description,
			"drivers":      img.Drivers,
			"hot_resize":   img.HotResize,
			"image_id":     img.Id,
			"link_to":      img.LinkTo,
			"image_name":   img.Name,
			"pool_name":    img.Pool,
			"sep_id":       img.SepId,
			"size":         img.Size,
			"status":       img.Status,
			"type":         img.Type,
			"username":     img.Username,
			"virtual":      img.Virtual,
		}
		res = append(res, temp)
	}
	return res
}

func dataSourceImageListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	imageList, err := utilityImageListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenImageList(imageList))

	return nil
}

func dataSourceImageListSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "optional account ID to include account images",
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
				Schema: dataSourceImageSchemaMake(),
			},
		},
	}

	return rets
}

func DataSourceImageList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceImageListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceImageListSchemaMake(),
	}
}

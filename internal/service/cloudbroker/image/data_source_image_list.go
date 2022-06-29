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
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenImageList(il ImageList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range il {
		temp := map[string]interface{}{
			"name":           item.Name,
			"url":            item.Url,
			"gid":            item.Gid,
			"guid":           item.Guid,
			"drivers":        item.Drivers,
			"image_id":       item.ImageId,
			"boot_type":      item.Boottype,
			"bootable":       item.Bootable,
			"image_type":     item.Imagetype,
			"status":         item.Status,
			"tech_status":    item.TechStatus,
			"version":        item.Version,
			"username":       item.Username,
			"username_dl":    item.UsernameDL,
			"password":       item.Password,
			"password_dl":    item.PasswordDL,
			"purge_attempts": item.PurgeAttempts,
			"architecture":   item.Architecture,
			"account_id":     item.AccountId,
			"computeci_id":   item.ComputeciId,
			"enabled":        item.Enabled,
			"reference_id":   item.ReferenceId,
			"res_id":         item.ResId,
			"res_name":       item.ResName,
			"rescuecd":       item.Rescuecd,
			"provider_name":  item.ProviderName,
			"milestones":     item.Milestones,
			"size":           item.Size,
			"sep_id":         item.SepId,
			"link_to":        item.LinkTo,
			"unc_path":       item.UNCPath,
			"pool_name":      item.PoolName,
			"hot_resize":     item.Hotresize,
			"history":        flattenHistory(item.History),
			"last_modified":  item.LastModified,
			"meta":           flattenMeta(item.Meta),
			"desc":           item.Desc,
			"shared_with":    item.SharedWith,
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

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

func flattenHistory(history []History) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, item := range history {
		t := map[string]interface{}{
			"id":        item.Id,
			"guid":      item.Guid,
			"timestamp": item.Timestamp,
		}

		temp = append(temp, t)
	}
	return temp
}

func flattenImage(d *schema.ResourceData, img *ImageExtend) {
	d.Set("unc_path", img.UNCPath)
	d.Set("ckey", img.CKey)
	d.Set("account_id", img.AccountId)
	d.Set("acl", img.Acl)
	d.Set("architecture", img.Architecture)
	d.Set("boot_type", img.BootType)
	d.Set("bootable", img.Bootable)
	d.Set("compute_ci_id", img.ComputeCiId)
	d.Set("deleted_time", img.DeletedTime)
	d.Set("desc", img.Description)
	d.Set("drivers", img.Drivers)
	d.Set("enabled", img.Enabled)
	d.Set("gid", img.GridId)
	d.Set("guid", img.GUID)
	d.Set("history", flattenHistory(img.History))
	d.Set("hot_resize", img.HotResize)
	d.Set("image_id", img.Id)
	d.Set("last_modified", img.LastModified)
	d.Set("link_to", img.LinkTo)
	d.Set("milestones", img.Milestones)
	d.Set("image_name", img.Name)
	d.Set("password", img.Password)
	d.Set("pool_name", img.Pool)
	d.Set("provider_name", img.ProviderName)
	d.Set("purge_attempts", img.PurgeAttempts)
	d.Set("res_id", img.ResId)
	d.Set("rescuecd", img.RescueCD)
	d.Set("sep_id", img.SepId)
	d.Set("shared_with", img.SharedWith)
	d.Set("size", img.Size)
	d.Set("status", img.Status)
	d.Set("tech_status", img.TechStatus)
	d.Set("type", img.Type)
	d.Set("username", img.Username)
	d.Set("version", img.Version)
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	image, err := utilityImageCheckPresence(ctx, d, m)
	if err != nil {

		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())

	flattenImage(d, image)

	return nil
}

func DataSourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceImageRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceImageExtendSchemaMake(),
	}
}

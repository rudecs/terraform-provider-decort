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

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenLBList(lbl LBList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(lbl))
	for _, lb := range lbl {
		temp := map[string]interface{}{
			"ha_mode":         lb.HAMode,
			"backends":        flattenLBBackends(lb.Backends),
			"created_by":      lb.CreatedBy,
			"created_time":    lb.CreatedTime,
			"deleted_by":      lb.DeletedBy,
			"deleted_time":    lb.DeletedTime,
			"desc":            lb.Description,
			"dp_api_user":     lb.DPAPIUser,
			"dp_api_password": lb.DPAPIPassword,
			"extnet_id":       lb.ExtnetId,
			"frontends":       flattenFrontends(lb.Frontends),
			"gid":             lb.GID,
			"guid":            lb.GUID,
			"image_id":        lb.ImageId,
			"milestones":      lb.Milestones,
			"name":            lb.Name,
			"primary_node":    flattenNode(lb.PrimaryNode),
			"rg_id":           lb.RGID,
			"rg_name":         lb.RGName,
			"secondary_node":  flattenNode(lb.SecondaryNode),
			"status":          lb.Status,
			"tech_status":     lb.TechStatus,
			"updated_by":      lb.UpdatedBy,
			"updated_time":    lb.UpdatedTime,
			"vins_id":         lb.VinsId,
		}
		res = append(res, temp)
	}
	return res
}

func dataSourceLBListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbList, err := utilityLBListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenLBList(lbList))

	return nil
}

func DataSourceLBList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceLBListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dsLBListSchemaMake(),
	}
}

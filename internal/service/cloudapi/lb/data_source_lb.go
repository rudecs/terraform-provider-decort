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
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenLB(d *schema.ResourceData, lb *LoadBalancer) {
	d.Set("ha_mode", lb.HAMode)
	d.Set("backends", flattenLBBackends(lb.Backends))
	d.Set("created_by", lb.CreatedBy)
	d.Set("created_time", lb.CreatedTime)
	d.Set("deleted_by", lb.DeletedBy)
	d.Set("deleted_time", lb.DeletedTime)
	d.Set("desc", lb.Description)
	d.Set("dp_api_user", lb.DPAPIUser)
	d.Set("extnet_id", lb.ExtnetId)
	d.Set("frontends", flattenFrontends(lb.Frontends))
	d.Set("gid", lb.GID)
	d.Set("guid", lb.GUID)
	d.Set("image_id", lb.ImageId)
	d.Set("milestones", lb.Milestones)
	d.Set("name", lb.Name)
	d.Set("primary_node", flattenNode(lb.PrimaryNode))
	d.Set("rg_id", lb.RGID)
	d.Set("rg_name", lb.RGName)
	d.Set("secondary_node", flattenNode(lb.SecondaryNode))
	d.Set("status", lb.Status)
	d.Set("tech_status", lb.TechStatus)
	d.Set("updated_by", lb.UpdatedBy)
	d.Set("updated_time", lb.UpdatedTime)
	d.Set("vins_id", lb.VinsId)
}

func dataSourceLBRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lb, err := utilityLBCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatUint(lb.ID, 10))

	flattenLB(d, lb)

	return nil
}

func DataSourceLB() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceLBRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dsLBSchemaMake(),
	}
}

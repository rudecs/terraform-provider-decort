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

func flattenImageListStacks(_ *schema.ResourceData, stack ImageListStacks) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, item := range stack {
		t := map[string]interface{}{
			"api_url":      item.ApiURL,
			"api_key":      item.ApiKey,
			"app_id":       item.AppId,
			"desc":         item.Desc,
			"drivers":      item.Drivers,
			"error":        item.Error,
			"guid":         item.Guid,
			"id":           item.Id,
			"images":       item.Images,
			"login":        item.Login,
			"name":         item.Name,
			"passwd":       item.Passwd,
			"reference_id": item.ReferenceId,
			"status":       item.Status,
			"type":         item.Type,
		}

		temp = append(temp, t)
	}
	return temp
}

func dataSourceImageListStacksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	imageListStacks, err := utilityImageListStacksCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenImageListStacks(d, imageListStacks))

	return nil
}

func dataSourceImageListStackSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"api_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"app_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"drivers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"error": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"images": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"login": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"passwd": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceImageListStacksSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "image id",
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
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataSourceImageListStackSchemaMake(),
			},
			Description: "items of stacks list",
		},
	}
}

func DataSourceImageListStacks() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceImageListStacksRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceImageListStacksSchemaMake(),
	}
}

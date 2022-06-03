/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBasicServiceGroupRead(d *schema.ResourceData, m interface{}) error {
	bsg, err := utilityBasicServiceGroupCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("account_id", bsg.AccountId)
	d.Set("account_name", bsg.AccountName)
	d.Set("computes", flattenBSGroupComputes(bsg.Computes))
	d.Set("consistency", bsg.Consistency)
	d.Set("cpu", bsg.CPU)
	d.Set("created_by", bsg.CreatedBy)
	d.Set("created_time", bsg.CreatedTime)
	d.Set("deleted_by", bsg.DeletedBy)
	d.Set("deleted_time", bsg.DeletedTime)
	d.Set("disk", bsg.Disk)
	d.Set("driver", bsg.Driver)
	d.Set("extnets", bsg.Extnets)
	d.Set("gid", bsg.GID)
	d.Set("guid", bsg.GUID)
	d.Set("image_id", bsg.ImageId)
	d.Set("milestones", bsg.Milestones)
	d.Set("compgroup_name", bsg.Name)
	d.Set("parents", bsg.Parents)
	d.Set("ram", bsg.RAM)
	d.Set("rg_id", bsg.RGID)
	d.Set("rg_name", bsg.RGName)
	d.Set("role", bsg.Role)
	d.Set("sep_id", bsg.SepId)
	d.Set("seq_no", bsg.SeqNo)
	d.Set("status", bsg.Status)
	d.Set("tech_status", bsg.TechStatus)
	d.Set("timeout_start", bsg.TimeoutStart)
	d.Set("updated_by", bsg.UpdatedBy)
	d.Set("updated_time", bsg.UpdatedTime)
	d.Set("vinses", bsg.Vinses)
	return nil
}

func flattenBSGroupOSUsers(bsgosus BasicServiceGroupOSUsers) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bsgosu := range bsgosus {
		temp := map[string]interface{}{
			"login":    bsgosu.Login,
			"password": bsgosu.Password,
		}
		res = append(res, temp)
	}

	return res
}

func flattenBSGroupComputes(bsgcs BasicServiceGroupComputes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bsgc := range bsgcs {
		temp := map[string]interface{}{
			"id":           bsgc.ID,
			"ip_addresses": bsgc.IPAdresses,
			"name":         bsgc.Name,
			"os_users":     flattenBSGroupOSUsers(bsgc.OSUsers),
		}
		res = append(res, temp)
	}
	return res
}

func dataSourceBasicServiceGroupSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"service_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"compgroup_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"account_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"ip_addresses": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"os_users": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"login": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"password": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"consistency": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"cpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"created_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deleted_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"driver": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"extnets": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"compgroup_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"parents": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"ram": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rg_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rg_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"role": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sep_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"seq_no": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"timeout_start": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"updated_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vinses": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}
	return res
}

func dataSourceBasicServiceGroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceBasicServiceGroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceBasicServiceGroupSchemaMake(),
	}
}

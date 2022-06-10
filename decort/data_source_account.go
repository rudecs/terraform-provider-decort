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

func dataSourceAccountRead(d *schema.ResourceData, m interface{}) error {
	acc, err := utilityAccountCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("dc_location", acc.DCLocation)
	d.Set("resources", flattenAccResources(acc.Resources))
	d.Set("ckey", acc.CKey)
	d.Set("meta", flattenMeta(acc.Meta))
	d.Set("acl", flattenAccAcl(acc.Acl))
	d.Set("company", acc.Company)
	d.Set("companyurl", acc.CompanyUrl)
	d.Set("created_by", acc.CreatedBy)
	d.Set("created_time", acc.CreatedTime)
	d.Set("deactivation_time", acc.DeactiovationTime)
	d.Set("deleted_by", acc.DeletedBy)
	d.Set("deleted_time", acc.DeletedTime)
	d.Set("displayname", acc.DisplayName)
	d.Set("guid", acc.GUID)
	d.Set("account_id", acc.ID)
	d.Set("account_name", acc.Name)
	d.Set("resource_limits", flattenRgResourceLimits(acc.ResourceLimits))
	d.Set("send_access_emails", acc.SendAccessEmails)
	d.Set("service_account", acc.ServiceAccount)
	d.Set("status", acc.Status)
	d.Set("updated_time", acc.UpdatedTime)
	d.Set("version", acc.Version)
	d.Set("vins", acc.Vins)
	d.Set("vinses", acc.Vinses)
	d.Set("computes", flattenAccComputes(acc.Computes))
	d.Set("machines", flattenAccMachines(acc.Machines))
	return nil
}

func flattenAccComputes(acs Computes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"started": acs.Started,
		"stopped": acs.Stopped,
	}
	res = append(res, temp)
	return res
}

func flattenAccMachines(ams Machines) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"running": ams.Running,
		"halted":  ams.Halted,
	}
	res = append(res, temp)
	return res
}

func flattenAccAcl(acls []AccountAclRecord) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, acls := range acls {
		temp := map[string]interface{}{
			"can_be_deleted": acls.CanBeDeleted,
			"explicit":       acls.IsExplicit,
			"guid":           acls.Guid,
			"right":          acls.Rights,
			"status":         acls.Status,
			"type":           acls.Type,
			"user_group_id":  acls.UgroupID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenAccResources(r Resources) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"current":  flattenAccResource(r.Current),
		"reserved": flattenAccResource(r.Reserved),
	}
	res = append(res, temp)
	return res
}

func flattenAccResource(r Resource) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cpu":        r.CPU,
		"disksize":   r.Disksize,
		"extips":     r.Extips,
		"exttraffic": r.Exttraffic,
		"gpu":        r.GPU,
		"ram":        r.RAM,
	}
	res = append(res, temp)
	return res
}

func dataSourceAccountSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"account_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"dc_location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"resources": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"current": {
						Type:     schema.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"disksize": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"extips": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"exttraffic": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"gpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"ram": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"reserved": {
						Type:     schema.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"disksize": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"extips": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"exttraffic": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"gpu": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"ram": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"meta": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"can_be_deleted": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"explicit": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"right": {
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
					"user_group_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"company": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"companyurl": {
			Type:     schema.TypeString,
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
		"deactivation_time": {
			Type:     schema.TypeFloat,
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
		"displayname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"account_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"resource_limits": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cu_c": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"cu_d": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"cu_i": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"cu_m": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"cu_np": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"gpu_units": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
				},
			},
		},
		"send_access_emails": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"service_account": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vins": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"started": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"stopped": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"machines": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"halted": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"running": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"vinses": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
	return res
}

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceAccountSchemaMake(),
	}
}

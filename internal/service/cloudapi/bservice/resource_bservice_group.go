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

package bservice

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceBasicServiceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceGroupCreate")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
	urlValues.Add("name", d.Get("compgroup_name").(string))

	urlValues.Add("count", strconv.Itoa(d.Get("comp_count").(int)))
	urlValues.Add("cpu", strconv.Itoa(d.Get("cpu").(int)))
	urlValues.Add("ram", strconv.Itoa(d.Get("ram").(int)))
	urlValues.Add("disk", strconv.Itoa(d.Get("disk").(int)))
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("driver", strings.ToUpper(d.Get("driver").(string)))

	if role, ok := d.GetOk("role"); ok {
		urlValues.Add("role", role.(string))
	}

	if timeoutStart, ok := d.GetOk("timeout_start"); ok {
		urlValues.Add("timeoutStart", strconv.Itoa(timeoutStart.(int)))
	}

	if vinses, ok := d.GetOk("vinses"); ok {
		vs := vinses.([]interface{})
		temp := ""
		l := len(vs)
		for i, v := range vs {
			s := strconv.Itoa(v.(int))
			if i != (l - 1) {
				s += ","
			}
			temp = temp + s
		}
		temp = "[" + temp + "]"
		urlValues.Add("vinses", temp)
	}
	if extnets, ok := d.GetOk("extnets"); ok {
		es := extnets.([]interface{})
		temp := ""
		l := len(es)
		for i, e := range es {
			s := strconv.Itoa(e.(int))
			if i != (l - 1) {
				s += ","
			}
			temp = temp + s
		}
		temp = "[" + temp + "]"
		urlValues.Add("extnets", temp)
	}

	compgroupId, err := c.DecortAPICall(ctx, "POST", bserviceGroupAddAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(compgroupId)
	d.Set("compgroup_id", compgroupId)

	diagnostics := resourceBasicServiceGroupRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceBasicServiceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceGroupRead")

	bsg, err := utilityBasicServiceGroupCheckPresence(ctx, d, m)
	if bsg == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

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
	d.Set("compgroup_id", bsg.ID)
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

func resourceBasicServiceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceGroupDelete")

	bsg, err := utilityBasicServiceGroupCheckPresence(ctx, d, m)
	if bsg == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
	urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))

	_, err = c.DecortAPICall(ctx, "POST", bserviceGroupRemoveAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceBasicServiceGroupEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceGroupEdit")
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}

	if d.HasChange("comp_count") {
		urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
		urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))
		urlValues.Add("count", strconv.Itoa(d.Get("comp_count").(int)))
		urlValues.Add("mode", strings.ToUpper(d.Get("mode").(string)))
		_, err := c.DecortAPICall(ctx, "POST", bserviceGroupResizeAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("start") {
		api := bserviceGroupStopAPI
		start := d.Get("start").(bool)
		if start {
			api = bserviceGroupStartAPI
		} else {
			urlValues.Add("force", strconv.FormatBool(d.Get("force_stop").(bool)))
		}
		urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
		urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChanges("compgroup_name", "ram", "cpu", "disk", "role") {
		urlValues.Add("name", d.Get("compgroup_name").(string))
		urlValues.Add("cpu", strconv.Itoa(d.Get("cpu").(int)))
		urlValues.Add("ram", strconv.Itoa(d.Get("ram").(int)))
		urlValues.Add("disk", strconv.Itoa(d.Get("disk").(int)))
		urlValues.Add("role", d.Get("role").(string))
		urlValues.Add("force", strconv.FormatBool(d.Get("force_update").(bool)))

		urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
		urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", bserviceGroupUpdateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("extnets") {
		extnets := d.Get("extnets").([]interface{})
		temp := ""
		l := len(extnets)
		for i, e := range extnets {
			s := strconv.Itoa(e.(int))
			if i != (l - 1) {
				s += ",\n"
			} else {
				s += "\n"
			}
			temp = temp + s
		}
		temp = "[" + temp + "]"

		urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
		urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))
		urlValues.Add("extnets", temp)
		_, err := c.DecortAPICall(ctx, "POST", bserviceGroupUpdateExtnetAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("vinses") {
		vinses := d.Get("vinses").([]interface{})
		temp := ""
		l := len(vinses)
		for i, v := range vinses {
			s := strconv.Itoa(v.(int))
			if i != (l - 1) {
				s += ",\n"
			} else {
				s += "\n"
			}
			temp = temp + s
		}
		temp = "[" + temp + "]"

		urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
		urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))
		urlValues.Add("vinses", temp)
		_, err := c.DecortAPICall(ctx, "POST", bserviceGroupUpdateVinsAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("parents") {
		deletedParents := make([]interface{}, 0)
		addedParents := make([]interface{}, 0)

		old, new := d.GetChange("parents")
		oldConv := old.([]interface{})
		newConv := new.([]interface{})
		for _, el := range oldConv {
			if !isContainsParent(newConv, el) {
				deletedParents = append(deletedParents, el)
			}
		}
		for _, el := range newConv {
			if !isContainsParent(oldConv, el) {
				addedParents = append(addedParents, el)
			}
		}

		if len(deletedParents) > 0 {
			for _, parent := range deletedParents {
				parentConv := parent.(int)

				urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
				urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))
				urlValues.Add("parentId", strconv.Itoa(parentConv))

				_, err := c.DecortAPICall(ctx, "POST", bserviceGroupParentRemoveAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

		if len(addedParents) > 0 {
			for _, parent := range addedParents {
				parentConv := parent.(int)
				urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
				urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))
				urlValues.Add("parentId", strconv.Itoa(parentConv))
				_, err := c.DecortAPICall(ctx, "POST", bserviceGroupParentAddAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}
	}

	if d.HasChange("remove_computes") {
		rcs := d.Get("remove_computes").([]interface{})
		if len(rcs) > 0 {
			for _, rc := range rcs {
				urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
				urlValues.Add("compgroupId", strconv.Itoa(d.Get("compgroup_id").(int)))
				urlValues.Add("computeId", strconv.Itoa(rc.(int)))

				_, err := c.DecortAPICall(ctx, "POST", bserviceGroupComputeRemoveAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}
	}

	return nil
}

func isContainsParent(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(int)
		elConv := el.(int)
		if elOldConv == elConv {
			return true
		}
	}
	return false
}

func resourceBasicServiceGroupSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the Basic Service to add a group to",
		},
		"compgroup_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Compute Group to add",
		},
		"comp_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "computes number. Defines how many computes must be there in the group",
		},
		"cpu": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "compute CPU number. All computes in the group have the same CPU count",
		},
		"ram": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "compute RAM volume in MB. All computes in the group have the same RAM volume",
		},
		"disk": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "compute boot disk size in GB",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "OS image ID to create computes from",
		},
		"driver": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "compute driver like a KVM_X86, KVM_PPC, etc.",
		},
		"role": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "group role tag. Can be empty string, does not have to be unique",
		},
		"timeout_start": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "time of Compute Group readiness",
		},
		"extnets": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "list of external networks to connect computes to",
		},
		"vinses": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "list of ViNSes to connect computes to",
		},
		"mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "RELATIVE",
			ValidateFunc: validation.StringInSlice([]string{"RELATIVE", "ABSOLUTE"}, false),
			Description:  "(RELATIVE;ABSOLUTE) either delta or absolute value of computes",
		},
		"start": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Start the specified Compute Group within BasicService",
		},
		"force_stop": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "force stop Compute Group",
		},
		"force_update": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "force resize Compute Group",
		},
		"parents": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"remove_computes": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"compgroup_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
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
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
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
		"updated_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func ResourceBasicServiceGroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceBasicServiceGroupCreate,
		ReadContext:   resourceBasicServiceGroupRead,
		UpdateContext: resourceBasicServiceGroupEdit,
		DeleteContext: resourceBasicServiceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout600s,
			Read:    &constants.Timeout300s,
			Update:  &constants.Timeout300s,
			Delete:  &constants.Timeout300s,
			Default: &constants.Timeout300s,
		},

		Schema: resourceBasicServiceGroupSchemaMake(),
	}
}

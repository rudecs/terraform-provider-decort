/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

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

	"encoding/json"
	"fmt"
	"log"
	// "net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func parseComputeDisks(disks []DiskRecord) []interface{} {
	length := len(disks)
	log.Debugf("parseComputeDisks: called for %d disks", length)

	result := make([]interface{}, length)
	if length == 0 {
		return result
	}

	elem := make(map[string]interface{})

	for i, value := range disks {
		// keys in this map should correspond to the Schema definition
		// as returned by dataSourceDiskSchemaMake()
		elem[" attribute "] = value. attribute
		...
		result[i] = elem
	}

	return result // this result will be used to d.Set("disks",) item of dataSourceCompute schema
}

func parseComputeInterfaces(ifaces []InterfaceRecord) []interface{} {
	length := len(ifaces)
	log.Debugf("parseComputeInterfaces: called for %d ifaces", length)

	result := make([]interface{}, length)
	if length == 0 {
		return result
	}

	elem := make(map[string]interface{})

	for i, value := range ifaces {
		// Keys in this map should correspond to the Schema definition
		// as returned by dataSourceInterfaceSchemaMake()
		elem[" attribute "] = value. attribute
		...
		result[i] = elem
	}

	return result // this result will be used to d.Set("interfaces",) item of dataSourceCompute schema
}


func flattenCompute(d *schema.ResourceData, comp_facts string) error {
	// This function expects that comp_facts string contains response from API compute/get,
	// i.e. detailed information about compute instance.
	// 
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceComputeExists(...) method
	model := ComputeGetResp{}
	log.Debugf("flattenCompute: ready to unmarshal string %q", comp_facts) 
	err := json.Unmarshal([]byte(comp_facts), &model)
	if err != nil {
		return err
	}

	log.Debugf("flattenCompute: model.ID %d, model.ResGroupID %d", model.ID, model.ResGroupID)
			   
	d.SetId(fmt.Sprintf("%d", model.ID))
	d.Set("compute_id", model.ID)
	d.Set("name", model.Name)
	d.Set("rg_id", model.ResGroupID)
	d.Set("rg_name", model.ResGroupName)
	d.Set("account_id", model.AccountID)
	d.Set("account_name", model.AccountName)
	d.Set("arch", model.Arch)
	d.Set("cpu", model.Cpu)
	d.Set("ram", model.Ram)
	d.Set("boot_disk_size", model.BootDiskSize)
	d.Set("image_id", model.ImageID)
	d.Set("description", model.Desc)
	d.Set("status", model.Status)
	d.Set("tech_status", model.TechStatus)

	if len(model.Disks) > 0 {
		log.Debugf("flattenCompute: calling parseComputeDisks for %d disks", len(model.Disks))
		if err = d.Set("disks", parseComputeDisks(model.Disks)); err != nil {
			return err
		}
	}

	if len(model.Interfaces) > 0 {
		log.Printf("flattenCompute: calling parseComputeInterfaces for %d interfaces", len(model.Interfaces))
		if err = d.Set("interfaces", parseComputeInterfaces(model.Interfaces)); err != nil {
			return err
		}
	}

	if len(model.GuestLogins) > 0 {
		log.Printf("flattenCompute: calling parseGuestLogins")
		guest_logins := parseGuestLogins(model.GuestLogins)
		if err = d.Set("guest_logins", guest_logins); err != nil {
			return err
		}

		default_login := guest_logins[0].(map[string]interface{})
		// set user & password attributes to the corresponding values of the 1st item in the list
		if err = d.Set("user", default_login["login"]); err != nil {
			return err
		}
		if err = d.Set("password", default_login["password"]); err != nil {
			return err
		}
	}

	return nil
}

func dataSourceComputeRead(d *schema.ResourceData, m interface{}) error {
	comp_facts, err := utilityComputeCheckPresence(d, m)
	if comp_facts == "" {
		// if empty string is returned from utilityComputeCheckPresence then there is no
		// such Compute and err tells so - just return it to the calling party 
		d.SetId("") // ensure ID is empty
		return err
	}

	return flattenCompute(d, comp_facts)
}

func dataSourceCompute() *schema.Resource {
	return &schema.Resource {
		SchemaVersion: 1,

		Read:   dataSourceComputeRead,

		Timeouts: &schema.ResourceTimeout {
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema {
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:  "Name of this compute instance. NOTE: this parameter is case sensitive.",
			},

			"compute_id": {
				Type:          schema.TypeInt, 
				Optional:      true,
				Description:   "ID of the compute instance. If ID is specified, name and resource group ID are ignored."
			},

			"rg_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ID of the resource group where this compute instance is located.",
			},

			"rg_name": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Name of the resource group where this compute instance is located.",
			},

			"account_id": {
				Type:         schema.TypeInt,
				Computed:     true,
				Description:  "ID of the account this compute instance belongs to.",
			},

			"account_name": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Name of the account this compute instance belongs to.",
			},

			"arch": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Hardware architecture of this compute instance.",
			},

			"cpu": {
				Type:         schema.TypeInt,
				Computed:     true,
				Description:  "Number of CPUs allocated for this compute instance.",
			},

			"ram": {
				Type:         schema.TypeInt,
				Computed:     true,
				Description:  "Amount of RAM in MB allocated for this compute instance.",
			},

			"image_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the OS image this compute instance is based on.",
			},

			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the OS image this compute instance is based on.",
			},

			"boot_disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This compute instance boot disk size in GB.",
			},

			"disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  dataSourceDiskSchemaMake(), // ID, type,  name, size, account ID, SEP ID, SEP type, pool, status, tech status, compute ID, image ID
				},
				Description: "Detailed specification for all disks attached to this compute instance (including bood disk).",
			},

			"guest_logins": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  guestLoginsSubresourceSchema(),
				},
				Description: "Details about the guest OS users provisioned together with this compute instance.",
			},

			"interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  interfaceSubresourceSchema(),
				},
				Description: "Specification for the virtual NICs configured on this compute instance.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-defined text description of this compute instance.",
			},

			"status": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Current model status of this compute instance.",
			},

			"tech_status": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Current technical status of this compute instance.",
			},

			/*
			"internal_ip": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Internal IP address of this Compute.",
			},
			*/
		},
	}
}
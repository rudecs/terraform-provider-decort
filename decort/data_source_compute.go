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

func flattenCompute(d *schema.ResourceData, comp_facts string) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceComputeExists(...) method
	model := ComputeGetResp{}
	log.Printf("flattenCompute: ready to unmarshal string %q", comp_facts) 
	err := json.Unmarshal([]byte(comp_facts), &model)
	if err != nil {
		return err
	}

	log.Printf("flattenCompute: model.ID %d, model.ResGroupID %d", model.ID, model.ResGroupID)
			   
	d.SetId(fmt.Sprintf("%d", model.ID))
	d.Set("name", model.Name)
	d.Set("rgid", model.ResGroupID)
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

	bootdisk_map := make(map[string]interface{})
	bootdisk_map["size"] = model.BootDisk
	bootdisk_map["label"] = "boot"
	bootdisk_map["pool"] = "default"
	bootdisk_map["provider"] = "default"
	
	if err = d.Set("boot_disk", []interface{}{bootdisk_map}); err != nil {
		return err
	}

	if len(model.DataDisks) > 0 {
		log.Printf("flattenCompute: calling flattenDataDisks")
		if err = d.Set("data_disks", flattenDataDisks(model.DataDisks)); err != nil {
			return err
		}
	}

	if len(model.NICs) > 0 {
		log.Printf("flattenCompute: calling flattenNICs")
		if err = d.Set("nics", flattenNICs(model.NICs)); err != nil {
			return err
		}
		log.Printf("flattenCompute: calling flattenNetworks")
		if err = d.Set("networks", flattenNetworks(model.NICs)); err != nil {
			return err
		}
	}

	if len(model.GuestLogins) > 0 {
		log.Printf("flattenCompute: calling flattenGuestLogins")
		guest_logins := flattenGuestLogins(model.GuestLogins)
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
				Required:      true,
				Description:  "Name of this compute instance. NOTE: this parameter is case sensitive.",
			},

			"rgid": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
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

			"networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  networkSubresourceSchema(),
				},
				Description: "Specification for the networks to connect this virtual machine to.",
			},

			"nics": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  nicSubresourceSchema(),
				},
				Description: "Specification for the virutal NICs allocated to this virtual machine.",
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
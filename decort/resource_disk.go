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
	// "encoding/json"
	"fmt"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDiskCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceDiskCreate: called for Disk name %q, Account ID %d", d.Get("name").(string), d.Get("account_id").(int))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	// accountId, gid, name, description, size, type, sep_id, pool
	urlValues.Add("accountId", fmt.Sprintf("%d", d.Get("account_id").(int)))
	urlValues.Add("gid", fmt.Sprintf("%d", DefaultGridID)) // we use default Grid ID, which was obtained along with DECORT Controller init
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("size", d.Get("size").(string))
	urlValues.Add("type", d.Get("type").(string))
	urlValues.Add("sep_id", fmt.Sprintf("%d", d.Get("sep_id").(int)))
	urlValues.Add("pool", d.Get("pool").(string))
	
	argVal, argSet := d.GetOk("description")
	if argSet {
		urlValues.Add("decs", argVal.(string))
	} 

	apiResp, err := controller.decortAPICall("POST", DiskCreateAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId(apiResp) // update ID of the resource to tell Terraform that the disk resource exists
	diskId, _ := strconv.Atoi(apiResp)

	log.Debugf("resourceDiskCreate: new Disk ID %d, name %q creation sequence complete", diskId, d.Get("name").(string))

	// We may reuse dataSourceDiskRead here as we maintain similarity 
	// between Disk resource and Disk data source schemas
	// Disk resource read function will also update resource ID on success, so that Terraform 
	// will know the resource exists (however, we already did it a few lines before)
	return dataSourceDiskRead(d, m)
}

func resourceDiskRead(d *schema.ResourceData, m interface{}) error {
	disk_facts, err := utilityDiskCheckPresence(d, m)
	if disk_facts == "" {
		// if empty string is returned from utilityDiskCheckPresence then there is no
		// such Disk and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty
		return err
	}

	return flattenDisk(d, disk_facts)
}

func resourceDiskUpdate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceDiskUpdate: called for disk name %q,  Account ID %d",
		d.Get("name").(string), d.Get("account_id").(int))

	log.Warn("resourceDiskUpdate: NOT IMPLEMENTED YET!")

	// we may reuse dataSourceDiskRead here as we maintain similarity 
	// between Compute resource and Compute data source schemas
	return dataSourceDiskRead(d, m) 
}

func resourceDiskDelete(d *schema.ResourceData, m interface{}) error {
	log.Warn("resourceDiskDelete: NOT IMPLEMENTED YET!")
	return nil
}

func resourceDiskExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	log.Debugf("resourceDiskExists: called for Disk name %q, Account ID %d",
		d.Get("name").(string), d.Get("account_id").(int))

	diskFacts, err := utilityDiskCheckPresence(d, m)
	if diskFacts == "" {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func resourceDiskSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of this disk. NOTE: disk names are NOT unique within an account. If disk ID is specified, disk name is ignored.",
		},

		"disk_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID of the disk to get. If disk ID is specified, then disk name and account ID are ignored.",
		},

		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID of the account this disk belongs to.",
		},

		"sep_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Storage end-point provider serving this disk. Cannot be changed for existing disk.",
		},

		"pool": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Pool where this disk is located. Cannot be changed for existing disk.",
		},

		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Size of the disk in GB. Note, that existing disks can only be grown in size.",
		},

		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "D",
			Description: "Optional type of this disk. Defaults to D, i.e. data disk. Cannot be changed for existing disks.",
		},

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional user-defined text description of this disk.",
		},

		// The rest of the attributes are all computed 
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the account this disk belongs to.",
		},

		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the image, which this disk was cloned from (if ever cloned).",
		},

		"sep_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of the storage end-point provider serving this disk.",
		},

		/*
			"snapshots": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:          &schema.Resource {
					Schema:    snapshotSubresourceSchemaMake(),
				},
				Description: "List of user-created snapshots for this disk."
			},

		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current model status of this disk.",
		},

		"tech_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current technical status of this disk.",
		},

		"compute_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the compute instance where this disk is attached to, or 0 for unattached disk.",
		},
		*/
	}

	return rets
}

func resourceDisk() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceDiskCreate,
		Read:   resourceDiskRead,
		Update: resourceDiskUpdate,
		Delete: resourceDiskDelete,
		Exists: resourceDiskExists,

		Timeouts: &schema.ResourceTimeout{
			Create:  &Timeout180s,
			Read:    &Timeout30s,
			Update:  &Timeout180s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

		Schema: resourceDiskSchemaMake(),
	}
}

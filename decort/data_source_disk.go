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
	// "net/url"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func flattenDisk(d *schema.ResourceData, disk_facts string) error {
	// This function expects disk_facts string to contain a response from disks/get API
	//
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceDiskExists(...) method. Use utilityDiskCheckPresence instead.
	model := DiskRecord{}
	log.Debugf("flattenDisk: ready to unmarshal string %q", disk_facts)
	err := json.Unmarshal([]byte(disk_facts), &model)
	if err != nil {
		return err
	}

	log.Debugf("flattenDisk: disk ID %d, disk AccountID %d", model.ID, model.AccountID)

	d.SetId(fmt.Sprintf("%d", model.ID))
	d.Set("disk_id", model.ID)
	d.Set("name", model.Name)
	d.Set("account_id", model.AccountID)
	d.Set("account_name", model.AccountName)
	d.Set("size", model.SizeMax)
	// d.Set("sizeUsed", model.SizeUsed)
	d.Set("type", model.Type)
	d.Set("image_id", model.ImageID)
	d.Set("sep_id", model.SepID)
	d.Set("sep_type", model.SepType)
	d.Set("pool", model.Pool)
	d.Set("compute_id", model.ComputeID)

	d.Set("description", model.Desc)
	d.Set("status", model.Status)
	d.Set("tech_status", model.TechStatus)

	/* we do not manage snapshots via Terraform yet, so keep this commented out for a while
	if len(model.Snapshots) > 0 {
		log.Debugf("flattenDisk: calling flattenDiskSnapshots")
		if err = d.Set("nics", flattenDiskSnapshots(model.Snapshots)); err != nil {
			return err
		}
	}
	*/

	return nil
}

func dataSourceDiskRead(d *schema.ResourceData, m interface{}) error {
	disk_facts, err := utilityDiskCheckPresence(d, m)
	if disk_facts == "" {
		// if empty string is returned from utilityDiskCheckPresence then there is no
		// such Disk and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty
		return err
	}

	return flattenDisk(d, disk_facts)
}

func dataSourceDiskSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of this disk. NOTE: disk names are NOT unique within an account.",
		},

		"disk_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID of the disk to get. If disk ID is specified, then name, account and account ID are ignored.",
		},

		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID of the account this disk belongs to.",
		},

		"account_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the account this disk belongs to. If account ID is specified, account name is ignored.",
		},

		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "User-defined text description of this disk.",
		},

		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the image, which this disk was cloned from.",
		},

		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Size of the disk in GB.",
		},

		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of this disk.",
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
		*/

		"sep_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Storage end-point provider serving this disk.",
		},

		"sep_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of the storage end-point provider serving this disk.",
		},

		"pool": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Pool where this disk is located.",
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
	}

	return rets
}

func dataSourceDisk() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceDiskRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceDiskSchemaMake(),
	}
}

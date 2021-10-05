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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDiskCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceDiskCreate: called for Disk name %q, Account ID %d", d.Get("name").(string), d.Get("account_id").(int))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	// accountId, gid, name, description, size, type, sep_id, pool
	urlValues.Add("accountId", fmt.Sprintf("%d", d.Get("account_id").(int)))
	urlValues.Add("gid", fmt.Sprintf("%d", DefaultGridID)) // we use default Grid ID, which was obtained along with DECORT Controller init
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("size", fmt.Sprintf("%d", d.Get("size").(int)))
	urlValues.Add("type", "D") // NOTE: only disks of Data type are managed via plugin
	urlValues.Add("sep_id", fmt.Sprintf("%d", d.Get("sep_id").(int)))
	urlValues.Add("pool", d.Get("pool").(string))
	
	argVal, argSet := d.GetOk("description")
	if argSet {
		urlValues.Add("description", argVal.(string))
	} 

	apiResp, err := controller.decortAPICall("POST", DisksCreateAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId(apiResp) // update ID of the resource to tell Terraform that the disk resource exists
	diskId, _ := strconv.Atoi(apiResp)

	log.Debugf("resourceDiskCreate: new Disk ID / name %d / %s creation sequence complete", diskId, d.Get("name").(string))

	// We may reuse dataSourceDiskRead here as we maintain similarity 
	// between Disk resource and Disk data source schemas
	// Disk resource read function will also update resource ID on success, so that Terraform 
	// will know the resource exists (however, we already did it a few lines before)
	return dataSourceDiskRead(d, m)
}

func resourceDiskRead(d *schema.ResourceData, m interface{}) error {
	diskFacts, err := utilityDiskCheckPresence(d, m)
	if diskFacts == "" {
		// if empty string is returned from utilityDiskCheckPresence then there is no
		// such Disk and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty
		return err
	}

	return flattenDisk(d, diskFacts)
}

func resourceDiskUpdate(d *schema.ResourceData, m interface{}) error {
	// Update will only change the following attributes of the disk:
	//  - Size; to keep data safe, shrinking disk is not allowed.
	//  - Name
	//
	// Attempt to change disk type will throw an error and mark disk
	// resource as partially updated
	log.Debugf("resourceDiskUpdate: called for Disk ID / name %s / %s,  Account ID %d",
		d.Id(), d.Get("name").(string), d.Get("account_id").(int))

	d.Partial(true)

	controller := m.(*ControllerCfg)

	oldSize, newSize := d.GetChange("size")
	if oldSize.(int) < newSize.(int) {
		log.Debugf("resourceDiskUpdate: resizing disk ID %s - %d GB -> %d GB", 
		           d.Id(), oldSize.(int), newSize.(int))
		sizeParams := &url.Values{}
		sizeParams.Add("diskId", d.Id())
		sizeParams.Add("size", fmt.Sprintf("%d", newSize.(int)))
		_, err := controller.decortAPICall("POST", DisksResizeAPI, sizeParams)
		if err != nil {
			return err
		}
		d.SetPartial("size")
	} else if oldSize.(int) > newSize.(int) {
		return fmt.Errorf("resourceDiskUpdate: Disk ID %s - reducing disk size is not allowed", d.Id())
	}

	oldName, newName := d.GetChange("name")
	if oldName.(string) != newName.(string) {
		log.Debugf("resourceDiskUpdate: renaming disk ID %d - %s -> %s", 
		           d.Get("disk_id").(int), oldName.(string), newName.(string))
		renameParams := &url.Values{}
		renameParams.Add("diskId", d.Id())
		renameParams.Add("name", newName.(string))
		_, err := controller.decortAPICall("POST", DisksRenameAPI, renameParams)
		if err != nil {
			return err
		}
		d.SetPartial("name")
	}

	/*
	NOTE: plugin will manage disks of type "Data" only, and type cannot be changed once disk is created

	oldType, newType := d.GetChange("type")
	if oldType.(string) != newType.(string) {
		return fmt.Errorf("resourceDiskUpdate: Disk ID %s - changing type of existing disk not allowed", d.Id())
	}
	*/

	d.Partial(false)

	// we may reuse dataSourceDiskRead here as we maintain similarity 
	// between Compute resource and Compute data source schemas
	return dataSourceDiskRead(d, m) 
}

func resourceDiskDelete(d *schema.ResourceData, m interface{}) error {
	// NOTE: this function tries to detach and destroy target Disk "permanently", so 
	// there is no way to restore it. 
	// If, however, the disk is attached to a compute, the method will
	// fail (by failing the underpinning DECORt API call, which is issued with detach=false)
	log.Debugf("resourceDiskDelete: called for Disk ID / name %d / %s, Account ID %d",
		d.Get("disk_id").(int), d.Get("name").(string), d.Get("account_id").(int))

	diskFacts, err := utilityDiskCheckPresence(d, m)
	if diskFacts == "" {
		// the specified Disk does not exist - in this case according to Terraform best practice
		// we exit from Destroy method without error
		return nil
	}

	params := &url.Values{}
	params.Add("diskId", d.Id())
	// NOTE: we are not force-detaching disk from a compute (if attached) thus protecting
	// data that may be on that disk from destruction.
	// However, this may change in the future, as TF state management logic may want
	// to delete disk resource BEFORE it is detached from compute instance, and, while
	// perfectly OK from data preservation viewpoint, this is breaking expected TF workflow
	// in the eyes of an experienced TF user 
	params.Add("detach", "0") 
	params.Add("permanently", "1")

	controller := m.(*ControllerCfg)
	_, err = controller.decortAPICall("POST", DisksDeleteAPI, params)
	if err != nil {
		return err
	}

	return nil
}

func resourceDiskExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	log.Debugf("resourceDiskExists: called for Disk ID / name %d / %s, Account ID %d",
		d.Get("disk_id").(int), d.Get("name").(string), d.Get("account_id").(int))

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
			ValidateFunc: validation.IntAtLeast(1),
			Description: "Storage end-point provider serving this disk. Cannot be changed for existing disk.",
		},

		"pool": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description: "Pool where this disk is located. Cannot be changed for existing disk.",
		},

		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntAtLeast(1),
			Description: "Size of the disk in GB. Note, that existing disks can only be grown in size.",
		},

		/* We moved "type" attribute to computed attributes section, as plugin manages disks of only 
		   one type - "D", e.g. data disks.
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "D",
			StateFunc:   stateFuncToUpper,
			ValidateFunc: validation.StringInSlice([]string{"B", "D"}, false),
			Description: "Optional type of this disk. Defaults to D, i.e. data disk. Cannot be changed for existing disks.",
		},
		*/

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Disk resource managed by Terraform",
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

		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of this disk.",
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

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

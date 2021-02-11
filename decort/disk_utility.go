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
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"
)


func utilityDiskCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	// This function tries to locate Disk by one of the following algorithms depending on 
	// the parameters passed:
	//    - if disk ID is specified -> by disk ID
	//    - if disk name is specifeid -> by disk name and either account ID or account name
	//
	// NOTE: disk names are not unique, so the first occurence of this name in the account will
	// be returned. There is no such ambiguity when locating disk by its ID.
	// 
	// If succeeded, it returns non empty string that contains JSON formatted facts about the disk
	// as returned by disks/get API call.
	// Otherwise it returns empty string and meaningful error.
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for resource's Exists method.
	//

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}

	disk_id, arg_set := d.GetOk("disk_id")
	if arg_set {
		// go straight for the disk by its ID
		log.Debugf("utilityDiskCheckPresence: locating disk by its ID %d", disk_id.(int))
		url_values.Add("diskId", fmt.Sprintf("%d", disk_id.(int)))
		disk_facts, err := controller.decortAPICall("POST", DisksGetAPI, url_values)
		if err != nil {
			return "", err
		}
		return disk_facts, nil
	}

	disk_name, arg_set := d.GetOk("name")
	if !arg_set {
		// no disk ID and no disk name - we cannot locate disk in this case
		return "", fmt.Error("Cannot locate disk if name is empty and no disk ID specified.")
	}

	// Valid account ID is required to locate disks
	// obtain Account ID by account name - it should not be zero on success
	validated_account_id, err := utilityGetAccountIdBySchema(d, m)
	if err != nil {
		return err
	}

	url_values.Add("accountId", fmt.Sprintf("%d", validated_account_id))
	disk_facts, err := controller.decortAPICall("POST", DisksListAPI, url_values)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityDiskCheckPresence: ready to unmarshal string %q", disk_facts) 

	disks_list := []DiskRecord
	err = json.Unmarshal([]byte(disk_facts), &disks_list)
	if err != nil {
		return "", err
	}

	// log.Printf("%#v", vm_list)
	log.Debugf("utilityDiskCheckPresence: traversing decoded JSON of length %d", len(disks_list))
	for _, item := range disks_list {
		// need to match disk by name, return the first match
		if item.Name == disk_name.(string) && item.Status != "DESTROYED" {
			log.Printf("utilityDiskCheckPresence: index %d, matched disk name %q", index, item.Name)
			// we found the disk we need - not get detailed information via API call to disks/get
			/*
			// TODO: this may not be optimal as it initiates one extra call to the DECORT controller
			// in spite of the fact that we already have all required information about the disk in
			// item variable
			//
			get_url_values := &url.Values{}
			get_url_values.Add("diskId", fmt.Sprintf("%d", item.ID))
			disk_facts, err = controller.decortAPICall("POST", DisksGetAPI, get_url_values)
			if err != nil {
				return "", err
			}
			return disk_facts, nil
			*/
			reencoded_item, err := json.Marshal(item)
			if err != nil {
				return "", err
			}
			return reencoded_item.(string), nil 
		}
	}

	return "", nil // there should be no error if disk does not exist
}
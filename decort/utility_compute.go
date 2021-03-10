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

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func (ctrl *ControllerCfg) utilityComputeExtraDisksConfigure(d *schema.ResourceData, do_delta bool) error {
	// d is filled with data according to computeResource schema, so extra disks config is retrieved via "extra_disks" key
	// If do_delta is true, this function will identify changes between new and existing specs for extra disks and try to 
	// update compute configuration accordingly

	// Note that this function will not abort on API errors, but will continue to configure (attach / detach) other individual 
	// disks via atomic API calls. However, it will not retry failed manipulation on the same disk.

	old_set, new_set := d.GetChange("extra_disks")

	old_disks := old_set.([]interface{}) // NB: "extra_disks" is an array of ints
	new_disks := new_set.([]interface{})

	apiErrCount := 0
	var lastSavedError error

	if !do_delta {
		log.Debugf("utilityComputeExtraDisksConfigure: called for Compute ID %s with do_delta = false", d.Id())
		if len(new_disks) < 1 {
			return nil
		}

		for _, disk := range new_disks {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("diskId", fmt.Sprintf("%d", disk.(int)))
			_, err := ctrl.decortAPICall("POST", ComputeDiskAttachAPI, urlValues)
			if err != nil {
				// failed to attach extra disk - partial resource update
				apiErrCount++
				lastSavedError = err
			}
		}

		if apiErrCount > 0 {
			log.Errorf("utilityComputeExtraDisksConfigure: there were %d error(s) when attaching disks to Compute ID %s. Last error was: %s", 
			           apiErrCount, d.Id(), lastSavedError)
			return lastSavedError
		}
		
		return nil
	}

	attach_list := make([]int, 0, MaxExtraDisksPerCompute)
	detach_list := make([]int, 0, MaxExtraDisksPerCompute)
	attIdx := 0
	detIdx := 0
	match := false
	
	for _, oDisk := range old_disks {
		match = false
		for _, nDisk := range new_disks {
			if oDisk.(int) == nDisk.(int) {
				match = true
				break
			}
		}
		if !match {
			detach_list[detIdx] = oDisk.(int)
			detIdx++
		}
	}
	log.Debugf("utilityComputeExtraDisksConfigure: detach list has %d items for Compute ID %s", len(detach_list), d.Id())

	for _, nDisk := range new_disks {
		match = false
		for _, oDisk := range old_disks {
			if nDisk.(int) == oDisk.(int) {
				match = true
				break
			}
		}
		if !match {
			attach_list[attIdx] = nDisk.(int)
			attIdx++
		}
	}
	log.Debugf("utilityComputeExtraDisksConfigure: attach list has %d items for Compute ID %s", len(attach_list), d.Id())

	for _, diskId := range detach_list {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("diskId", fmt.Sprintf("%d", diskId))
		_, err := ctrl.decortAPICall("POST", ComputeDiskDetachAPI, urlValues)
		if err != nil {
			// failed to detach disk - there will be partial resource update
			log.Debugf("utilityComputeExtraDisksConfigure: failed to detach disk ID %d from Compute ID %s: %s", diskId, d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	for _, diskId := range attach_list {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("diskId", fmt.Sprintf("%d", diskId))
		_, err := ctrl.decortAPICall("POST", ComputeDiskAttachAPI, urlValues)
		if err != nil {
			// failed to attach disk - there will be partial resource update
			log.Debugf("utilityComputeExtraDisksConfigure: failed to attach disk ID %d to Compute ID %s: %s", diskId, d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	if apiErrCount > 0 {
		log.Errorf("utilityComputeExtraDisksConfigure: there were %d error(s) when managing disks on Compute ID %s. Last error was: %s", 
				   apiErrCount, d.Id(), lastSavedError)
		return lastSavedError
	}

	return nil
}

// TODO: implement do_delta logic
func (ctrl *ControllerCfg) utilityComputeNetworksConfigure(d *schema.ResourceData, do_delta bool) error {
	// "d" is filled with data according to computeResource schema, so extra networks config is retrieved via "networks" key
	// If do_delta is true, this function will identify changes between new and existing specs for network and try to 
	// update compute configuration accordingly
	argVal, argSet := d.GetOk("networks") 
	if !argSet || len(argVal.([]interface{})) < 1 {
		return nil
	}

	net_list := argVal.([]interface{}) // networks" is ar array of maps; for keys see func networkSubresourceSchemaMake() definition 

	for _, net := range net_list {
		urlValues := &url.Values{}
		net_data := net.(map[string]interface{}) 
		urlValues.Add("computeId", d.Id())
		urlValues.Add("netType", net_data["net_type"].(string))
		urlValues.Add("netId", fmt.Sprintf("%d", net_data["net_id"].(int)))
		ipaddr, ipSet := net_data["ip_address"] // "ip_address" key is optional
		if ipSet {
			urlValues.Add("ipAddr", ipaddr.(string))
		}
		_, err := ctrl.decortAPICall("POST", ComputeNetAttachAPI, urlValues)
		if err != nil {
			// failed to attach network - partial resource update
			return err
		}
	}
	return nil
}

func utilityComputeCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	// This function tries to locate Compute by one of the following approaches:
	// - if compute_id is specified - locate by compute ID
	// - if compute_name is specified - locate by a combination of compute name and resource
	//   group ID
	//
	// If succeeded, it returns non-empty string that contains JSON formatted facts about the
	// Compute as returned by compute/get API call.
	// Otherwise it returns empty string and meaningful error.
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for resource's Exists method.
	//

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	computeId, argSet := d.GetOk("compute_id") // NB: compute_id is NOT present in computeResource schema!
	if argSet {
		// compute ID is specified, try to get compute instance straight by this ID
		log.Debugf("utilityComputeCheckPresence: locating compute by its ID %d", computeId.(int))
		urlValues.Add("computeId", fmt.Sprintf("%d", computeId.(int)))
		computeFacts, err := controller.decortAPICall("POST", ComputeGetAPI, urlValues)
		if err != nil {
			return "", err
		}
		return computeFacts, nil
	}

	computeName, argSet := d.GetOk("name")
	if !argSet {
		return "", fmt.Errorf("Cannot locate compute instance if name is empty and no compute ID specified")
	}

	rgId, argSet := d.GetOk("rg_id")
	if !argSet {
		return "", fmt.Errorf("Cannot locate compute by name %s if no resource group ID is set", computeName.(string))
	}
	
	urlValues.Add("rgId", fmt.Sprintf("%d", rgId))
	apiResp, err := controller.decortAPICall("POST", RgListComputesAPI, urlValues)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityComputeCheckPresence: ready to unmarshal string %q", apiResp)

	computeList := RgListComputesResp{}
	err = json.Unmarshal([]byte(apiResp), &computeList)
	if err != nil {
		return "", err
	}

	// log.Printf("%#v", computeList)
	log.Debugf("utilityComputeCheckPresence: traversing decoded JSON of length %d", len(computeList))
	for index, item := range computeList {
		// need to match Compute by name, skip Computes with the same name in DESTROYED satus
		if item.Name == computeName.(string) && item.Status != "DESTROYED" {
			log.Debugf("utilityComputeCheckPresence: index %d, matched name %q", index, item.Name)
			// we found the Compute we need - now get detailed information via compute/get API
			cgetValues := &url.Values{}
			cgetValues.Add("computeId", fmt.Sprintf("%d", item.ID))
			apiResp, err = controller.decortAPICall("POST", ComputeGetAPI, cgetValues)
			if err != nil {
				return "", err
			}
			return apiResp, nil
		}
	}

	return "", nil // there should be no error if Compute does not exist
	// return "", fmt.Errorf("Cannot find Compute name %q in resource group ID %d", name, rgid)
}

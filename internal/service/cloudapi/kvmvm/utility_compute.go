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

package kvmvm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityComputeExtraDisksConfigure(ctx context.Context, d *schema.ResourceData, m interface{}, do_delta bool) error {
	// d is filled with data according to computeResource schema, so extra disks config is retrieved via "extra_disks" key
	// If do_delta is true, this function will identify changes between new and existing specs for extra disks and try to
	// update compute configuration accordingly
	// Otherwise it will apply whatever is found in the new set of "extra_disks" right away.
	// Primary use of do_delta=false is when calling this function from compute Create handler.

	// Note that this function will not abort on API errors, but will continue to configure (attach / detach) other individual
	// disks via atomic API calls. However, it will not retry failed manipulation on the same disk.
	c := m.(*controller.ControllerCfg)

	log.Debugf("utilityComputeExtraDisksConfigure: called for Compute ID %s with do_delta = %t", d.Id(), do_delta)

	// NB: as of rc-1.25 "extra_disks" are TypeSet with the elem of TypeInt
	old_set, new_set := d.GetChange("extra_disks")

	apiErrCount := 0
	var lastSavedError error

	if !do_delta {
		if new_set.(*schema.Set).Len() < 1 {
			return nil
		}

		for _, disk := range new_set.(*schema.Set).List() {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("diskId", fmt.Sprintf("%d", disk.(int)))
			_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAttachAPI, urlValues)
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

	detach_set := old_set.(*schema.Set).Difference(new_set.(*schema.Set))
	log.Debugf("utilityComputeExtraDisksConfigure: detach set has %d items for Compute ID %s", detach_set.Len(), d.Id())
	for _, diskId := range detach_set.List() {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("diskId", fmt.Sprintf("%d", diskId.(int)))
		_, err := c.DecortAPICall(ctx, "POST", ComputeDiskDetachAPI, urlValues)
		if err != nil {
			// failed to detach disk - there will be partial resource update
			log.Errorf("utilityComputeExtraDisksConfigure: failed to detach disk ID %d from Compute ID %s: %s", diskId.(int), d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	attach_set := new_set.(*schema.Set).Difference(old_set.(*schema.Set))
	log.Debugf("utilityComputeExtraDisksConfigure: attach set has %d items for Compute ID %s", attach_set.Len(), d.Id())
	for _, diskId := range attach_set.List() {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("diskId", fmt.Sprintf("%d", diskId.(int)))
		_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAttachAPI, urlValues)
		if err != nil {
			// failed to attach disk - there will be partial resource update
			log.Errorf("utilityComputeExtraDisksConfigure: failed to attach disk ID %d to Compute ID %s: %s", diskId.(int), d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	if apiErrCount > 0 {
		log.Errorf("utilityComputeExtraDisksConfigure: there were %d error(s) when managing disks of Compute ID %s. Last error was: %s",
			apiErrCount, d.Id(), lastSavedError)
		return lastSavedError
	}

	return nil
}

func utilityComputeNetworksConfigure(ctx context.Context, d *schema.ResourceData, m interface{}, do_delta bool) error {
	// "d" is filled with data according to computeResource schema, so extra networks config is retrieved via "network" key
	// If do_delta is true, this function will identify changes between new and existing specs for network and try to
	// update compute configuration accordingly
	// Otherwise it will apply whatever is found in the new set of "network" right away.
	// Primary use of do_delta=false is when calling this function from compute Create handler.

	c := m.(*controller.ControllerCfg)

	old_set, new_set := d.GetChange("network")

	apiErrCount := 0
	var lastSavedError error

	if !do_delta {
		if new_set.(*schema.Set).Len() < 1 {
			return nil
		}

		for _, runner := range new_set.(*schema.Set).List() {
			urlValues := &url.Values{}
			net_data := runner.(map[string]interface{})
			urlValues.Add("computeId", d.Id())
			urlValues.Add("netType", net_data["net_type"].(string))
			urlValues.Add("netId", fmt.Sprintf("%d", net_data["net_id"].(int)))
			ipaddr, ipSet := net_data["ip_address"] // "ip_address" key is optional
			if ipSet {
				urlValues.Add("ipAddr", ipaddr.(string))
			}
			_, err := c.DecortAPICall(ctx, "POST", ComputeNetAttachAPI, urlValues)
			if err != nil {
				// failed to attach network - partial resource update
				apiErrCount++
				lastSavedError = err
			}
		}

		if apiErrCount > 0 {
			log.Errorf("utilityComputeNetworksConfigure: there were %d error(s) when managing networks of Compute ID %s. Last error was: %s",
				apiErrCount, d.Id(), lastSavedError)
			return lastSavedError
		}
		return nil
	}

	detach_set := old_set.(*schema.Set).Difference(new_set.(*schema.Set))
	log.Debugf("utilityComputeNetworksConfigure: detach set has %d items for Compute ID %s", detach_set.Len(), d.Id())
	for _, runner := range detach_set.List() {
		urlValues := &url.Values{}
		net_data := runner.(map[string]interface{})
		urlValues.Add("computeId", d.Id())
		urlValues.Add("ipAddr", net_data["ip_address"].(string))
		urlValues.Add("mac", net_data["mac"].(string))
		_, err := c.DecortAPICall(ctx, "POST", ComputeNetDetachAPI, urlValues)
		if err != nil {
			// failed to detach this network - there will be partial resource update
			log.Errorf("utilityComputeNetworksConfigure: failed to detach net ID %d of type %s from Compute ID %s: %s",
				net_data["net_id"].(int), net_data["net_type"].(string), d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	attach_set := new_set.(*schema.Set).Difference(old_set.(*schema.Set))
	log.Debugf("utilityComputeNetworksConfigure: attach set has %d items for Compute ID %s", attach_set.Len(), d.Id())
	for _, runner := range attach_set.List() {
		urlValues := &url.Values{}
		net_data := runner.(map[string]interface{})
		urlValues.Add("computeId", d.Id())
		urlValues.Add("netId", fmt.Sprintf("%d", net_data["net_id"].(int)))
		urlValues.Add("netType", net_data["net_type"].(string))
		if net_data["ip_address"].(string) != "" {
			urlValues.Add("ipAddr", net_data["ip_address"].(string))
		}
		_, err := c.DecortAPICall(ctx, "POST", ComputeNetAttachAPI, urlValues)
		if err != nil {
			// failed to attach this network - there will be partial resource update
			log.Errorf("utilityComputeNetworksConfigure: failed to attach net ID %d of type %s to Compute ID %s: %s",
				net_data["net_id"].(int), net_data["net_type"].(string), d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	if apiErrCount > 0 {
		log.Errorf("utilityComputeNetworksConfigure: there were %d error(s) when managing networks of Compute ID %s. Last error was: %s",
			apiErrCount, d.Id(), lastSavedError)
		return lastSavedError
	}

	return nil
}

func utilityComputeCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
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

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	// make it possible to use "read" & "check presence" functions with compute ID set so
	// that Import of Compute resource is possible
	idSet := false
	theId, err := strconv.Atoi(d.Id())
	if err != nil || theId <= 0 {
		computeId, argSet := d.GetOk("compute_id") // NB: compute_id is NOT present in computeResource schema!
		if argSet {
			theId = computeId.(int)
			idSet = true
		}
	} else {
		idSet = true
	}

	if idSet {
		// compute ID is specified, try to get compute instance straight by this ID
		log.Debugf("utilityComputeCheckPresence: locating compute by its ID %d", theId)
		urlValues.Add("computeId", fmt.Sprintf("%d", theId))
		computeFacts, err := c.DecortAPICall(ctx, "POST", ComputeGetAPI, urlValues)
		if err != nil {
			return "", err
		}
		return computeFacts, nil
	}

	// ID was not set in the schema upon entering this function - work through Compute name
	// and RG ID
	computeName, argSet := d.GetOk("name")
	if !argSet {
		return "", fmt.Errorf("cannot locate compute instance if name is empty and no compute ID specified")
	}

	rgId, argSet := d.GetOk("rg_id")
	if !argSet {
		return "", fmt.Errorf("cannot locate compute by name %s if no resource group ID is set", computeName.(string))
	}

	urlValues.Add("rgId", fmt.Sprintf("%d", rgId))
	apiResp, err := c.DecortAPICall(ctx, "POST", RgListComputesAPI, urlValues)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityComputeCheckPresence: ready to unmarshal string %s", apiResp)

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
			log.Debugf("utilityComputeCheckPresence: index %d, matched name %s", index, item.Name)
			// we found the Compute we need - now get detailed information via compute/get API
			cgetValues := &url.Values{}
			cgetValues.Add("computeId", fmt.Sprintf("%d", item.ID))
			apiResp, err = c.DecortAPICall(ctx, "POST", ComputeGetAPI, cgetValues)
			if err != nil {
				return "", err
			}
			return apiResp, nil
		}
	}

	return "", nil // there should be no error if Compute does not exist
}

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

/*

func (ctrl *ControllerCfg) utilityVmDisksProvision(mcfg *MachineConfig) error {
	for index, disk := range mcfg.DataDisks {
		urlValues := &url.Values{}
		// urlValues.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		urlValues.Add("accountId", fmt.Sprintf("%d", mcfg.TenantID))
		urlValues.Add("gid", fmt.Sprintf("%d", mcfg.GridID))
		urlValues.Add("name", fmt.Sprintf("%s", disk.Label))
		urlValues.Add("description", fmt.Sprintf("Data disk for VM ID %d / VM Name: %s", mcfg.ID, mcfg.Name))
		urlValues.Add("size", fmt.Sprintf("%d", disk.Size))
		urlValues.Add("type", "D")
		// urlValues.Add("iops", )

		disk_id_resp, err := ctrl.decortAPICall("POST", DiskCreateAPI, urlValues)
		if err != nil {
			// failed to create disk - partial resource update
			return err
		}
		// disk created - API call returns disk ID as a string - use it to update
		// disk ID in the corresponding MachineConfig.DiskConfig record

		mcfg.DataDisks[index].ID, err = strconv.Atoi(disk_id_resp)
		if err != nil {
			// failed to convert disk ID into proper integer value - partial resource update
			return err
		}

		// now that we have disk created and stored its ID in the mcfg.DataDisks[index].ID
		// we can attempt attaching the disk to the VM
		urlValues = &url.Values{}
		// urlValues.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		urlValues.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		urlValues.Add("diskId", disk_id_resp)
		_, err = ctrl.decortAPICall("POST", DiskAttachAPI, urlValues)
		if err != nil {
			// failed to attach disk - partial resource update
			return err
		}
	}
	return nil
}


func (ctrl *ControllerCfg) utilityVmPortforwardsProvision(mcfg *MachineConfig) error {
	for _, rule := range mcfg.PortForwards {
		urlValues := &url.Values{}
		urlValues.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		urlValues.Add("cloudspaceId", fmt.Sprintf("%d", mcfg.ResGroupID))
		urlValues.Add("publicIp", mcfg.ExtIP) // this may be obsoleted by Resource group implementation
		urlValues.Add("publicPort", fmt.Sprintf("%d", rule.ExtPort))
		urlValues.Add("localPort", fmt.Sprintf("%d", rule.IntPort))
		urlValues.Add("protocol", rule.Proto)
		_, err := ctrl.decortAPICall("POST", PortforwardingCreateAPI, urlValues)
		if err != nil {
			// failed to create port forward rule - partial resource update
			return err
		}
	}
	return nil
}

func (ctrl *ControllerCfg) utilityVmNetworksProvision(mcfg *MachineConfig) error {
	for _, net := range mcfg.Networks {
		urlValues := &url.Values{}
		urlValues.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		urlValues.Add("externalNetworkId", fmt.Sprintf("%d", net.NetworkID))
		_, err := ctrl.decortAPICall("POST", AttachExternalNetworkAPI, urlValues)
		if err != nil {
			// failed to attach network - partial resource update
			return err
		}
	}
	return nil
}

*/

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

	computeId, argSet := d.GetOk("compute_id")
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

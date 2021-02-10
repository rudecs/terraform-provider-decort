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

/*

func (ctrl *ControllerCfg) utilityVmDisksProvision(mcfg *MachineConfig) error {
	for index, disk := range mcfg.DataDisks {
		url_values := &url.Values{}
		// url_values.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		url_values.Add("accountId", fmt.Sprintf("%d", mcfg.TenantID))
		url_values.Add("gid", fmt.Sprintf("%d", mcfg.GridID))
		url_values.Add("name", fmt.Sprintf("%s", disk.Label))
		url_values.Add("description", fmt.Sprintf("Data disk for VM ID %d / VM Name: %s", mcfg.ID, mcfg.Name))
		url_values.Add("size", fmt.Sprintf("%d", disk.Size))
		url_values.Add("type", "D")
		// url_values.Add("iops", )

		disk_id_resp, err := ctrl.decortAPICall("POST", DiskCreateAPI, url_values)
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
		url_values = &url.Values{}
		// url_values.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		url_values.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		url_values.Add("diskId", disk_id_resp)
		_, err = ctrl.decortAPICall("POST", DiskAttachAPI, url_values)
		if err != nil {
			// failed to attach disk - partial resource update
			return err
		}
	}
	return nil
}


func (ctrl *ControllerCfg) utilityVmPortforwardsProvision(mcfg *MachineConfig) error {
	for _, rule := range mcfg.PortForwards {
		url_values := &url.Values{}
		url_values.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		url_values.Add("cloudspaceId", fmt.Sprintf("%d", mcfg.ResGroupID))
		url_values.Add("publicIp", mcfg.ExtIP) // this may be obsoleted by Resource group implementation 
		url_values.Add("publicPort", fmt.Sprintf("%d", rule.ExtPort))
		url_values.Add("localPort", fmt.Sprintf("%d", rule.IntPort))
		url_values.Add("protocol", rule.Proto)
		_, err := ctrl.decortAPICall("POST", PortforwardingCreateAPI, url_values)
		if err != nil {
			// failed to create port forward rule - partial resource update
			return err
		}
	}
	return nil
}

func (ctrl *ControllerCfg) utilityVmNetworksProvision(mcfg *MachineConfig) error {
	for _, net := range mcfg.Networks {
		url_values := &url.Values{}
		url_values.Add("machineId", fmt.Sprintf("%d", mcfg.ID))
		url_values.Add("externalNetworkId", fmt.Sprintf("%d", net.NetworkID))
		_, err := ctrl.decortAPICall("POST", AttachExternalNetworkAPI, url_values)
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

	compute_id, arg_set := d.GetOk("compute_id")
	if arg_set {
		// compute ID is specified, try to get compute instance straight by this ID
		log.Debugf("utilityComputeCheckPresence: locating compute by its ID %d", compute_id.(int))
		url_values.Add("computeId", fmt.Sprintf("%d", compute_id.(int)))
		compute_facts, err := controller.decortAPICall("POST", ComputeGetAPI, url_values)
		if err != nil {
			return "", err
		}
		return compute_facts, nil
	}

	compute_name, arg_set := d.GetOk("name")
	if !arg_set {
		return "", fmt.Error("Cannot locate compute instance if name is empty and no compute ID specified.")
	}

	rg_id, arg_set := d.GetOk("rg_id")
	if !arg_set {
		return "", fmt.Error("Cannot locate compute by name %s if no resource group ID is set", compute_name.(string))
	}


	controller := m.(*ControllerCfg)
	list_url_values := &url.Values{}
	list_url_values.Add("rgId", fmt.Sprintf("%d",rg_id))
	api_resp, err := controller.decortAPICall("POST", RgListComputesAPI, list_url_values)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityComputeCheckPresence: ready to unmarshal string %q", api_resp) 

	comp_list := RgListComputesResp{}
	err = json.Unmarshal([]byte(api_resp), &comp_list)
	if err != nil {
		return "", err
	}

	// log.Printf("%#v", comp_list)
	log.Debugf("utilityComputeCheckPresence: traversing decoded JSON of length %d", len(comp_list))
	for _, item := range comp_list {
		// need to match Compute by name, skip Computes with the same name in DESTROYED satus
		if item.Name == name && item.Status != "DESTROYED" {
			log.Debugf("utilityComputeCheckPresence: index %d, matched name %q", index, item.Name)
			// we found the Compute we need - now get detailed information via compute/get API
			get_url_values := &url.Values{}
			get_url_values.Add("computeId", fmt.Sprintf("%d", item.ID))
			api_resp, err = controller.decortAPICall("POST", ComputeGetAPI, get_url_values)
			if err != nil {
				return "", err
			}
			return api_resp, nil
		}
	}

	return "", nil // there should be no error if Compute does not exist
	// return "", fmt.Errorf("Cannot find Compute name %q in resource group ID %d", name, rgid)
}
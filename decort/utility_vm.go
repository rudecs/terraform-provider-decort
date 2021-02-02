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

func utilityVmCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	// This function tries to locate VM by its name and resource group ID
	// if succeeded, it returns non empty string that contains JSON formatted facts about the VM
	// as returned by machines/get API call.
	// Otherwise it returns empty string and meaningful error.
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for resource's Exists method.
	//
	name := d.Get("name").(string)
	rgid := d.Get("rgid").(int)

	controller := m.(*ControllerCfg)
	list_url_values := &url.Values{}
	list_url_values.Add("cloudspaceId", fmt.Sprintf("%d",rgid))
	body_string, err := controller.decortAPICall("POST", MachinesListAPI, list_url_values)
	if err != nil {
		return "", err
	}

	// log.Printf("%s", body_string)
	// log.Printf("dataSourceVmRead: ready to decode mashines/list response body")
	vm_list := MachinesListResp{}
	err = json.Unmarshal([]byte(body_string), &vm_list)
	if err != nil {
		return "", err
	}

	// log.Printf("%#v", vm_list)
	// log.Printf("dataSourceVmRead: traversing decoded JSON of length %d", len(vm_list))
	for _, item := range vm_list {
		// need to match VM by name, skip VMs with the same name in DESTROYED satus
		if item.Name == name && item.Status != "DESTROYED" {
			// log.Printf("dataSourceVmRead: index %d, matched name %q", index, item.Name)
			// we found the VM we need - not get detailed information via API call to cloudapi/machines/get
			get_url_values := &url.Values{}
			get_url_values.Add("machineId", fmt.Sprintf("%d", item.ID))
			body_string, err = controller.decortAPICall("POST", MachinesGetAPI, get_url_values)
			if err != nil {
				return "", err
			}
			return body_string, nil
		}
	}

	return "", nil // there should be no error if VM does not exist
	// return "", fmt.Errorf("Cannot find VM name %q in resource group ID %d", name, rgid)
}
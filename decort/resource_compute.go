/*
Copyright (c) 2019-2020 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceComputeCreate(d *schema.ResourceData, m interface{}) error {
	/*
	machine := &MachineConfig{
		ResGroupID:  d.Get("rgid").(int),
		Name:        d.Get("name").(string),
		Cpu:         d.Get("cpu").(int),
		Ram:         d.Get("ram").(int),
		ImageID:     d.Get("image_id").(int),
		Description: d.Get("description").(string),
	}
	// BootDisk
	// DataDisks
	// Networks
	// PortForwards
	// SshKeyData string
	log.Printf("resourceComputeCreate: called for VM name %q, ResGroupID %d", machine.Name, machine.ResGroupID)

	var subres_list []interface{}
	var subres_data map[string]interface{}
	var arg_value interface{}
	var arg_set bool
	// boot disk list is a required argument and has only one element,
	// which is of type diskSubresourceSchema
	subres_list = d.Get("boot_disk").([]interface{})
	subres_data = subres_list[0].(map[string]interface{})
	machine.BootDisk.Label = subres_data["label"].(string)
	machine.BootDisk.Size = subres_data["size"].(int)
	machine.BootDisk.Pool = subres_data["pool"].(string)
	machine.BootDisk.Provider = subres_data["provider"].(string)

	arg_value, arg_set = d.GetOk("data_disks")
	if arg_set {
		log.Printf("resourceComputeCreate: calling makeDisksConfig")
		machine.DataDisks, _ = makeDisksConfig(arg_value.([]interface{}))
	}

	arg_value, arg_set = d.GetOk("networks")
	if arg_set {
		log.Printf("resourceComputeCreate: calling makeNetworksConfig")
		machine.Networks, _ = makeNetworksConfig(arg_value.([]interface{}))
	}

	arg_value, arg_set = d.GetOk("port_forwards")
	if arg_set {
		log.Printf("resourceComputeCreate: calling makePortforwardsConfig")
		machine.PortForwards, _ = makePortforwardsConfig(arg_value.([]interface{}))
	}

	arg_value, arg_set = d.GetOk("ssh_keys")
	if arg_set {
		log.Printf("resourceComputeCreate: calling makeSshKeysConfig")
		machine.SshKeys, _ = makeSshKeysConfig(arg_value.([]interface{}))
	}

	// create basic VM (i.e. without port forwards and ext network connections - those will be done
	// by separate API calls)
	d.Partial(true)
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("cloudspaceId", fmt.Sprintf("%d", machine.ResGroupID))
	urlValues.Add("name", machine.Name)
	urlValues.Add("description", machine.Description)
	urlValues.Add("vcpus", fmt.Sprintf("%d", machine.Cpu))
	urlValues.Add("memory", fmt.Sprintf("%d", machine.Ram))
	urlValues.Add("imageId", fmt.Sprintf("%d", machine.ImageID))
	urlValues.Add("disksize", fmt.Sprintf("%d", machine.BootDisk.Size))
	if len(machine.SshKeys) > 0 {
		urlValues.Add("userdata", makeSshKeysArgString(machine.SshKeys))
	}
	api_resp, err := controller.decortAPICall("POST", MachineCreateAPI, urlValues)
	if err != nil {
		return err
	}
	d.SetId(api_resp) // machines/create API plainly returns ID of the new VM on success
	machine.ID, _ = strconv.Atoi(api_resp)
	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("cpu")
	d.SetPartial("ram")
	d.SetPartial("image_id")
	d.SetPartial("boot_disk")
	if len(machine.SshKeys) > 0 {
		d.SetPartial("ssh_keys")
	}

	log.Printf("resourceComputeCreate: new VM ID %d, name %q created", machine.ID, machine.Name)

	if len(machine.DataDisks) > 0 || len(machine.PortForwards) > 0 {
		// for data disk or port foreards provisioning we have to know Tenant ID
		// and Grid ID so we call utilityResgroupConfigGet method to populate these
		// fields in the machine structure that will be passed to provisionVmDisks or
		// provisionVmPortforwards
		log.Printf("resourceComputeCreate: calling utilityResgroupConfigGet")
		resgroup, err := controller.utilityResgroupConfigGet(machine.ResGroupID)
		if err == nil {
			machine.TenantID = resgroup.TenantID
			machine.GridID = resgroup.GridID
			machine.ExtIP = resgroup.ExtIP
			log.Printf("resourceComputeCreate: tenant ID %d, GridID %d, ExtIP %q",
				machine.TenantID, machine.GridID, machine.ExtIP)
		}
	}

	//
	// Configure data disks
	disks_ok := true
	if len(machine.DataDisks) > 0 {
		log.Printf("resourceComputeCreate: calling utilityVmDisksProvision for disk count %d", len(machine.DataDisks))
		if machine.TenantID == 0 {
			// if TenantID is still 0 it means that we failed to get Resgroup Facts by
			// a previous call to utilityResgroupGetFacts,
			// hence we do not have technical ability to provision data disks
			disks_ok = false
		} else {
			// provisionVmDisks accomplishes two steps for each data disk specification
			// 1) creates the disks
			// 2) attaches them to the VM
			err = controller.utilityVmDisksProvision(machine)
			if err != nil {
				disks_ok = false
			}
		}
	}

	if disks_ok {
		d.SetPartial("data_disks")
	}

	//
	// Configure port forward rules
	pfws_ok := true
	if len(machine.PortForwards) > 0 {
		log.Printf("resourceComputeCreate: calling utilityVmPortforwardsProvision for pfw rules count %d", len(machine.PortForwards))
		if machine.ExtIP == "" {
			// if ExtIP is still empty it means that we failed to get Resgroup Facts by
			// a previous call to utilityResgroupGetFacts,
			// hence we do not have technical ability to provision port forwards
			pfws_ok = false
		} else {
			err := controller.utilityVmPortforwardsProvision(machine)
			if err != nil {
				pfws_ok = false
			}
		}
	}
	if pfws_ok {
		//  there were no errors reported when configuring port forwards
		d.SetPartial("port_forwards")
	}

	//
	// Configure external networks
	// NOTE: currently only one external network can be attached to each VM, so in the current
	// implementation we ignore all but the 1st network definition
	nets_ok := true
	if len(machine.Networks) > 0 {
		log.Printf("resourceComputeCreate: calling utilityVmNetworksProvision for networks count %d", len(machine.Networks))
		err := controller.utilityVmNetworksProvision(machine)
		if err != nil {
			nets_ok = false
		}
	}
	if nets_ok {
		// there were no errors reported when configuring networks
		d.SetPartial("networks")
	}

	if disks_ok && nets_ok && pfws_ok {
		// if there were no errors in setting any of the subresources, we may leave Partial mode
		d.Partial(false)
	}
	*/

	// resourceComputeRead will also update resource ID on success, so that Terraform will know
	// that resource exists
	return resourceComputeRead(d, m)
}

func resourceComputeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("resourceComputeRead: called for VM name %q, ResGroupID %d",
		d.Get("name").(string), d.Get("rgid").(int))

	compFacts, err := utilityComputeCheckPresence(d, m)
	if compFacts == "" {
		if err != nil {
			return err
		}
		// VM was not found
		return nil
	}

	if err = flattenCompute(d, compFacts); err != nil {
		return err
	}
	log.Printf("resourceComputeRead: after flattenCompute: VM ID %s, VM name %q, ResGroupID %d",
		d.Id(), d.Get("name").(string), d.Get("rgid").(int))

	// Not all parameters, that we may need, are returned by machines/get API
	// Continue with further reading of VM subresource parameters:
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	/*
		// Obtain information on external networks
		urlValues.Add("machineId", d.Id())
		body_string, err := controller.decortAPICall("POST", VmExtNetworksListAPI, urlValues)
		if err != nil {
			return err
		}

		net_list := ExtNetworksResp{}
		err = json.Unmarshal([]byte(body_string), &net_list)
		if err != nil {
			return err
		}

		if len(net_list) > 0 {
			if err = d.Set("networks", flattenNetworks(net_list)); err != nil {
				return err
			}
		}
	*/

	/*
		// Ext networks flattening is now done inside flattenCompute because it is currently based
		// on data read into NICs component by machine/get API call

		if err = d.Set("networks", flattenNetworks()); err != nil {
			return err
		}
	*/

	//
	// Obtain information on port forwards
	/*
	urlValues.Add("cloudspaceId", fmt.Sprintf("%d", d.Get("rgid")))
	urlValues.Add("machineId", d.Id())
	pfw_list := PortforwardsResp{}
	body_string, err := controller.decortAPICall("POST", PortforwardsListAPI, urlValues)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body_string), &pfw_list)
	if err != nil {
		return err
	}

	if len(pfw_list) > 0 {
		if err = d.Set("port_forwards", flattenPortforwards(pfw_list)); err != nil {
			return err
		}
	}
	*/

	return nil
}

func resourceComputeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("resourceComputeUpdate: called for VM name %q, ResGroupID %d",
		d.Get("name").(string), d.Get("rgid").(int))

	return resourceComputeRead(d, m)
}

func resourceComputeDelete(d *schema.ResourceData, m interface{}) error {
	// NOTE: this method destroys target Compute instance with flag "permanently", so 
	// there is no way to restore destroyed Compute
	log.Printf("resourceComputeDelete: called for VM name %q, ResGroupID %d",
		d.Get("name").(string), d.Get("rgid").(int))

	compFacts, err := utilityComputeCheckPresence(d, m)
	if compFacts == "" {
		// the target Compute does not exist - in this case according to Terraform best practice
		// we exit from Destroy method without error
		return nil
	}

	params := &url.Values{}
	params.Add("computeId", d.Id())
	params.Add("permanently", "true")

	controller := m.(*ControllerCfg)
	compFacts, err = controller.decortAPICall("POST", ComputeDeleteAPI, params)
	if err != nil {
		return err
	}

	return nil
}

func resourceComputeExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	log.Printf("resourceComputeExist: called for VM name %q, ResGroupID %d",
		d.Get("name").(string), d.Get("rgid").(int))

	compFacts, err := utilityComputeCheckPresence(d, m)
	if compFacts == "" {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func resourceCompute() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceComputeCreate,
		Read:   resourceComputeRead,
		Update: resourceComputeUpdate,
		Delete: resourceComputeDelete,
		Exists: resourceComputeExists,

		Timeouts: &schema.ResourceTimeout{
			Create:  &Timeout180s,
			Read:    &Timeout30s,
			Update:  &Timeout180s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this compute. This parameter is case sensitive and must be unique in the resource group.",
			},

			"rg_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the resource group where this compute should be deployed.",
			},

			"arch": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Hardware architecture of this compute instance.",
			},

			"cpu": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 64),
				Description:  "Number of CPUs to allocate to this compute instance.",
			},

			"ram": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(512),
				Description:  "Amount of RAM in MB to allocate to this compute instance.",
			},

			"image_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the OS image to base this compute instance on.",
			},

			"boot_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Size of the boot disk on this compute instance.",
			},

			"extra_disks": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: MaxExtraDisksPerCompute,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Optional list of IDs of the extra disks to attach to this compute.",
			},

			"ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: MaxSshKeysPerCompute,
				Elem: &schema.Resource{
					Schema: sshSubresourceSchemaMake(),
				},
				Description: "SSH keys to authorize on this compute instance.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of this compute instance.",
			},

			// The rest are Compute properties, which are "computed" once it is created
			"rg_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the resource group where this compute instance is located.",
			},

			"account_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the account this compute instance belongs to.",
			},

			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the account this compute instance belongs to.",
			},

			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceDiskSchemaMake(), // ID, type,  name, size, account ID, SEP ID, SEP type, pool, status, tech status, compute ID, image ID
				},
				Description: "Detailed specification for all disks attached to this compute instance (including bood disk).",
			},

			"interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: interfaceSubresourceSchemaMake(),
				},
				Description: "Specification for the virtual NICs configured on this compute instance.",
			},

			"guest_logins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: loginsSubresourceSchemaMake(),
				},
				Description: "Specification for guest logins on this compute instance.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current model status of this compute instance.",
			},

			"tech_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current technical status of this compute instance.",
			},
		},
	}
}

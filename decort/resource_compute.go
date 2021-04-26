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

func resourceComputeCreate(d *schema.ResourceData, m interface{}) error {
	// we assume all mandatory parameters it takes to create a comptue instance are properly 
	// specified - we rely on schema "Required" attributes to let Terraform validate them for us
	
	log.Debugf("resourceComputeCreate: called for Compute name %q, RG ID %d", d.Get("name").(string), d.Get("rg_id").(int))

	// create basic Compute (i.e. without extra disks and network connections - those will be attached
	// by subsequent individual API calls).
	// creating Compute is a multi-step workflow, which may fail at some step, so we use "partial" feature of Terraform
	d.Partial(true) 
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("rgId", fmt.Sprintf("%d", d.Get("rg_id").(int)))
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("cpu", fmt.Sprintf("%d", d.Get("cpu").(int)))
	urlValues.Add("ram", fmt.Sprintf("%d", d.Get("ram").(int)))
	urlValues.Add("imageId", fmt.Sprintf("%d", d.Get("image_id").(int)))
	urlValues.Add("bootDisk", fmt.Sprintf("%d", d.Get("boot_disk_size").(int)))
	urlValues.Add("netType", "NONE") // at the 1st step create isolated compute
	// urlValues.Add("start", "false") // at the 1st step create compute in a stopped state

	argVal, argSet := d.GetOk("description") 
	if argSet {
		urlValues.Add("desc", argVal.(string))
	}

	sshKeysVal, sshKeysSet := d.GetOk("ssh_keys") 
	if sshKeysSet {
		// process SSH Key settings and set API values accordingly
		log.Debugf("resourceComputeCreate: calling makeSshKeysArgString to setup SSH keys for guest login(s)")
		urlValues.Add("userdata", makeSshKeysArgString(sshKeysVal.([]interface{})))
	}

	computeCreateAPI := KvmX86CreateAPI
	arch := d.Get("arch").(string)
	if arch == "KVM_PPC" {
		computeCreateAPI = KvmPPCCreateAPI
		log.Debugf("resourceComputeCreate: creating Compute of type KVM VM PowerPC")
	} else { // note that we do not validate arch value for explicit "KVM_X86" here
		log.Debugf("resourceComputeCreate: creating Compute of type KVM VM x86")
	}
	
	apiResp, err := controller.decortAPICall("POST", computeCreateAPI, urlValues)
	if err != nil {
		return err
	}
	// Compute create API returns ID of the new Compute instance on success

	d.SetId(apiResp) // update ID of the resource to tell Terraform that the resource exists, albeit partially
	compId, _ := strconv.Atoi(apiResp)
	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("cpu")
	d.SetPartial("ram")
	d.SetPartial("image_id")
	d.SetPartial("boot_disk_size")
	if sshKeysSet {
		d.SetPartial("ssh_keys")
	}

	log.Debugf("resourceComputeCreate: new simple Compute ID %d, name %s created", compId, d.Get("name").(string))

	// Configure data disks if any
	extraDisksOk := true
	argVal, argSet = d.GetOk("extra_disks") 
	if argSet && len(argVal.([]interface{})) > 0 {
		// urlValues.Add("desc", argVal.(string))
		log.Debugf("resourceComputeCreate: calling utilityComputeExtraDisksConfigure to attach %d extra disk(s)", len(argVal.([]interface{})))
		err = controller.utilityComputeExtraDisksConfigure(d, false) // do_delta=false, as we are working on a new compute
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching extra disks to a new Compute ID %s: %s", compId, err)
			extraDisksOk = false
		}
	}
	if extraDisksOk {
		d.SetPartial("extra_disks")
	}

	// Configure external networks if any
	netsOk := true
	argVal, argSet = d.GetOk("network") 
	if argSet && len(argVal.([]interface{})) > 0 {
		log.Debugf("resourceComputeCreate: calling utilityComputeNetworksConfigure to attach %d network(s)", len(argVal.([]interface{})))
		err = controller.utilityComputeNetworksConfigure(d, false) // do_delta=false, as we are working on a new compute
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching networks to a new Compute ID %d: %s", compId, err)
			netsOk = false
		}
	}
	if netsOk {
		// there were no errors reported when configuring networks
		d.SetPartial("network")
	}

	if extraDisksOk && netsOk {
		// if there were no errors in setting any of the subresources, we may leave Partial mode
		d.Partial(false)
	}
	
	log.Debugf("resourceComputeCreate: new Compute ID %d, name %s creation sequence complete", compId, d.Get("name").(string))

	// We may reuse dataSourceComputeRead here as we maintain similarity 
	// between Compute resource and Compute data source schemas
	// Compute read function will also update resource ID on success, so that Terraform 
	// will know the resource exists
	return dataSourceComputeRead(d, m) 
}

func resourceComputeRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceComputeRead: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	compFacts, err := utilityComputeCheckPresence(d, m)
	if compFacts == "" {
		if err != nil {
			return err
		}
		// Compute with such name and RG ID was not found
		return nil
	}

	if err = flattenCompute(d, compFacts); err != nil {
		return err
	}

	log.Debugf("resourceComputeRead: after flattenCompute: Compute ID %s, name %q, RG ID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	return nil
}

func resourceComputeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceComputeUpdate: called for Compute ID %s / name %s, RGID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	controller := m.(*ControllerCfg)

	/* 
	1. Resize CPU/RAM
	2. Resize (grow) boot disk
	3. Update extra disks
	4. Update networks
	5. Update port forwards 
	*/

	// 1. Resize CPU/RAM
	params := &url.Values{}
	doUpdate := false
	params.Add("computeId", d.Id())

	oldCpu, newCpu := d.GetChange("cpu")
	if oldCpu.(int) != newCpu.(int) {
		params.Add("cpu", fmt.Sprintf("%d", newCpu.(int)))
		doUpdate = true
	} else {
		params.Add("cpu", "0") // no change to CPU allocation
	}
	
	oldRam, newRam := d.GetChange("ram")
	if oldRam.(int) != newRam.(int) {
		params.Add("ram", fmt.Sprintf("%d", newRam.(int)))
		doUpdate = true
	} else {
		params.Add("ram", "0")
	}

	if doUpdate {
		log.Debugf("resourceComputeUpdate: changing CPU %d -> %d and/or RAM %d -> %d",
		           oldCpu.(int), newCpu.(int),
				   oldRam.(int), newRam.(int))
		_, err := controller.decortAPICall("POST", ComputeResizeAPI, params)
		if err != nil {
			return err
		}
	}

	// 2. Resize (grow) Boot disk	
	oldSize, newSize := d.GetChange("boot_disk_size")
	if oldSize.(int) < newSize.(int) {
		bdsParams := &url.Values{}
		bdsParams.Add("diskId", fmt.Sprintf("%d", d.Get("boot_disk_id").(int)))
		bdsParams.Add("size", fmt.Sprintf("%d", newSize.(int)))
		log.Debugf("resourceComputeUpdate: compute ID %s, boot disk ID %d resize %d -> %d",
		           d.Id(), d.Get("boot_disk_id").(int), oldSize.(int), newSize.(int))
		_, err := controller.decortAPICall("POST", DisksResizeAPI, params)
		if err != nil {
			return err
		}
	} else if oldSize.(int) > newSize.(int) {
		log.Warnf("resourceComputeUpdate: compute ID %d - shrinking boot disk is not allowed", d.Id())
	}

	// 3. Calculate and apply changes to data disks
	err := controller.utilityComputeExtraDisksConfigure(d, true) // pass do_delta = true to apply changes, if any
	if err != nil {
		return err
	}

	// 4. Calculate and apply changes to network connections
	err = controller.utilityComputeNetworksConfigure(d, true) // pass do_delta = true to apply changes, if any
	if err != nil {
		return err
	}

	// we may reuse dataSourceComputeRead here as we maintain similarity 
	// between Compute resource and Compute data source schemas
	return dataSourceComputeRead(d, m) 
}

func resourceComputeDelete(d *schema.ResourceData, m interface{}) error {
	// NOTE: this function destroys target Compute instance "permanently", so 
	// there is no way to restore it. It also destroys all extra disks
	// attached to this compute, so "User, ye be warned!"
	log.Debugf("resourceComputeDelete: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

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
	_, err = controller.decortAPICall("POST", ComputeDeleteAPI, params)
	if err != nil {
		return err
	}

	return nil
}

func resourceComputeExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	log.Debugf("resourceComputeExist: called for Compute name %q, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

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
				ValidateFunc: validation.StringInSlice([]string{"KVM_X86", "KVM_PPC"}, false), // observe case while validating
				Description: "Hardware architecture of this compute instance.",
			},

			"cpu": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, MaxCpusPerCompute),
				Description:  "Number of CPUs to allocate to this compute instance.",
			},

			"ram": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(MinRamPerCompute),
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
				Description: "This compute instance boot disk size in GB.",
			},

			"extra_disks": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: MaxExtraDisksPerCompute,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Optional list of IDs of extra disks to attach to this compute. You may specify several extra disks.",
			},

			"network": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: MaxNetworksPerCompute,
				Elem: &schema.Resource{
					Schema: networkSubresourceSchemaMake(),
				},
				Description: "Optional network connection(s) for this compute. You may specify several network blocks, one for each connection.",
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

			"boot_disk_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This compute instance boot disk ID.",
			},

			/*
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
			*/
		},
	}
}

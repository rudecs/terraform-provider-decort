/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/statefuncs"
	"github.com/rudecs/terraform-provider-decort/internal/status"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func cloudInitDiffSupperss(key, oldVal, newVal string, d *schema.ResourceData) bool {
	if oldVal == "" && newVal != "applied" {
		// if old value for "cloud_init" resource is empty string, it means that we are creating new compute
		// and there is a chance that the user will want custom cloud init parameters - so we check if
		// cloud_init is explicitly set in TF file by making sure that its new value is different from "applied",
		// which is a reserved key word.
		log.Debugf("cloudInitDiffSupperss: key=%s, oldVal=%q, newVal=%q -> suppress=FALSE", key, oldVal, newVal)
		return false // there is a difference between stored and new value
	}
	log.Debugf("cloudInitDiffSupperss: key=%s, oldVal=%q, newVal=%q -> suppress=TRUE", key, oldVal, newVal)
	return true // suppress difference
}

func resourceComputeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// we assume all mandatory parameters it takes to create a comptue instance are properly
	// specified - we rely on schema "Required" attributes to let Terraform validate them for us

	log.Debugf("resourceComputeCreate: called for Compute name %q, RG ID %d", d.Get("name").(string), d.Get("rg_id").(int))

	// create basic Compute (i.e. without extra disks and network connections - those will be attached
	// by subsequent individual API calls).
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("rgId", fmt.Sprintf("%d", d.Get("rg_id").(int)))
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("cpu", fmt.Sprintf("%d", d.Get("cpu").(int)))
	urlValues.Add("ram", fmt.Sprintf("%d", d.Get("ram").(int)))
	urlValues.Add("imageId", fmt.Sprintf("%d", d.Get("image_id").(int)))
	urlValues.Add("netType", "NONE")
	urlValues.Add("start", "0") // at the 1st step create compute in a stopped state

	argVal, argSet := d.GetOk("description")
	if argSet {
		urlValues.Add("desc", argVal.(string))
	}

	if sepID, ok := d.GetOk("sep_id"); ok {
		urlValues.Add("sepId", strconv.Itoa(sepID.(int)))
	}

	if pool, ok := d.GetOk("pool"); ok {
		urlValues.Add("pool", pool.(string))
	}

	if ipaType, ok := d.GetOk("ipa_type"); ok {
		urlValues.Add("ipaType", ipaType.(string))
	}

	if bootSize, ok := d.GetOk("boot_disk_size"); ok {
		urlValues.Add("bootDisk", fmt.Sprintf("%d", bootSize.(int)))
	}

	if IS, ok := d.GetOk("is"); ok {
		urlValues.Add("IS", IS.(string))
	}
	if networks, ok := d.GetOk("network"); ok {
		if networks.(*schema.Set).Len() > 0 {
			ns := networks.(*schema.Set).List()
			defaultNetwork := ns[0].(map[string]interface{})
			urlValues.Set("netType", defaultNetwork["net_type"].(string))
			urlValues.Add("netId", fmt.Sprintf("%d", defaultNetwork["net_id"].(int)))
			ipaddr, ipSet := defaultNetwork["ip_address"] // "ip_address" key is optional
			if ipSet {
				urlValues.Add("ipAddr", ipaddr.(string))
			}

		}
	}

	/*
		sshKeysVal, sshKeysSet := d.GetOk("ssh_keys")
		if sshKeysSet {
			// process SSH Key settings and set API values accordingly
			log.Debugf("resourceComputeCreate: calling makeSshKeysArgString to setup SSH keys for guest login(s)")
			urlValues.Add("userdata", makeSshKeysArgString(sshKeysVal.([]interface{})))
		}
	*/

	computeCreateAPI := KvmX86CreateAPI
	driver := d.Get("driver").(string)
	if driver == "KVM_PPC" {
		computeCreateAPI = KvmPPCCreateAPI
		log.Debugf("resourceComputeCreate: creating Compute of type KVM VM PowerPC")
	} else { // note that we do not validate arch value for explicit "KVM_X86" here
		log.Debugf("resourceComputeCreate: creating Compute of type KVM VM x86")
	}

	argVal, argSet = d.GetOk("cloud_init")
	if argSet {
		// userdata must not be empty string and must not be a reserved keyword "applied"
		userdata := argVal.(string)
		if userdata != "" && userdata != "applied" {
			urlValues.Add("userdata", userdata)
		}
	}

	apiResp, err := c.DecortAPICall(ctx, "POST", computeCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	urlValues = &url.Values{}
	// Compute create API returns ID of the new Compute instance on success

	d.SetId(apiResp) // update ID of the resource to tell Terraform that the resource exists, albeit partially
	compId, _ := strconv.Atoi(apiResp)

	cleanup := false
	defer func() {
		if cleanup {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("permanently", "1")
			urlValues.Add("detachDisks", "1")

			if _, err := c.DecortAPICall(ctx, "POST", ComputeDeleteAPI, urlValues); err != nil {
				log.Errorf("resourceComputeCreate: could not delete compute after failed creation: %v", err)
			}
			d.SetId("")
			urlValues = &url.Values{}
		}
	}()

	log.Debugf("resourceComputeCreate: new simple Compute ID %d, name %s created", compId, d.Get("name").(string))

	// Configure data disks if any
	argVal, argSet = d.GetOk("extra_disks")
	if argSet && argVal.(*schema.Set).Len() > 0 {
		// urlValues.Add("desc", argVal.(string))
		log.Debugf("resourceComputeCreate: calling utilityComputeExtraDisksConfigure to attach %d extra disk(s)", argVal.(*schema.Set).Len())
		err = utilityComputeExtraDisksConfigure(ctx, d, m, false) // do_delta=false, as we are working on a new compute
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching extra disk(s) to a new Compute ID %d: %v", compId, err)
			cleanup = true
			return diag.FromErr(err)
		}
	}
	// Configure external networks if any
	argVal, argSet = d.GetOk("network")
	if argSet && argVal.(*schema.Set).Len() > 0 {
		log.Debugf("resourceComputeCreate: calling utilityComputeNetworksConfigure to attach %d network(s)", argVal.(*schema.Set).Len())
		err = utilityComputeNetworksConfigure(ctx, d, m, false, true) // do_delta=false, as we are working on a new compute
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching networks to a new Compute ID %d: %s", compId, err)
			cleanup = true
			return diag.FromErr(err)
		}
	}

	// Note bene: we created compute in a STOPPED state (this is required to properly attach 1st network interface),
	// now we need to start it before we report the sequence complete
	if d.Get("started").(bool) {
		reqValues := &url.Values{}
		reqValues.Add("computeId", fmt.Sprintf("%d", compId))
		log.Debugf("resourceComputeCreate: starting Compute ID %d after completing its resource configuration", compId)
		if _, err := c.DecortAPICall(ctx, "POST", ComputeStartAPI, reqValues); err != nil {
			cleanup = true
			return diag.FromErr(err)
		}
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		api := ComputeDisableAPI
		if enabled.(bool) {
			api = ComputeEnableAPI
		}
		urlValues := &url.Values{}
		urlValues.Add("computeId", fmt.Sprintf("%d", compId))
		log.Debugf("resourceComputeCreate: enable=%t Compute ID %d after completing its resource configuration", compId, enabled)
		if _, err := c.DecortAPICall(ctx, "POST", api, urlValues); err != nil {
			return diag.FromErr(err)
		}

	}

	if !cleanup {
		if disks, ok := d.GetOk("disks"); ok {
			log.Debugf("resourceComputeCreate: Create disks on ComputeID: %d", compId)
			addedDisks := disks.([]interface{})
			if len(addedDisks) > 0 {
				for _, disk := range addedDisks {
					diskConv := disk.(map[string]interface{})

					urlValues.Add("computeId", d.Id())
					urlValues.Add("diskName", diskConv["disk_name"].(string))
					urlValues.Add("size", strconv.Itoa(diskConv["size"].(int)))
					if diskConv["disk_type"].(string) != "" {
						urlValues.Add("diskType", diskConv["disk_type"].(string))
					}
					if diskConv["sep_id"].(int) != 0 {
						urlValues.Add("sepId", strconv.Itoa(diskConv["sep_id"].(int)))
					}
					if diskConv["pool"].(string) != "" {
						urlValues.Add("pool", diskConv["pool"].(string))
					}
					if diskConv["desc"].(string) != "" {
						urlValues.Add("desc", diskConv["desc"].(string))
					}
					if diskConv["image_id"].(int) != 0 {
						urlValues.Add("imageId", strconv.Itoa(diskConv["image_id"].(int)))
					}
					_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAddAPI, urlValues)
					if err != nil {
						cleanup = true
						return diag.FromErr(err)
					}
					urlValues = &url.Values{}
				}
			}
		}
	}

	log.Debugf("resourceComputeCreate: new Compute ID %d, name %s creation sequence complete", compId, d.Get("name").(string))

	// We may reuse dataSourceComputeRead here as we maintain similarity
	// between Compute resource and Compute data source schemas
	// Compute read function will also update resource ID on success, so that Terraform
	// will know the resource exists
	return resourceComputeRead(ctx, d, m)
}

func resourceComputeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceComputeRead: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)

	compFacts, err := utilityComputeCheckPresence(ctx, d, m)
	if compFacts == "" {
		if err != nil {
			return diag.FromErr(err)
		}
		// Compute with such name and RG ID was not found
		return nil
	}

	compute := &ComputeGetResp{}
	err = json.Unmarshal([]byte(compFacts), compute)

	log.Debugf("resourceComputeRead: compute is: %+v", compute)
	if err != nil {
		return diag.FromErr(err)
	}

	switch compute.Status {
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", ComputeRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall(ctx, "POST", ComputeEnableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	case status.Destroyed:
		d.SetId("")
		return resourceComputeCreate(ctx, d, m)
	case status.Disabled:
		log.Debugf("The compute is in status: %s, may troubles can be occured with update. Please, enable compute first.", compute.Status)
	case status.Redeploying:
	case status.Deleting:
	case status.Destroying:
		return diag.Errorf("The compute is in progress with status: %s", compute.Status)
	case status.Modeled:
		return diag.Errorf("The compute is in status: %s, please, contant the support for more information", compute.Status)
	}

	compFacts, err = utilityComputeCheckPresence(ctx, d, m)
	log.Debugf("resourceComputeRead: after changes compute is: %s", compFacts)
	if compFacts == "" {
		if err != nil {
			return diag.FromErr(err)
		}
		// Compute with such name and RG ID was not found
		return nil
	}

	if err = flattenCompute(d, compFacts); err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("resourceComputeRead: after flattenCompute: Compute ID %s, name %q, RG ID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	return nil
}

func resourceComputeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceComputeUpdate: called for Compute ID %s / name %s, RGID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)

	computeRaw, err := utilityComputeCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	compute := &ComputeGetResp{}
	err = json.Unmarshal([]byte(computeRaw), compute)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled")
		api := ComputeDisableAPI
		if enabled.(bool) {
			api = ComputeEnableAPI
		}
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		log.Debugf("resourceComputeUpdate: enable=%t Compute ID %s after completing its resource configuration", d.Id(), enabled)
		if _, err := c.DecortAPICall(ctx, "POST", api, urlValues); err != nil {
			return diag.FromErr(err)
		}
	}

	// check compute statuses
	switch compute.Status {
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", ComputeRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall(ctx, "POST", ComputeEnableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	case status.Destroyed:
		d.SetId("")
		return resourceComputeCreate(ctx, d, m)
	case status.Disabled:
		log.Debugf("The compute is in status: %s, may troubles can be occured with update. Please, enable compute first.", compute.Status)
	case status.Redeploying:
	case status.Deleting:
	case status.Destroying:
		return diag.Errorf("The compute is in progress with status: %s", compute.Status)
	case status.Modeled:
		return diag.Errorf("The compute is in status: %s, please, contant the support for more information", compute.Status)
	}

	/*
		1. Resize CPU/RAM
		2. Resize (grow) boot disk
		3. Update extra disks
		4. Update networks
		5. Start/stop
	*/

	// 1. Resize CPU/RAM
	urlValues := &url.Values{}
	doUpdate := false
	urlValues.Add("computeId", d.Id())

	oldCpu, newCpu := d.GetChange("cpu")
	if oldCpu.(int) != newCpu.(int) {
		urlValues.Add("cpu", fmt.Sprintf("%d", newCpu.(int)))
		doUpdate = true
	} else {
		urlValues.Add("cpu", "0") // no change to CPU allocation
	}

	oldRam, newRam := d.GetChange("ram")
	if oldRam.(int) != newRam.(int) {
		urlValues.Add("ram", fmt.Sprintf("%d", newRam.(int)))
		doUpdate = true
	} else {
		urlValues.Add("ram", "0")
	}

	if doUpdate {
		log.Debugf("resourceComputeUpdate: changing CPU %d -> %d and/or RAM %d -> %d",
			oldCpu.(int), newCpu.(int),
			oldRam.(int), newRam.(int))
		urlValues.Add("force", "true")
		_, err := c.DecortAPICall(ctx, "POST", ComputeResizeAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
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
		_, err := c.DecortAPICall(ctx, "POST", DisksResizeAPI, bdsParams)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if oldSize.(int) > newSize.(int) {
		log.Warnf("resourceComputeUpdate: compute ID %s - shrinking boot disk is not allowed", d.Id())
	}

	// 3. Calculate and apply changes to data disks
	if d.HasChange("extra_disks") {
		err := utilityComputeExtraDisksConfigure(ctx, d, m, true) // pass do_delta = true to apply changes, if any
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// 4. Calculate and apply changes to network connections
	err = utilityComputeNetworksConfigure(ctx, d, m, true, false) // pass do_delta = true to apply changes, if any
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("description") || d.HasChange("name") {
		updateParams := &url.Values{}
		updateParams.Add("computeId", d.Id())
		updateParams.Add("name", d.Get("name").(string))
		updateParams.Add("desc", d.Get("description").(string))
		if _, err := c.DecortAPICall(ctx, "POST", ComputeUpdateAPI, updateParams); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("started") {
		params := &url.Values{}
		params.Add("computeId", d.Id())
		if d.Get("started").(bool) {
			if _, err := c.DecortAPICall(ctx, "POST", ComputeStartAPI, params); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if _, err := c.DecortAPICall(ctx, "POST", ComputeStopAPI, params); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("disks") {
		deletedDisks := make([]interface{}, 0)
		addedDisks := make([]interface{}, 0)

		oldDisks, newDisks := d.GetChange("disks")
		oldConv := oldDisks.([]interface{})
		newConv := newDisks.([]interface{})

		for _, el := range oldConv {
			if !isContainsDisk(newConv, el) {
				deletedDisks = append(deletedDisks, el)
			}
		}

		for _, el := range newConv {
			if !isContainsDisk(oldConv, el) {
				addedDisks = append(addedDisks, el)
			}
		}

		if len(deletedDisks) > 0 {
			urlValues.Add("computeId", d.Id())
			urlValues.Add("force", "false")
			_, err := c.DecortAPICall(ctx, "POST", ComputeStopAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
			urlValues = &url.Values{}

			for _, disk := range deletedDisks {
				diskConv := disk.(map[string]interface{})
				if diskConv["disk_name"].(string) == "bootdisk" {
					continue
				}
				urlValues.Add("computeId", d.Id())
				urlValues.Add("diskId", strconv.Itoa(diskConv["disk_id"].(int)))
				urlValues.Add("permanently", strconv.FormatBool(diskConv["permanently"].(bool)))
				_, err := c.DecortAPICall(ctx, "POST", ComputeDiskDeleteAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("altBootId", "0")
			_, err = c.DecortAPICall(ctx, "POST", ComputeStartAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
			urlValues = &url.Values{}
		}

		if len(addedDisks) > 0 {
			for _, disk := range addedDisks {
				diskConv := disk.(map[string]interface{})
				if diskConv["disk_name"].(string) == "bootdisk" {
					continue
				}
				urlValues.Add("computeId", d.Id())
				urlValues.Add("diskName", diskConv["disk_name"].(string))
				urlValues.Add("size", strconv.Itoa(diskConv["size"].(int)))
				if diskConv["disk_type"].(string) != "" {
					urlValues.Add("diskType", diskConv["disk_type"].(string))
				}
				if diskConv["sep_id"].(int) != 0 {
					urlValues.Add("sepId", strconv.Itoa(diskConv["sep_id"].(int)))
				}
				if diskConv["pool"].(string) != "" {
					urlValues.Add("pool", diskConv["pool"].(string))
				}
				if diskConv["desc"].(string) != "" {
					urlValues.Add("desc", diskConv["desc"].(string))
				}
				if diskConv["image_id"].(int) != 0 {
					urlValues.Add("imageId", strconv.Itoa(diskConv["image_id"].(int)))
				}
				_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAddAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}
	}

	// we may reuse dataSourceComputeRead here as we maintain similarity
	// between Compute resource and Compute data source schemas
	return resourceComputeRead(ctx, d, m)
}

func isContainsDisk(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["disk_name"].(string) == elConv["disk_name"].(string) {
			return true
		}
	}
	return false
}

func resourceComputeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// NOTE: this function destroys target Compute instance "permanently", so
	// there is no way to restore it.
	// If compute being destroyed has some extra disks attached, they are
	// detached from the compute
	log.Debugf("resourceComputeDelete: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)

	params := &url.Values{}
	params.Add("computeId", d.Id())
	params.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))
	params.Add("detachDisks", strconv.FormatBool(d.Get("detach_disks").(bool)))

	if _, err := c.DecortAPICall(ctx, "POST", ComputeDeleteAPI, params); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceComputeSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of this compute. Compute names are case sensitive and must be unique in the resource group.",
		},

		"rg_id": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "ID of the resource group where this compute should be deployed.",
		},

		"driver": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			StateFunc:    statefuncs.StateFuncToUpper,
			ValidateFunc: validation.StringInSlice([]string{"KVM_X86", "KVM_PPC"}, false), // observe case while validating
			Description:  "Hardware architecture of this compute instance.",
		},

		"cpu": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, constants.MaxCpusPerCompute),
			Description:  "Number of CPUs to allocate to this compute instance.",
		},

		"ram": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(constants.MinRamPerCompute),
			Description:  "Amount of RAM in MB to allocate to this compute instance.",
		},

		"image_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the OS image to base this compute instance on.",
		},

		"boot_disk_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "This compute instance boot disk size in GB. Make sure it is large enough to accomodate selected OS image.",
		},

		"disks": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name for disk",
					},
					"size": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "Disk size in GiB",
					},
					"disk_type": {
						Type:         schema.TypeString,
						Computed:     true,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"B", "D"}, false),
						Description:  "The type of disk in terms of its role in compute: 'B=Boot, D=Data'",
					},
					"sep_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Optional:    true,
						Description: "Storage endpoint provider ID; by default the same with boot disk",
					},
					"pool": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Pool name; by default will be chosen automatically",
					},
					"desc": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Optional description",
					},
					"image_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Optional:    true,
						Description: "Specify image id for create disk from template",
					},
					"disk_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Disk ID",
					},
					"permanently": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Disk deletion status",
					},
				},
			},
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "ID of SEP to create bootDisk on. Uses image's sepId if not set.",
		},

		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Pool to use if sepId is set, can be also empty if needed to be chosen by system.",
		},

		"extra_disks": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: constants.MaxExtraDisksPerCompute,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "Optional list of IDs of extra disks to attach to this compute. You may specify several extra disks.",
		},

		"network": {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			MaxItems: constants.MaxNetworksPerCompute,
			Elem: &schema.Resource{
				Schema: networkSubresourceSchemaMake(),
			},
			Description: "Optional network connection(s) for this compute. You may specify several network blocks, one for each connection.",
		},

		/*
			"ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: MaxSshKeysPerCompute,
				Elem: &schema.Resource{
					Schema: sshSubresourceSchemaMake(),
				},
				Description: "SSH keys to authorize on this compute instance.",
			},
		*/

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional text description of this compute instance.",
		},

		"cloud_init": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional cloud_init parameters. Applied when creating new compute instance only, ignored in all other cases.",
			//Default:          "applied",
			//DiffSuppressFunc: cloudInitDiffSupperss,
		},

		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "If true - enable compute, else - disable",
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

		"os_users": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: osUsersSubresourceSchemaMake(),
			},
			Description: "Guest OS users provisioned on this compute instance.",
		},

		"started": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Is compute started.",
		},
		"detach_disks": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"permanently": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"is": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "system name",
		},
		"ipa_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "compute purpose",
		},
	}
	return rets
}

func ResourceCompute() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceComputeCreate,
		ReadContext:   resourceComputeRead,
		UpdateContext: resourceComputeUpdate,
		DeleteContext: resourceComputeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout600s,
			Read:    &constants.Timeout300s,
			Update:  &constants.Timeout300s,
			Delete:  &constants.Timeout300s,
			Default: &constants.Timeout300s,
		},

		Schema: ResourceComputeSchemaMake(),
	}
}

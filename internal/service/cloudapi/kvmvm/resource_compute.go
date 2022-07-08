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

	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/statefuncs"
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
	urlValues.Add("bootDisk", fmt.Sprintf("%d", d.Get("boot_disk_size").(int)))
	urlValues.Add("netType", "NONE") // at the 1st step create isolated compute
	urlValues.Add("start", "0")      // at the 1st step create compute in a stopped state

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
	// Compute create API returns ID of the new Compute instance on success

	d.SetId(apiResp) // update ID of the resource to tell Terraform that the resource exists, albeit partially
	compId, _ := strconv.Atoi(apiResp)

	log.Debugf("resourceComputeCreate: new simple Compute ID %d, name %s created", compId, d.Get("name").(string))

	// Configure data disks if any
	argVal, argSet = d.GetOk("extra_disks")
	if argSet && argVal.(*schema.Set).Len() > 0 {
		// urlValues.Add("desc", argVal.(string))
		log.Debugf("resourceComputeCreate: calling utilityComputeExtraDisksConfigure to attach %d extra disk(s)", argVal.(*schema.Set).Len())
		err = utilityComputeExtraDisksConfigure(ctx, d, m, false) // do_delta=false, as we are working on a new compute
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching extra disk(s) to a new Compute ID %d: %v", compId, err)
			return diag.FromErr(err)
		}
	}
	// Configure external networks if any
	argVal, argSet = d.GetOk("network")
	if argSet && argVal.(*schema.Set).Len() > 0 {
		log.Debugf("resourceComputeCreate: calling utilityComputeNetworksConfigure to attach %d network(s)", argVal.(*schema.Set).Len())
		err = utilityComputeNetworksConfigure(ctx, d, m, false) // do_delta=false, as we are working on a new compute
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching networks to a new Compute ID %d: %s", compId, err)
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
			return diag.FromErr(err)
		}
	}

	log.Debugf("resourceComputeCreate: new Compute ID %d, name %s creation sequence complete", compId, d.Get("name").(string))

	// We may reuse dataSourceComputeRead here as we maintain similarity
	// between Compute resource and Compute data source schemas
	// Compute read function will also update resource ID on success, so that Terraform
	// will know the resource exists
	return dataSourceComputeRead(ctx, d, m)
}

func resourceComputeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceComputeRead: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	compFacts, err := utilityComputeCheckPresence(ctx, d, m)
	if compFacts == "" {
		if err != nil {
			return diag.FromErr(err)
		}
		// Compute with such name and RG ID was not found
		return nil
	}

	if diagnostic := flattenCompute(d, compFacts); diagnostic != nil {
		return diagnostic
	}

	log.Debugf("resourceComputeRead: after flattenCompute: Compute ID %s, name %q, RG ID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	return nil
}

func resourceComputeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceComputeUpdate: called for Compute ID %s / name %s, RGID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)

	/*
		1. Resize CPU/RAM
		2. Resize (grow) boot disk
		3. Update extra disks
		4. Update networks
		5. Start/stop
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
		params.Add("force", "true")
		_, err := c.DecortAPICall(ctx, "POST", ComputeResizeAPI, params)
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
	err := utilityComputeExtraDisksConfigure(ctx, d, m, true) // pass do_delta = true to apply changes, if any
	if err != nil {
		return diag.FromErr(err)
	}

	// 4. Calculate and apply changes to network connections
	err = utilityComputeNetworksConfigure(ctx, d, m, true) // pass do_delta = true to apply changes, if any
	if err != nil {
		return diag.FromErr(err)
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

	// we may reuse dataSourceComputeRead here as we maintain similarity
	// between Compute resource and Compute data source schemas
	return dataSourceComputeRead(ctx, d, m)
}

func resourceComputeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// NOTE: this function destroys target Compute instance "permanently", so
	// there is no way to restore it.
	// If compute being destroyed has some extra disks attached, they are
	// detached from the compute
	log.Debugf("resourceComputeDelete: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	compFacts, err := utilityComputeCheckPresence(ctx, d, m)
	if compFacts == "" {
		if err != nil {
			return diag.FromErr(err)
		}
		// the target Compute does not exist - in this case according to Terraform best practice
		// we exit from Destroy method without error
		return nil
	}

	c := m.(*controller.ControllerCfg)

	model := ComputeGetResp{}
	log.Debugf("resourceComputeDelete: ready to unmarshal string %s", compFacts)
	err = json.Unmarshal([]byte(compFacts), &model)
	if err == nil && len(model.Disks) > 0 {
		// prepare to detach data disks from compute - do it only if compFacts unmarshalled
		// properly and the resulting model contains non-empty Disks list
		for _, diskFacts := range model.Disks {
			if diskFacts.Type == "B" {
				// boot disk is never detached on compute delete
				continue
			}

			log.Debugf("resourceComputeDelete: ready to detach data disk ID %d from compute ID %s", diskFacts.ID, d.Id())

			detachParams := &url.Values{}
			detachParams.Add("computeId", d.Id())
			detachParams.Add("diskId", fmt.Sprintf("%d", diskFacts.ID))

			_, err = c.DecortAPICall(ctx, "POST", ComputeDiskDetachAPI, detachParams)
			if err != nil {
				// We do not fail compute deletion on data disk detach errors
				log.Errorf("resourceComputeDelete: error when detaching Disk ID %d: %s", diskFacts.ID, err)
			}
		}
	}

	params := &url.Values{}
	params.Add("computeId", d.Id())
	params.Add("permanently", "1")
	// TODO: this is for the upcoming API update - params.Add("detachdisks", "1")

	_, err = c.DecortAPICall(ctx, "POST", ComputeDeleteAPI, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceComputeExists(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	log.Debugf("resourceComputeExist: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	compFacts, err := utilityComputeCheckPresence(ctx, d, m)
	if compFacts == "" {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func ResourceCompute() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceComputeCreate,
		ReadContext:   resourceComputeRead,
		UpdateContext: resourceComputeUpdate,
		DeleteContext: resourceComputeDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout180s,
			Read:    &constants.Timeout30s,
			Update:  &constants.Timeout180s,
			Delete:  &constants.Timeout60s,
			Default: &constants.Timeout60s,
		},

		Schema: map[string]*schema.Schema{
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
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the OS image to base this compute instance on.",
			},

			"boot_disk_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "This compute instance boot disk size in GB. Make sure it is large enough to accomodate selected OS image.",
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "applied",
				DiffSuppressFunc: cloudInitDiffSupperss,
				Description:      "Optional cloud_init parameters. Applied when creating new compute instance only, ignored in all other cases.",
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
				Default:     true,
				Description: "Is compute started.",
			},
		},
	}
}

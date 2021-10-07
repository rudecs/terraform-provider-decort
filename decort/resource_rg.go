/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions. All Rights Reserved.

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
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceResgroupCreate(d *schema.ResourceData, m interface{}) error {
	// First validate that we have all parameters required to create the new Resource Group

	// Valid account ID is required to create new resource group
	// obtain Account ID by account name - it should not be zero on success
	validated_account_id, err := utilityGetAccountIdBySchema(d, m)
	if err != nil {
		return err
	}

	rg_name, arg_set := d.GetOk("name")
	if !arg_set {
		return fmt.Errorf("Cannot create new RG: missing name.")
	}

	/* Current version of provider works with default grid id (same is true for disk resources)
	grid_id, arg_set := d.GetOk("grid_id")
	if !arg_set {
		return fmt.Errorf("Cannot create new RG %q in account ID %d: missing Grid ID.",
			rg_name.(string), validated_account_id)
	}
	if grid_id.(int) < 1 {
		grid_id = DefaultGridID
	}
	*/

	// all required parameters are set in the schema - we can continue with RG creation
	log.Debugf("resourceResgroupCreate: called for RG name %s, account ID %d",
		rg_name.(string), validated_account_id)

	// quota settings are optional
	set_quota := false
	var quota_record QuotaRecord
	arg_value, arg_set := d.GetOk("quota")
	if arg_set {
		log.Debugf("resourceResgroupCreate: setting Quota on RG requested")
		quota_record, _ = makeQuotaRecord(arg_value.([]interface{}))
		set_quota = true
	}

	controller := m.(*ControllerCfg)
	log.Debugf("resourceResgroupCreate: called by user %q for RG name %s, account ID %d",
		controller.getDecortUsername(),
		rg_name.(string), validated_account_id)

	url_values := &url.Values{}
	url_values.Add("accountId", fmt.Sprintf("%d", validated_account_id))
	url_values.Add("name", rg_name.(string))
	url_values.Add("gid", fmt.Sprintf("%d", DefaultGridID)) // use default Grid ID, similar to disk resource mgmt convention
	url_values.Add("owner", controller.getDecortUsername())

	// pass quota values as set
	if set_quota {
		url_values.Add("maxCPUCapacity", fmt.Sprintf("%d", quota_record.Cpu))
		url_values.Add("maxVDiskCapacity", fmt.Sprintf("%d", quota_record.Disk))
		url_values.Add("maxMemoryCapacity", fmt.Sprintf("%f", quota_record.Ram)) // RAM quota is float; this may change in the future
		url_values.Add("maxNetworkPeerTransfer", fmt.Sprintf("%d", quota_record.ExtTraffic))
		url_values.Add("maxNumPublicIP", fmt.Sprintf("%d", quota_record.ExtIPs))
		// url_values.Add("???", fmt.Sprintf("%d", quota_record.GpuUnits))
	}

	// parse and handle network settings
	def_net_type, arg_set := d.GetOk("def_net_type")
	if arg_set {
		url_values.Add("def_net", def_net_type.(string)) // NOTE: in API default network type is set by "def_net" parameter
	}

	ipcidr, arg_set := d.GetOk("ipcidr")
	if arg_set {
		url_values.Add("ipcidr", ipcidr.(string))
	}

	ext_net_id, arg_set := d.GetOk("ext_net_id")
	if arg_set {
		url_values.Add("extNetId", fmt.Sprintf("%d", ext_net_id.(int)))
	}

	ext_ip, arg_set := d.GetOk("ext_ip")
	if arg_set {
		url_values.Add("extIp", ext_ip.(string))
	}

	api_resp, err, _ := controller.decortAPICall("POST", ResgroupCreateAPI, url_values)
	if err != nil {
		return err
	}

	d.SetId(api_resp) // rg/create API returns ID of the newly creted resource group on success
	// rg.ID, _ = strconv.Atoi(api_resp)

	// re-read newly created RG to make sure schema contains complete and up to date set of specifications
	return resourceResgroupRead(d, m)
}

func resourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceResgroupRead: called for RG name %s, account ID %s",
		d.Get("name").(string), d.Get("account_id").(int))
		
	rg_facts, err := utilityResgroupCheckPresence(d, m)
	if rg_facts == "" {
		// if empty string is returned from utilityResgroupCheckPresence then there is no
		// such resource group and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty
		return err
	}

	return flattenResgroup(d, rg_facts)
}

func resourceResgroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceResgroupUpdate: called for RG name %s, account ID %d",
		d.Get("name").(string), d.Get("account_id").(int))

	/* NOTE: we do not allow changing the following attributes of an existing RG via terraform:
	   - def_net_type
	   - ipcidr
	   - ext_net_id
	   - ext_ip

	   The following code fragment checks if any of these have been changed and generates error.
	*/
	for _, attr := range []string{"def_net_type", "ipcidr", "ext_ip"} {
		attr_new, attr_old := d.GetChange("def_net_type")
		if attr_new.(string) != attr_old.(string) {
			return fmt.Errorf("resourceResgroupUpdate: RG ID %s: changing %s for existing RG is not allowed", d.Id(), attr)
		}
	}

	attr_new, attr_old := d.GetChange("ext_net_id")
	if attr_new.(int) != attr_old.(int) {
		return fmt.Errorf("resourceResgroupUpdate: RG ID %s: changing ext_net_id for existing RG is not allowed", d.Id())
	}
	

	do_general_update := false      // will be true if general RG update is necessary (API rg/update)

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	url_values.Add("rgId", d.Id())

	name_new, name_set := d.GetOk("name")
	if name_set {
		log.Debugf("resourceResgroupUpdate: name specified - looking for deltas from the old settings.")
		name_old, _ := d.GetChange("name")
		if name_old.(string) != name_new.(string) {
			do_general_update = true
			url_values.Add("name", name_new.(string))
		}
	}

	quota_value, quota_set := d.GetOk("quota")
	if quota_set {
		log.Debugf("resourceResgroupUpdate: quota specified - looking for deltas from the old quota.")
		quotarecord_new, _ := makeQuotaRecord(quota_value.([]interface{}))
		quota_value_old, _ := d.GetChange("quota") // returns old as 1st, new as 2nd return value
		quotarecord_old, _ := makeQuotaRecord(quota_value_old.([]interface{}))

		if quotarecord_new.Cpu != quotarecord_old.Cpu {
			do_general_update = true
			log.Debugf("resourceResgroupUpdate: Cpu diff %d <- %d", quotarecord_new.Cpu, quotarecord_old.Cpu)
			url_values.Add("maxCPUCapacity", fmt.Sprintf("%d", quotarecord_new.Cpu))
		}

		if quotarecord_new.Disk != quotarecord_old.Disk {
			do_general_update = true
			log.Debugf("resourceResgroupUpdate: Disk diff %d <- %d", quotarecord_new.Disk, quotarecord_old.Disk)
			url_values.Add("maxVDiskCapacity", fmt.Sprintf("%d", quotarecord_new.Disk))
		}

		if quotarecord_new.Ram != quotarecord_old.Ram { // NB: quota on RAM is stored as float32, in units of MB
			do_general_update = true
			log.Debugf("resourceResgroupUpdate: Ram diff %f <- %f", quotarecord_new.Ram, quotarecord_old.Ram)
			url_values.Add("maxMemoryCapacity", fmt.Sprintf("%f", quotarecord_new.Ram))
		}

		if quotarecord_new.ExtTraffic != quotarecord_old.ExtTraffic {
			do_general_update = true
			log.Debugf("resourceResgroupUpdate: ExtTraffic diff %d <- %d", quotarecord_new.ExtTraffic, quotarecord_old.ExtTraffic)
			url_values.Add("maxNetworkPeerTransfer", fmt.Sprintf("%d", quotarecord_new.ExtTraffic))
		}

		if quotarecord_new.ExtIPs != quotarecord_old.ExtIPs {
			do_general_update = true
			log.Debugf("resourceResgroupUpdate: ExtIPs diff %d <- %d", quotarecord_new.ExtIPs, quotarecord_old.ExtIPs)
			url_values.Add("maxNumPublicIP", fmt.Sprintf("%d", quotarecord_new.ExtIPs))
		}
	}

	desc_new, desc_set := d.GetOk("description")
	if desc_set {
		log.Debugf("resourceResgroupUpdate: description specified - looking for deltas from the old settings.")
		desc_old, _ := d.GetChange("description")
		if desc_old.(string) != desc_new.(string) {
			do_general_update = true
			url_values.Add("desc", desc_new.(string))
		}
	}

	if do_general_update {
		log.Debugf("resourceResgroupUpdate: detected delta between new and old RG specs - updating the RG")
		_, err, _ := controller.decortAPICall("POST", ResgroupUpdateAPI, url_values)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("resourceResgroupUpdate: no difference between old and new state - no update on the RG will be done")
	}

	return resourceResgroupRead(d, m)
}

func resourceResgroupDelete(d *schema.ResourceData, m interface{}) error {
	// NOTE: this method forcibly destroys target resource group with flag "permanently", so there is no way to
	// restore the destroyed resource group as well all Computes & VINSes that existed in it
	log.Debugf("resourceResgroupDelete: called for RG name %s, account ID %s",
		d.Get("name").(string), d.Get("account_id").(int))

	rg_facts, err := utilityResgroupCheckPresence(d, m)
	if rg_facts == "" {
		// the target RG does not exist - in this case according to Terraform best practice
		// we exit from Destroy method without error
		return nil
	}

	url_values := &url.Values{}
	url_values.Add("rgId", d.Id())
	url_values.Add("force", "1")
	url_values.Add("permanently", "1")
	url_values.Add("reason", "Destroyed by DECORT Terraform provider")

	controller := m.(*ControllerCfg)
	_, err, _ = controller.decortAPICall("POST", ResgroupDeleteAPI, url_values)
	if err != nil {
		return err
	}

	return nil
}

func resourceResgroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should NOT modify ResourceData argument
	rg_facts, err := utilityResgroupCheckPresence(d, m)
	if rg_facts == "" {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func resourceResgroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceResgroupCreate,
		Read:   resourceResgroupRead,
		Update: resourceResgroupUpdate,
		Delete: resourceResgroupDelete,
		Exists: resourceResgroupExists,

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
				Description: "Name of this resource group. Names are case sensitive and unique within the context of a account.",
			},

			"account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ValidateFunc: validation.IntAtLeast(1),
				Description: "Unique ID of the account, which this resource group belongs to.",
			},

			"def_net_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "PRIVATE",
				ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC", "NONE"}, false),
				Description: "Type of the network, which this resource group will use as default for its computes - PRIVATE or PUBLIC or NONE.",
			},

			"def_net_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the default network for this resource group (if any).",
			},

			"ipcidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the netowrk inside the private network segment (aka ViNS) if def_net_type=PRIVATE",
			},

			"ext_net_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the external network for default ViNS. Pass 0 if def_net_type=PUBLIC or no external connection required for the defult ViNS when def_net_type=PRIVATE",
			},

			"ext_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address on the external netowrk to request when def_net_type=PRIVATE and ext_net_id is not 0",
			},

			/* commented out, as in this version of provider we use default Grid ID
			"grid_id": { 
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,     // if 0 is passed, default Grid ID will be used
				// DefaultFunc: utilityResgroupGetDefaultGridID,
				ForceNew:    true,  // change of Grid ID will require new RG
				Description: "Unique ID of the grid, where this resource group is deployed.",
			},
			*/

			"quota": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: quotaRgSubresourceSchemaMake(),
				},
				Description: "Quota settings for this resource group.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User-defined text description of this resource group.",
			},

			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the account, which this resource group belongs to.",
			},

			/*
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of this resource group.",
			},

			"vins": {
				Type:     schema.TypeList, // this is a list of ints
				Computed: true,
				MaxItems: LimitMaxVinsPerResgroup,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of VINs deployed in this resource group.",
			},

			"computes": {
				Type:     schema.TypeList, // this is a list of ints
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of computes deployed in this resource group.",
			},
			*/
		},
	}
}

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
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ipcidrDiffSupperss(key, oldVal, newVal string, d *schema.ResourceData) bool {
	if oldVal == "" && newVal != "" {
		// if old value for "ipcidr" resource is empty string, it means that we are creating new ViNS
		// and there is a chance that the user will want specific IP address range for this ViNS -  
		// check if "ipcidr" is explicitly set in TF file to a non-empty string.
		log.Debugf("ipcidrDiffSupperss: key=%s, oldVal=%q, newVal=%q -> suppress=FALSE", key, oldVal, newVal)
		return false // there is a difference between stored and new value
	}
	log.Debugf("ipcidrDiffSupperss: key=%s, oldVal=%q, newVal=%q -> suppress=TRUE", key, oldVal, newVal)
	return true // suppress difference
}

func resourceVinsCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceVinsCreate: called for ViNS name %s, Account ID %d, RG ID %d",
	           d.Get("name").(string), d.Get("account_id").(int), d.Get("rg_id").(int))

	apiToCall := VinsCreateInAccountAPI

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("name", d.Get("name").(string))

	argVal, argSet := d.GetOk("rg_id")
	if argSet && argVal.(int) > 0 {
		apiToCall = VinsCreateInRgAPI
		urlValues.Add("rgId", fmt.Sprintf("%d", argVal.(int)))
	} else {
		// RG ID either not set at all or set to 0 - user may want ViNS at account level 
		argVal, argSet = d.GetOk("account_id")
		if !argSet || argVal.(int) <= 0 {
			// No valid Account ID (and no RG ID either) - cannot create ViNS
			return fmt.Errorf("resourceVinsCreate: ViNS name %s - no valid account and/or resource group ID specified", d.Id())
		}
		urlValues.Add("accountId", fmt.Sprintf("%d", argVal.(int)))
	}

	argVal, argSet = d.GetOk("ext_net_id") // NB: even if ext_net_id value is explicitly set to 0, argSet = false anyway
	if argSet {
		if argVal.(int) > 0 {
			// connect to specific external network
			urlValues.Add("extNetId", fmt.Sprintf("%d", argVal.(int)))
			/* 
			 Commented out, as we've made "ext_net_ip" parameter non-configurable via Terraform!

			// in case of specific ext net connection user may also want a particular IP address
			argVal, argSet = d.GetOk("ext_net_ip")
			if argSet && argVal.(string) != "" {
				urlValues.Add("extIp", argVal.(string))
			}
			*/
		} else {
			// ext_net_id is set to a negative value - connect to default external network
			// no particular IP address selection in this case
			urlValues.Add("extNetId", "0")
		}
	}

	argVal, argSet = d.GetOk("ipcidr")
	if argSet && argVal.(string) != "" {
		log.Debugf("resourceVinsCreate: ipcidr is set to %s", argVal.(string))
		urlValues.Add("ipcidr", argVal.(string))
	}
	
	argVal, argSet = d.GetOk("description")
	if argSet {
		urlValues.Add("desc", argVal.(string))
	}

	apiResp, err, _ := controller.decortAPICall("POST", apiToCall, urlValues)
	if err != nil {
		return err
	}

	d.SetId(apiResp) // update ID of the resource to tell Terraform that the ViNS resource exists
	vinsId, _ := strconv.Atoi(apiResp)

	log.Debugf("resourceVinsCreate: new ViNS ID / name %d / %s creation sequence complete", vinsId, d.Get("name").(string))

	// We may reuse dataSourceVinsRead here as we maintain similarity 
	// between ViNS resource and ViNS data source schemas
	// ViNS resource read function will also update resource ID on success, so that Terraform 
	// will know the resource exists (however, we already did it a few lines before)
	return dataSourceVinsRead(d, m)
}

func resourceVinsRead(d *schema.ResourceData, m interface{}) error {
	vinsFacts, err := utilityVinsCheckPresence(d, m)
	if vinsFacts == "" {
		// if empty string is returned from utilityVinsCheckPresence then there is no
		// such ViNS and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty
		return err
	}

	return flattenVins(d, vinsFacts)
}

func resourceVinsUpdate(d *schema.ResourceData, m interface{}) error {

	log.Debugf("resourceVinsUpdate: called for ViNS ID / name %s / %s,  Account ID %d, RG ID %d",
		d.Id(), d.Get("name").(string), d.Get("account_id").(int),  d.Get("rg_id").(int))

	controller := m.(*ControllerCfg)

	d.Partial(true)
	
	// 1. Handle external network connection change
	oldExtNetId, newExtNedId := d.GetChange("ext_net_id")
	if oldExtNetId.(int) != newExtNedId.(int) {
		log.Debugf("resourceVinsUpdate: changing ViNS ID %s - ext_net_id %d -> %d", d.Id(), oldExtNetId.(int), newExtNedId.(int))

		extnetParams := &url.Values{}
		extnetParams.Add("vinsId", d.Id())

		if oldExtNetId.(int) > 0 {
			// there was preexisting external net connection - disconnect ViNS
			_, err, _ := controller.decortAPICall("POST", VinsExtNetDisconnectAPI, extnetParams)
			if err != nil {
				return err
			}
		}

		if newExtNedId.(int) > 0 {
			// new external network connection requested - connect ViNS
			extnetParams.Add("netId", fmt.Sprintf("%d", newExtNedId.(int)))
			_, err, _ := controller.decortAPICall("POST", VinsExtNetConnectAPI, extnetParams)
			if err != nil {
				return err
			}
		}
		
		d.SetPartial("ext_net_id")
	}

	d.Partial(false)

	// we may reuse dataSourceVinsRead here as we maintain similarity 
	// between Compute resource and Compute data source schemas
	return dataSourceVinsRead(d, m) 
}

func resourceVinsDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceVinsDelete: called for ViNS ID / name %s / %s, Account ID %d, RG ID %d",
		d.Id(), d.Get("name").(string), d.Get("account_id").(int), d.Get("rg_id").(int))

	vinsFacts, err := utilityVinsCheckPresence(d, m)
	if vinsFacts == "" {
		// the specified ViNS does not exist - in this case according to Terraform best practice
		// we exit from Destroy method without error
		return nil
	}

	params := &url.Values{}
	params.Add("vinsId", d.Id())
	params.Add("force", "1")        // disconnect all computes before deleting ViNS
	params.Add("permanently", "1")  // delete ViNS immediately bypassing recycle bin

	controller := m.(*ControllerCfg)
	_, err, _ = controller.decortAPICall("POST", VinsDeleteAPI, params)
	if err != nil {
		return err
	}

	return nil
}

func resourceVinsExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	log.Debugf("resourceVinsExists: called for ViNS name %s, Account ID %d, RG ID %d",
		d.Get("name").(string), d.Get("account_id").(int), d.Get("rg_id").(int))

	vinsFacts, err := utilityVinsCheckPresence(d, m)
	if vinsFacts == "" {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func resourceVinsSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description: "Name of the ViNS. Names are case sensitive and unique within the context of an account or resource group.",
		},

		/* we do not need ViNS ID as an argument because if we already know this ID, it is not practical to call resource provider.
		   Resource Import will work anyway, as it obtains the ID of ViNS to be imported through another mechanism.
		"vins_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Unique ID of the ViNS. If ViNS ID is specified, then ViNS name, rg_id and account_id are ignored.",
		},
		*/

		"rg_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "ID of the resource group, where this ViNS belongs to. Non-zero for ViNS created at resource group level, 0 otherwise.",
		},

		"account_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.IntAtLeast(1),
			Description: "ID of the account, which this ViNS belongs to. For ViNS created at account level, resource group ID is 0.",
		},

		"ext_net_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntAtLeast(0),
			Description: "ID of the external network this ViNS is connected to. Pass 0 if no external connection required.",
		},

		"ipcidr": {
			Type:        schema.TypeString,
			Optional:    true,
			DiffSuppressFunc: ipcidrDiffSupperss,
			Description: "Network address to use by this ViNS. This parameter is only valid when creating new ViNS.",
		},

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Optional user-defined text description of this ViNS.",
		},

		// the rest of attributes are computed
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the account, which this ViNS belongs to.",
		},

		"ext_ip_addr": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP address of the external connection (valid for ViNS connected to external network, ignored otherwise).",
		},
	}

	return rets
}

func resourceVins() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceVinsCreate,
		Read:   resourceVinsRead,
		Update: resourceVinsUpdate,
		Delete: resourceVinsDelete,
		Exists: resourceVinsExists,

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

		Schema: resourceVinsSchemaMake(),
	}
}

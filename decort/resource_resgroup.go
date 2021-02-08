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
	"log"
	"net/url"
	"strconv"
	"strings"
	
	"github.com/hashicorp/terraform/helper/schema"

)

func resourceResgroupCreate(d *schema.ResourceData, m interface{}) error {
	// First validate that we have all parameters required to create the new Resource Group
	arg_set := false
	account_name, arg_set := d.GetOk("account")
	if !arg_set {
		return  fmt.Errorf("Cannot create new RG: missing account.")
	}
	rg_name, arg_set := d.GetOk("name")
	if !arg_set {
		return  fmt.Errorf("Cannot create new RG: missing name.")
	}
	grid_id, arg_set := d.GetOk("grid_id")
	if !arg_set {
		return  fmt.Errorf("Cannot create new RG %q for account %q: missing Grid ID.", 
		                    rg_name.(string), account_name.(string))
	}

	// all required parameters are set in the schema - we can continue with RG creation
	log.Debugf("resourceResgroupCreate: called for RG name %q, account name %q", 
			   account_name.(string), rg_name.(string))
			   
	// Valid account ID is required to create new resource group
	// obtain Account ID by account name - it should not be zero on success
	validated_account_id, err := utilityGetAccountIdByName(account_name.(string), m)
	if err != nil {
		return err
	}

	// quota settings are optional
	set_quota := false
	var quota_record QuotaRecord 
	arg_value, arg_set = d.GetOk("quota")
	if arg_set {
		log.Debugf("resourceResgroupCreate: setting Quota on RG requested")
		quota_record, _ = makeQuotaRecord(arg_value.([]interface{}))
		set_quota = true
	}

	controller := m.(*ControllerCfg)
	log.Debugf("resourceResgroupCreate: called by user %q for RG name %q, account  %q / ID %d, Grid ID %d",
	            controller.getdecortUsername(),
				rg_name.(string), account_name.(string), validated_account_id, gird_id.(int))
	/*
	type ResgroupCreateParam struct {
	AccountID int          `json:"accountId"`
	GridId int             `json:"gid"`
	Name string            `json:"name"`
	Ram int                `json:"maxMemoryCapacity"`
	Disk int               `json:"maxVDiskCapacity"`
	Cpu int                `json:"maxCPUCapacity"`
	NetTraffic int         `json:"maxNetworkPeerTransfer"`
	ExtIPs int             `json:"maxNumPublicIP"`
	Owner string           `json:"owner"`
	DefNet string          `json:"def_net"`
	IPCidr string          `json:"ipcidr"`
	Desc string            `json:"decs"`
	Reason string          `json:"reason"`
	ExtNetID int           `json:"extNetId"`
	ExtIP string           `json:"extIp"`	
} 
	*/
				
	url_values := &url.Values{}
	url_values.Add("accountId", fmt.Sprintf("%d", validated_account_id))
	url_values.Add("name", rg_name.(string))
	url_values.Add("gid", fmt.Sprintf("%d", grid_id.(int)))
	url_values.Add("owner", controller.getdecortUsername())
	
	// pass quota values as set
	if set_quota {
		url_values.Add("maxCPUCapacity", fmt.Sprintf("%d", quota_record.Cpu))
		url_values.Add("maxVDiskCapacity", fmt.Sprintf("%d", quota_record.Disk))
		url_values.Add("maxMemoryCapacity", fmt.Sprintf("%d", quota_record.Ram))
		url_values.Add("maxNetworkPeerTransfer", fmt.Sprintf("%d", quota_record.ExtTraffic))
		url_values.Add("maxNumPublicIP", fmt.Sprintf("%d", quota_record.ExtIPs))
		// url_values.Add("???", fmt.Sprintf("%d", quota_record.GpuUnits))
	}

	// parse and handle network settings
	def_net_type, arg_set = d.GetOk("def_net_type")
	if arg_set {
		ulr_values.Add("def_net", def_net_type.(string))
	}

	ipcidr, arg_set = d.GetOk("ipcidr")
	if arg_set {
		ulr_values.Add("ipcidr", ipcidr.(string))
	}

	ext_net_id, arg_set = d.GetOk("ext_net_id")
	if arg_set {
		ulr_values.Add("extNetId", ext_net_id.(int))
	}

	ext_ip, arg_set = d.GetOk("ext_ip")
	if arg_set {
		ulr_values.Add("extIp", ext_ip.(string))
	}
	
	api_resp, err := controller.decortAPICall("POST", ResgroupCreateAPI, url_values)
	if err != nil {
		return err
	}

	d.SetId(api_resp) // rg/create API returns ID of the newly creted resource group on success
	rg.ID, _ = strconv.Atoi(api_resp)

	// re-read newly created RG to make sure schema contains complete and up to date set of specifications
	return resourceResgroupRead(d, m)
}

func resourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceResgroupRead: called for RG name %q, account name %q", 
	           d.Get("name").(string), d.Get("account").(string))
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
	log.Debugf("resourceResgroupUpdate: called for RG name %q, account name %q", 
			   d.Get("name").(string), d.Get("account").(string))

	do_update := false

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	url_values.Add("rgId", d.Id())

	name_new, name_set := d.GetOk("name")
	if name_set {
		log.Debugf("resourceResgroupUpdate: name specified - looking for deltas from the old settings.")
		name_old, _ := d.GetChange("name")
		if name_old.(string) != name_new.(string) {
			do_update := true
			url_values.Add("name", name_new.(string))
		}
	}
	
	quota_value, quota_set := d.GetOk("quota")
	if quota_set {
		log.Debugf("resourceResgroupUpdate: quota specified - looking for deltas from the old quota.")
		quotarecord_new, _ := makeQuotaRecord(quota_value.([]interface{}))
		quota_value_old, _ = d.GetChange("quota") // returns old as 1st, new as 2nd return value
		quotarecord_old, _ := makeQuotaRecord(quota_value_old.([]interface{}))

		if quotarecord_new.Cpu != quotarecord_old.Cpu {
			do_update = true
			log.Debugf("resourceResgroupUpdate: Cpu diff %d <- %d", quotarecord_new.Cpu, quotarecord_old.Cpu)
			url_values.Add("maxCPUCapacity", fmt.Sprintf("%d", quotarecord_new.Cpu))
		}
	
		if quotarecord_new.Disk != quotarecord_old.Disk {
			do_update = true
			log.Debugf("resourceResgroupUpdate: Disk diff %d <- %d", quotarecord_new.Disk, quotarecord_old.Disk)
			url_values.Add("maxVDiskCapacity", fmt.Sprintf("%d", quotarecord_new.Disk))
		}
	
		if quotarecord_new.Ram != quotarecord_old.Ram {
			do_update = true
			log.Debugf("resourceResgroupUpdate: Ram diff %f <- %f", quotarecord_new.Ram, quotarecord_old.Ram)
			url_values.Add("maxMemoryCapacity", fmt.Sprintf("%f", quotarecord_new.Ram))
		}
	
		if quotarecord_new.ExtTraffic != quotarecord_old.ExtTraffic {
			do_update = true
			log.Debugf("resourceResgroupUpdate: NetTraffic diff %d <- %d", quotarecord_new.ExtTraffic, quotarecord_old.ExtTraffic)
			url_values.Add("maxNetworkPeerTransfer", fmt.Sprintf("%d", quotarecord_new.NetTraffic))
		}
	
		if quotarecord_new.ExtIPs != quotarecord_old.ExtIPs {
			do_update = true
			log.Debugf("resourceResgroupUpdate: ExtIPs diff %d <- %d", quotarecord_new.ExtIPs, quotarecord_old.ExtIPs)
			url_values.Add("maxNumPublicIP", fmt.Sprintf("%d", quotarecord_new.ExtIPs))
		}
	}

	desc_new, desc_set := d.GetOk("desc")
	if desc_set {
		log.Debugf("resourceResgroupUpdate: description specified - looking for deltas from the old settings.")
		desc_old, _ := d.GetChange("desc")
		if desc_old.(string) != desc_new.(string) {
			do_update := true
			url_values.Add("desc", desc_new.(string))
		}
	}

	if do_update {
		log.Debugf("resourceResgroupUpdate: detected delta between new and old RG specs - updating the RG")
		_, err := controller.decortAPICall("POST", ResgroupUpdateAPI, url_values)
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
	log.Debugf("resourceResgroupDelete: called for RG name %q, account name %q", 
			   d.Get("name").(string), d.Get("account").(string))

	rg_facts, err := utilityResgroupCheckPresence(d, m)
	if rg_facts == "" {
		// the target RG does not exist - in this case according to Terraform best practice 
		// we exit from Destroy method without error
		return nil
	}

	url_values := &url.Values{}
	url_values.Add("rgId", d.Id())
	url_values.Add("force", "true")
	url_values.Add("permanently", "true")
	url_values.Add("reason", "Destroyed by DECORT Terraform provider")

	controller := m.(*ControllerCfg)
	_, err = controller.decortAPICall("POST", ResgroupDeleteAPI, url_values)
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
	return &schema.Resource {
		SchemaVersion: 1,

		Create: resourceResgroupCreate,
		Read:   resourceResgroupRead,
		Update: resourceResgroupUpdate,
		Delete: resourceResgroupDelete,
		Exists: resourceResgroupExists,

		Timeouts: &schema.ResourceTimeout {
			Create:  &Timeout180s,
			Read:    &Timeout30s,
			Update:  &Timeout180s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema {
			"name": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this resource group. Names are case sensitive and unique within the context of a account.",
			},

			"account": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the account, which this resource group belongs to.",
			},

			"def_net": &schema.Schema {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "PRIVATE"
				Description: "Type of the network, which this resource group will use as default for its computes - PRIVATE or PUBLIC or NONE.",
			},

			"ipcidr": &schema.Schema {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the netowrk inside the private network segment (aka ViNS) if def_net=PRIVATE",
			},

			"ext_net_id": &schema.Schema {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the external network, which this resource group will use as default for its computes if def_net=PUBLIC",
			},

			"ext_ip": &schema.Schema {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address on the external netowrk to request, if def_net=PUBLIC",
			},

			"account_id": &schema.Schema {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the account, which this resource group belongs to.",
			},

			"grid_id": &schema.Schema {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique ID of the grid, where this resource group is deployed.",
			},

			"quota": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        &schema.Resource {
					Schema:  quotasSubresourceSchema(),
				},
				Description: "Quota settings for this resource group.",
			},

			"desc": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "User-defined text description of this resource group."
			},

			"status": { 
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Current status of this resource group.",
			},

			"def_net_id": &schema.Schema {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the default network for this resource group (if any).",
			},

			"vins": {
				Type:          schema.TypeList,
				Computed:      true,
				MaxItems:      LimitMaxVinsPerResgroup,
				Elem:          &schema.Resource {
					Schema: vinsRgSubresourceSchema() // this is a list of ints
				},
				Description: "List of VINs deployed in this resource group.",
			},

			"computes": {
				Type:          schema.TypeList,
				Computed:      true,
				Elem:          &schema.Resource {
					Schema: computesRgSubresourceSchema() //this is a list of ints
				},
				Description: "List of computes deployed in this resource group."
			},
		},
	}
}
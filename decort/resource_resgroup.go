/*
Copyright (c) 2019-2020 Digital Energy Cloud Solutions. All Rights Reserved.

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
	
	"github.com/hashicorp/terraform/helper/schema"

)

func resourceResgroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceResgroupCreate: called for res group name %q, account name %q", 
			   d.Get("name").(string), d.Get("account").(string))
			   
	rg := &ResgroupConfig{
		Name:         d.Get("name").(string),
		AccountName:  d.Get("account").(string),
	}

	// validate that we have all parameters required to create the new Resource Group
	// location code is required to create new resource group
	arg_value, arg_set := d.GetOk("location")
	if arg_set {
		rg.Location = arg_value.(string)
	} else {
		return  fmt.Errorf("Cannot create new RG %q for account %q: missing location parameter.", 
		                    rg.Name, rg.AccountName)
	}
	// account ID is required to create new resource group
	// obtain Account ID by account name - it should not be zero on success
	account_id, err := utilityGetAccountIdByName(rg.AccountName, m)
	if err != nil {
		return err
	}
	rg.AccountID = account_id

	set_quotas := false
	arg_value, arg_set = d.GetOk("quotas")
	if arg_set {
		log.Debugf("resourceResgroupCreate: calling makeQuotaConfig")
		rg.Quota, _ = makeQuotaConfig(arg_value.([]interface{}))
		set_quotas = true
	}

	controller := m.(*ControllerCfg)
	log.Debugf("resourceResgroupCreate: called by user %q for RG name %q, for account  %q / ID %d, location %q",
	            controller.getdecortUsername(),
				rg.Name, d.Get("account").(string), rg.AccountID, rg.Location)
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
	url_values.Add("accountId", fmt.Sprintf("%d", rg.AccountID))
	url_values.Add("name", rg.Name)
	url_values.Add("gid", rg.Location)
	url_values.Add("owner", controller.getdecortUsername())
	url_values.Add("def_net", "NONE")
	// pass quota values as set
	if set_quotas {
		url_values.Add("maxCPUCapacity", fmt.Sprintf("%d", rg.Quota.Cpu))
		url_values.Add("maxVDiskCapacity", fmt.Sprintf("%d", rg.Quota.Disk))
		url_values.Add("maxMemoryCapacity", fmt.Sprintf("%f", rg.Quota.Ram))
		url_values.Add("maxNetworkPeerTransfer", fmt.Sprintf("%d", rg.Quota.NetTraffic))
		url_values.Add("maxNumPublicIP", fmt.Sprintf("%d", rg.Quota.ExtIPs))
	}
	// pass externalnetworkid if set
	arg_value, arg_set = d.GetOk("extnet_id")
	if arg_set {
		url_values.Add("extNetId", fmt.Sprintf("%d", arg_value))
	}
	
	api_resp, err := controller.decortAPICall("POST", ResgroupCreateAPI, url_values)
	if err != nil {
		return err
	}

	d.SetId(api_resp) // rg/create API returns ID of the newly creted resource group on success
	rg.ID, _ = strconv.Atoi(api_resp)

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
	// this method will only update quotas, if any are set
	log.Debugf("resourceResgroupUpdate: called for RG name %q, account name %q", 
			   d.Get("name").(string), d.Get("account").(string))

	quota_value, arg_set := d.GetOk("quotas")
	if !arg_set {
		// if there are no quotas set explicitly in the resource configuration - no change will be done
		log.Debugf("resourceResgroupUpdate: quotas are not set in the resource config - no update on this resource will be done")
		return resourceResgroupRead(d, m)
	}
	quotaconfig_new, _ := makeQuotaConfig(quota_value.([]interface{}))

	quota_value, _ = d.GetChange("quotas") // returns old as 1st, new as 2nd argument
	quotaconfig_old, _ := makeQuotaConfig(quota_value.([]interface{}))

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	url_values.Add("cloudspaceId", d.Id())
	url_values.Add("name", d.Get("name").(string))
	
	do_update := false

	if quotaconfig_new.Cpu != quotaconfig_old.Cpu {
		do_update = true
		log.Debugf("resourceResgroupUpdate: Cpu diff %d <- %d", quotaconfig_new.Cpu, quotaconfig_old.Cpu)
		url_values.Add("maxCPUCapacity", fmt.Sprintf("%d", quotaconfig_new.Cpu))
	}

	if quotaconfig_new.Disk != quotaconfig_old.Disk {
		do_update = true
		log.Debugf("resourceResgroupUpdate: Disk diff %d <- %d", quotaconfig_new.Disk, quotaconfig_old.Disk)
		url_values.Add("maxVDiskCapacity", fmt.Sprintf("%d", quotaconfig_new.Disk))
	}

	if quotaconfig_new.Ram != quotaconfig_old.Ram {
		do_update = true
		log.Debugf("resourceResgroupUpdate: Ram diff %f <- %f", quotaconfig_new.Ram, quotaconfig_old.Ram)
		url_values.Add("maxMemoryCapacity", fmt.Sprintf("%f", quotaconfig_new.Ram))
	}

	if quotaconfig_new.NetTraffic != quotaconfig_old.NetTraffic {
		do_update = true
		log.Debugf("resourceResgroupUpdate: NetTraffic diff %d <- %d", quotaconfig_new.NetTraffic, quotaconfig_old.NetTraffic)
		url_values.Add("maxNetworkPeerTransfer", fmt.Sprintf("%d", quotaconfig_new.NetTraffic))
	}

	if quotaconfig_new.ExtIPs != quotaconfig_old.ExtIPs {
		do_update = true
		log.Debugf("resourceResgroupUpdate: ExtIPs diff %d <- %d", quotaconfig_new.ExtIPs, quotaconfig_old.ExtIPs)
		url_values.Add("maxNumPublicIP", fmt.Sprintf("%d", quotaconfig_new.ExtIPs))
	}

	if do_update {
		log.Debugf("resourceResgroupUpdate: some new quotas are set - updating the resource")
		_, err := controller.decortAPICall("POST", ResgroupUpdateAPI, url_values)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("resourceResgroupUpdate: no difference in quotas between old and new state - no update on this resource will be done")
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

	params := &url.Values{}
	params.Add("rgId", d.Id())
	params.Add("force", "true")
	params.Add("permanently", "true")
	params.Add("reason", "Destroyed by DECORT Terraform provider")

	controller := m.(*ControllerCfg)
	vm_facts, err = controller.decortAPICall("POST", ResgroupDeleteAPI, params)
	if err != nil {
		return err
	}

	return nil
}

func resourceResgroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify ResourceData argument
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

			"extnet_id": &schema.Schema {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the external network, which this resource group will be connected to by default.",
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

			"quotas": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        &schema.Resource {
					Schema:  quotasSubresourceSchema(),
				},
				Description: "Quotas on the resources for this resource group.",
			},
		},
	}
}
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

package decort

import (

	// "encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus" 
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourcePfwCreate(d *schema.ResourceData, m interface{}) error {
	compId := d.Get("compute_id")

	rules_set, ok := d.GetOk("rules")
	if !ok || rules_set.(*schema.Set).Len() == 0 {
		log.Debugf("resourcePfwCreate: empty new PFW rules set requested for compute ID %d - nothing to create", compId.(int))
		return nil
	}

	log.Debugf("resourcePfwCreate: ready to setup %d PFW rules for compute ID %d", 
	           rules_set.(*schema.Set).Len(), compId.(int))

	controller := m.(*ControllerCfg)
	apiErrCount := 0
	var lastSavedError error

	for _, runner := range rules_set.(*schema.Set).List() {
		rule := runner.(map[string]interface{})
		params := &url.Values{}
		params.Add("computeId", fmt.Sprintf("%d", compId.(int)))
		params.Add("publicPortStart", fmt.Sprintf("%d", rule["pub_port_start"].(int)))
		params.Add("publicPortEnd", fmt.Sprintf("%d", rule["pub_port_end"].(int)))
		params.Add("localBasePort", fmt.Sprintf("%d", rule["local_port"].(int)))
		params.Add("proto", rule["proto"].(string))
		log.Debugf("resourcePfwCreate: ready to add rule %d:%d -> %d %s for Compute ID %d",
	               rule["pub_port_start"].(int),rule["pub_port_end"].(int),
				   rule["local_port"].(int), rule["proto"].(string),
				   compId.(int))
		_, err, _ := controller.decortAPICall("POST", ComputePfwAddAPI, params)
		if err != nil {
			log.Errorf("resourcePfwCreate: error adding rule %d:%d -> %d %s for Compute ID %d: %s",
	               rule["pub_port_start"].(int),rule["pub_port_end"].(int),
				   rule["local_port"].(int), rule["proto"].(string),
				   compId.(int),
				   err)
			apiErrCount++
			lastSavedError = err
		}
	}
	
	if apiErrCount > 0 {
		log.Errorf("resourcePfwCreate: there were %d error(s) adding PFW rules to Compute ID %s. Last error was: %s", 
				   apiErrCount, compId.(int), lastSavedError)
		return lastSavedError
	}

	return nil
}

func resourcePfwRead(d *schema.ResourceData, m interface{}) error {
	pfwFacts, err := utilityPfwCheckPresence(d, m)
	if pfwFacts == "" {
		// if empty string is returned from dataSourcePfwRead then we got no
		// PFW rules. It could also be because there was some error, which
		// is indicated by non-nil err value
		d.SetId("") // ensure ID is empty in this case anyway
		return err
	}

	return flattenPfw(d, pfwFacts)
}

func resourcePfwUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO: update not implemented yet
	compId := d.Get("compute_id")
	return fmt.Errorf("resourcePfwUpdate: method is not implemented yet (Compute ID %d)", compId.(int))
}

func resourcePfwDelete(d *schema.ResourceData, m interface{}) error {
	compId := d.Get("compute_id")

	rules_set, ok := d.GetOk("rules")
	if !ok || rules_set.(*schema.Set).Len() == 0 {
		log.Debugf("resourcePfwCreate: no PFW rules defined for compute ID %d - nothing to delete", compId.(int))
		return nil
	}

	log.Debugf("resourcePfwDelete: ready to delete %d PFW rules from compute ID %d", 
	           rules_set.(*schema.Set).Len(), compId.(int))

	controller := m.(*ControllerCfg)
	apiErrCount := 0
	var lastSavedError error

	for _, runner := range rules_set.(*schema.Set).List() {
		rule := runner.(map[string]interface{})
		params := &url.Values{}
		params.Add("computeId", fmt.Sprintf("%d", compId.(int)))
		params.Add("ruleId", fmt.Sprintf("%d", rule["id"].(int)))
		log.Debugf("resourcePfwCreate: ready to delete rule ID%s (%d:%d -> %d %s) from Compute ID %d",
		           rule["id"].(int),
	               rule["pub_port_start"].(int),rule["pub_port_end"].(int),
				   rule["local_port"].(int), rule["proto"].(string),
				   compId.(int))
		_, err, _ := controller.decortAPICall("POST", ComputePfwDelAPI, params)
		if err != nil {
			log.Errorf("resourcePfwDelete: error deleting rule ID %d (%d:%d -> %d %s) from Compute ID %d: %s",
			       rule["id"].(int),
	               rule["pub_port_start"].(int),rule["pub_port_end"].(int),
				   rule["local_port"].(int), rule["proto"].(string),
				   compId.(int),
				   err)
			apiErrCount++
			lastSavedError = err
		}
	}
	
	if apiErrCount > 0 {
		log.Errorf("resourcePfwDelete: there were %d error(s) when deleting PFW rules from Compute ID %s. Last error was: %s", 
				   apiErrCount, compId.(int), lastSavedError)
		return lastSavedError
	}

	return nil
}

func resourcePfwExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	log.Debugf("resourcePfwExists: called for Compute ID %d", d.Get("compute_id").(int))

	pfwFacts, err := utilityPfwCheckPresence(d, m)
	if pfwFacts == "" {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func resourcePfw() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourcePfwCreate,
		Read:   resourcePfwRead,
		Update: resourcePfwUpdate,
		Delete: resourcePfwDelete,
		Exists: resourcePfwExists,

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
			"compute_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the compute instance to configure port forwarding rules for.",
			},

			"vins_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the ViNS to configure port forwarding rules on. Compute must be already plugged into this ViNS and ViNS must have external network connection.",
			},

			"rule": {
				Type:         schema.TypeSet,
				Optional:     true,
				Elem:         &schema.Resource{
					Schema: rulesSubresourceSchemaMake(),
				},
				Description:  "Port forwarding rule. You may specify several rules, one in each such block.",
			},
		},
	}
}

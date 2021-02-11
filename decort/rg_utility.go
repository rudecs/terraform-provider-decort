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
	"log"
	"net/url"

	// "strconv"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"
)

func (ctrl *ControllerCfg) utilityResgroupConfigGet(rgid int) (*ResgroupGetResp, error) {
	url_values := &url.Values{}
	url_values.Add("rgId", fmt.Sprintf("%d", rgid))
	resgroup_facts, err := ctrl.decortAPICall("POST", ResgroupGetAPI, url_values)
	if err != nil {
		return nil, err
	}

	log.Debugf("utilityResgroupConfigGet: ready to unmarshal string %q", resgroup_facts)
	model := &ResgroupGetResp{}
	err = json.Unmarshal([]byte(resgroup_facts), model)
	if err != nil {
		return nil, err
	}

	/*
		ret := &ResgroupConfig{}
		ret.AccountID = model.AccountID
		ret.Location = model.Location
		ret.Name = model.Name
		ret.ID = rgid
		ret.GridID = model.GridID
		ret.ExtIP = model.ExtIP   // legacy field for VDC - this will eventually become obsoleted by true Resource Groups
		// Quota ResgroupQuotaConfig
		// Network NetworkConfig
	*/
	log.Debugf("utilityResgroupConfigGet: account ID %d, GridID %d, Name %s",
		model.AccountID, model.GridID, model.Name)

	return model, nil
}

// On success this function returns a string, as returned by API rg/get, which could be unmarshalled
// into ResgroupGetResp structure
func utilityResgroupCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	// This function tries to locate resource group by one of the following algorithms depending
	// on the parameters passed:
	//    - if resource group ID is specified -> by RG ID
	//    - if resource group name is specifeid -> by RG name and either account ID or account name
	//
	// If succeeded, it returns non empty string that contains JSON formatted facts about the
	// resource group as returned by cloudspaces/get API call.
	// Otherwise it returns empty string and meaningful error.
	//
	// NOTE: As our provider always deletes RGs permanently, there is no "restore" method and
	// consequently we are not interested in matching RGs in DELETED state. Hence, we call
	// .../rg/list API with includedeleted=false
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for the Terraform resource Exists method.
	//

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}

	rg_id, arg_set := d.GetOk("rg_id")
	if arg_set {
		// go straight for the RG by its ID
		log.Debugf("utilityResgroupCheckPresence: locating RG by its ID %d", rg_id.(int))
		url_values.Add("rgId", fmt.Sprintf("%d", rg_id.(int)))
		rg_facts, err := controller.decortAPICall("POST", ResgroupGetAPI, url_values)
		if err != nil {
			return "", err
		}
		return rg_facts, nil
	}

	rg_name, arg_set := d.GetOk("name")
	if !arg_set {
		// no RG ID and no RG name - we cannot locate resource group in this case
		return "", fmt.Error("Cannot check resource group presence if name is empty and no resource group ID specified.")
	}

	// Valid account ID is required to locate a resource group
	// obtain Account ID by account name - it should not be zero on success
	validated_account_id, err := utilityGetAccountIdBySchema(d, m)
	if err != nil {
		return err
	}

	url_values.Add("includedeleted", "false")
	body_string, err := controller.decortAPICall("POST", ResgroupListAPI, url_values)
	if err != nil {
		return "", err
	}
	// log.Debugf("%s", body_string)
	// log.Debugf("utilityResgroupCheckPresence: ready to decode response body from %q", ResgroupListAPI)
	model := ResgroupListResp{}
	err = json.Unmarshal([]byte(body_string), &model)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityResgroupCheckPresence: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// match by RG name & account ID
		if item.Name == rg_name.(string) && item.AccountID == validated_account_id {
			log.Debugf("utilityResgroupCheckPresence: match RG name %q / ID %d, account ID %d at index %d",
				item.Name, item.ID, item.AccountID, index)

			// not all required information is returned by rg/list API, so we need to initiate one more
			// call to rg/get to obtain extra data to complete Resource population.
			// Namely, we need resource quota settings
			req_values := &url.Values{}
			req_values.Add("rgId", fmt.Sprintf("%d", item.ID))
			body_string, err := controller.decortAPICall("POST", ResgroupGetAPI, req_values)
			if err != nil {
				return "", err
			}

			return body_string, nil
		}
	}

	return "", fmt.Errorf("Cannot find RG name %q owned by account ID %d", name, validated_account_id)
}

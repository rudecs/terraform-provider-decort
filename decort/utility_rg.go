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
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func (ctrl *ControllerCfg) utilityResgroupConfigGet(rgid int) (*ResgroupGetResp, error) {
	urlValues := &url.Values{}
	urlValues.Add("rgId", fmt.Sprintf("%d", rgid))
	rgFacts, err := ctrl.decortAPICall("POST", ResgroupGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	log.Debugf("utilityResgroupConfigGet: ready to unmarshal string %q", rgFacts)
	model := &ResgroupGetResp{}
	err = json.Unmarshal([]byte(rgFacts), model)
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
	// resource group as returned by rg/get API call.
	// Otherwise it returns empty string and a meaningful error.
	//
	// NOTE: As our provider always deletes RGs permanently, there is no "restore" method and
	// consequently we are not interested in matching RGs in DELETED state. Hence, we call
	// .../rg/list API with includedeleted=false
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for the Terraform resource Exists method.
	//

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	rgId, argSet := d.GetOk("rg_id")
	if argSet {
		// go straight for the RG by its ID
		log.Debugf("utilityResgroupCheckPresence: locating RG by its ID %d", rgId.(int))
		urlValues.Add("rgId", fmt.Sprintf("%d", rgId.(int)))
		rgFacts, err := controller.decortAPICall("POST", ResgroupGetAPI, urlValues)
		if err != nil {
			return "", err
		}
		return rgFacts, nil
	}

	rgName, argSet := d.GetOk("name")
	if !argSet {
		// no RG ID and no RG name - we cannot locate resource group in this case
		return "", fmt.Errorf("Cannot check resource group presence if name is empty and no resource group ID specified")
	}

	// Valid account ID is required to locate a resource group
	// obtain Account ID by account name - it should not be zero on success
	validatedAccountId, err := utilityGetAccountIdBySchema(d, m)
	if err != nil {
		return "", err
	}

	urlValues.Add("includedeleted", "false")
	apiResp, err := controller.decortAPICall("POST", ResgroupListAPI, urlValues)
	if err != nil {
		return "", err
	}
	// log.Debugf("%s", apiResp)
	log.Debugf("utilityResgroupCheckPresence: ready to decode response body from %s", ResgroupListAPI)
	model := ResgroupListResp{}
	err = json.Unmarshal([]byte(apiResp), &model)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityResgroupCheckPresence: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// match by RG name & account ID
		if item.Name == rgName.(string) && item.AccountID == validatedAccountId {
			log.Debugf("utilityResgroupCheckPresence: match RG name %q / ID %d, account ID %d at index %d",
				item.Name, item.ID, item.AccountID, index)

			// not all required information is returned by rg/list API, so we need to initiate one more
			// call to rg/get to obtain extra data to complete Resource population.
			// Namely, we need resource quota settings
			reqValues := &url.Values{}
			reqValues.Add("rgId", fmt.Sprintf("%d", item.ID))
			apiResp, err := controller.decortAPICall("POST", ResgroupGetAPI, reqValues)
			if err != nil {
				return "", err
			}

			return apiResp, nil
		}
	}

	return "", fmt.Errorf("Cannot find RG name %q owned by account ID %d", rgName, validatedAccountId)
}

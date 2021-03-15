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

func utilityAccountCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	accId, argSet := d.GetOk("account_id")
	if argSet {
		// get Account right away by its ID
		log.Debugf("utilityAccountCheckPresence: locating Account by its ID %d", accId.(int))
		urlValues.Add("accountId", fmt.Sprintf("%d", accId.(int)))
		apiResp, err := controller.decortAPICall("POST", AccountsGetAPI, urlValues)
		if err != nil {
			return "", err
		}
		return apiResp, nil
	}

	accName, argSet := d.GetOk("name")
	if !argSet {
		// neither ID nor name - no account for you!
		return "", fmt.Errorf("Cannot check account presence if name is empty and no account ID specified")
	}

	apiResp, err := controller.decortAPICall("POST", AccountsListAPI, urlValues)
	if err != nil {
		return "", err
	}
	// log.Debugf("%s", apiResp)
	// log.Debugf("utilityAccountCheckPresence: ready to decode response body from %q", AccountsListAPI)
	accList := AccountsListResp{}
	err = json.Unmarshal([]byte(apiResp), &accList)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityAccountCheckPresence: traversing decoded Json of length %d", len(accList))
	for index, item := range accList {
		// match by account name
		if item.Name == accName.(string) {
			log.Debugf("utilityAccountCheckPresence: match account name %q / ID %d at index %d",
				item.Name, item.ID, index)

			// NB: unlike accounts/get API, accounts/list API returns abridged set of account info,
			// for instance it does not return quotas

			reencodedItem, err := json.Marshal(item)
			if err != nil {
				return "", err
			}
			return string(reencodedItem[:]), nil
		}
	}

	return "", fmt.Errorf("Cannot find account name %q", accName.(string))
}

func utilityGetAccountIdBySchema(d *schema.ResourceData, m interface{}) (int, error) {
	/*
		This function expects schema that contains the following two elements:

		"account_name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    Optional,
			Description: "Name of the account, ....",
		},

		"account_id": &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Unique ID of the account, ....",
		},

		Then it will check, which argument is set, and if account name is present, it will
		initiate API calls to the DECORT cloud controller and try to match relevant account
		by the name.

	*/

	accId, argSet := d.GetOk("account_id")
	if argSet {
		if accId.(int) > 0 {
			return accId.(int), nil
		}
		return 0, fmt.Errorf("Account ID must be positive")
	}

	accName, argSet := d.GetOk("account_name")
	if !argSet {
		return 0, fmt.Errorf("Either non-empty account name or valid account ID must be specified")
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	apiResp, err := controller.decortAPICall("POST", AccountsListAPI, urlValues)
	if err != nil {
		return 0, err
	}

	model := AccountsListResp{}
	err = json.Unmarshal([]byte(apiResp), &model)
	if err != nil {
		return 0, err
	}

	log.Debugf("utilityGetAccountIdBySchema: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// need to match Account by name
		if item.Name == accName.(string) {
			log.Debugf("utilityGetAccountIdBySchema: match Account name %q / ID %d at index %d",
				item.Name, item.ID, index)
			return item.ID, nil
		}
	}

	return 0, fmt.Errorf("Cannot find account %q for the current user. Check account name and your access rights", accName.(string))
}

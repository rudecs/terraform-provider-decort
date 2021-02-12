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

func utilityAccountCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	controller := m.(*ControllerCfg)
	url_values := &url.Values{}

	acc_id, arg_set := d.GetOk("account_id")
	if arg_set {
		// get Account right away by its ID
		log.Debugf("utilityAccountCheckPresence: locating Account by its ID %d", acc_id.(int))
		url_values.Add("accountId", fmt.Sprintf("%d", acc_id.(int)))
		api_resp, err := controller.decortAPICall("POST", AccountsGetAPI, url_values)
		if err != nil {
			return "", err
		}
		return api_resp, nil
	}

	acc_name, arg_set := d.GetOk("name")
	if !arg_set {
		// neither ID nor name - no account for you!
		return "", fmt.Error("Cannot check account presence if name is empty and no account ID specified.")
	}

	api_resp, err := controller.decortAPICall("POST", AccountsListAPI, url_values)
	if err != nil {
		return "", err
	}
	// log.Debugf("%s", api_resp)
	// log.Debugf("utilityAccountCheckPresence: ready to decode response body from %q", AccountsListAPI)
	acc_list := AccountsListResp{}
	err = json.Unmarshal([]byte(api_resp), &acc_list)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityAccountCheckPresence: traversing decoded Json of length %d", len(model))
	for index, item := range acc_list {
		// match by account name
		if item.Name == acc_name.(string) {
			log.Debugf("utilityAccountCheckPresence: match account name %q / ID %d at index %d",
				item.Name, item.ID, index)

			// NB: unlike accounts/get API, accounts/list API returns abridged set of account info,
			// for instance it does not return quotas

			reencoded_item, err := json.Marshal(item)
			if err != nil {
				return "", err
			}
			return reencoded_item.(string), nil
		}
	}

	return "", fmt.Errorf("Cannot find account name %q owned by account ID %d", name, validated_account_id)
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

	account_id, arg_set := d.GetOk("account_id")
	if arg_set {
		if account_id.(int) > 0 {
			return account_id.(int), nil
		}
		return 0, fmt.Error("Account ID must be positive.")
	}

	account_name, arg_set := d.GetOk("account_name")
	if !arg_set {
		return 0, fmt.Error("Non-empty account name or positive account ID must be specified.")
	}

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	body_string, err := controller.decortAPICall("POST", AccountsListAPI, url_values)
	if err != nil {
		return 0, err
	}

	model := AccountsListResp{}
	err = json.Unmarshal([]byte(body_string), &model)
	if err != nil {
		return 0, err
	}

	log.Debugf("utilityGetAccountIdBySchema: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// need to match Account by name
		if item.Name == account_name.(string) {
			log.Debugf("utilityGetAccountIdBySchema: match Account name %q / ID %d at index %d",
				item.Name, item.ID, index)
			return item.ID, nil
		}
	}

	return 0, fmt.Errorf("Cannot find account %q for the current user. Check account name and your access rights", account_name.(string))
}

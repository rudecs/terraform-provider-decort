/*
Copyright (c) 2019-2020 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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

func (ctrl *ControllerCfg) utilityResgroupConfigGet(rgid int) (*ResgroupConfig, error) {
	url_values := &url.Values{}
	url_values.Add("cloudspaceId", fmt.Sprintf("%d", rgid))
	resgroup_facts, err := ctrl.decortAPICall("POST", CloudspacesGetAPI, url_values)
	if err != nil {
		return nil, err
	}

	log.Printf("utilityResgroupConfigGet: ready to unmarshal string %q", resgroup_facts)
	model := CloudspacesGetResp{}
	err = json.Unmarshal([]byte(resgroup_facts), &model)
	if err != nil {
		return nil, err
	}

	ret := &ResgroupConfig{}
	ret.TenantID = model.TenantID
	ret.Location = model.Location
	ret.Name = model.Name
	ret.ID = rgid
	ret.GridID = model.GridID
	ret.ExtIP = model.ExtIP   // legacy field for VDC - this will eventually become obsoleted by true Resource Groups
	// Quota ResgroupQuotaConfig
	// Network NetworkConfig
	log.Printf("utilityResgroupConfigGet: tenant ID %d, GridID %d, ExtIP %q", 
	           model.TenantID, model.GridID, model.ExtIP)

	return ret, nil
}

func utilityResgroupCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	// This function tries to locate resource group by its name and tenant name.
	// If succeeded, it returns non empty string that contains JSON formatted facts about the 
	// resource group as returned by cloudspaces/get API call.
	// Otherwise it returns empty string and meaningful error.
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for the resource's Exists method.
	//
	name := d.Get("name").(string)
	tenant_name := d.Get("tenant").(string)

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	url_values.Add("includedeleted", "false")
	body_string, err := controller.decortAPICall("POST", CloudspacesListAPI, url_values)
	if err != nil {
		return "", err
	}

	log.Printf("%s", body_string)
	log.Printf("utilityResgroupCheckPresence: ready to decode response body from %q", CloudspacesListAPI)
	model := CloudspacesListResp{}
	err = json.Unmarshal([]byte(body_string), &model)
	if err != nil {
		return "", err
	}

	log.Printf("utilityResgroupCheckPresence: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// need to match VDC by name & tenant name
		if item.Name == name && item.TenantName == tenant_name {
			log.Printf("utilityResgroupCheckPresence: match ResGroup name %q / ID %d, tenant %q at index %d", 
					   item.Name, item.ID, item.TenantName, index)

			// not all required information is returned by cloudspaces/list API, so we need to initiate one more
			// call to cloudspaces/get to obtain extra data to complete Resource population.
			// Namely, we need to extract resource quota settings
			req_values := &url.Values{} 
			req_values.Add("cloudspaceId", fmt.Sprintf("%d", item.ID))
			body_string, err := controller.decortAPICall("POST", CloudspacesGetAPI, req_values)
			if err != nil {
				return "", err
			}

			return body_string, nil
		}
	}

	return "", fmt.Errorf("Cannot find resource group name %q owned by tenant %q", name, tenant_name)
}

func utilityGetTenantIdByName(tenant_name string, m interface{}) (int, error) {
	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	body_string, err := controller.decortAPICall("POST", TenantsListAPI, url_values)
	if err != nil {
		return 0, err
	}

	model := TenantsListResp{}
	err = json.Unmarshal([]byte(body_string), &model)
	if err != nil {
		return 0, err
	}

	log.Printf("utilityGetTenantIdByName: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// need to match Tenant by name
		if item.Name == tenant_name {
			log.Printf("utilityGetTenantIdByName: match Tenant name %q / ID %d at index %d", 
					   item.Name, item.ID, index)
			return item.ID, nil
		}
	}

	return 0, fmt.Errorf("Cannot find tenant %q for the current user. Check tenant value and your access rights", tenant_name)
}
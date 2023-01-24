/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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
Terraform DECORT provider - manage resources provided by DECORT (Digital Energy Cloud
Orchestration Technology) with Terraform by Hashicorp.

Source code: https://github.com/rudecs/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://github.com/rudecs/terraform-provider-decort/wiki
*/

package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/controller"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// On success this function returns a string, as returned by API rg/get, which could be unmarshalled
// into ResgroupGetResp structure
func utilityResgroupCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*ResgroupGetResp, error) {
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

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if d.Id() != "" {
		urlValues.Add("rgId", d.Id())
	} else {
		urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	}

	rgData := &ResgroupGetResp{}
	rgRaw, err := c.DecortAPICall(ctx, "POST", ResgroupGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(rgRaw), rgData)
	if err != nil {
		return nil, err
	}
	return rgData, nil
}

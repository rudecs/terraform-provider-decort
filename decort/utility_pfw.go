/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Petr Krutov, <petr.krutov@digitalenergy.online>

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
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func utilityPfwCheckPresence(d *schema.ResourceData, m interface{}) (*PfwRecord, error) {
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	resp, err := controller.decortAPICall("POST", ComputePfwListAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	idS := d.Id()
	id, err := strconv.Atoi(idS)
	if err != nil {
		return nil, err
	}

	var pfws []PfwRecord
	if err := json.Unmarshal([]byte(resp), &pfws); err != nil {
		return nil, err
	}

	for _, pfw := range pfws {
		if pfw.ID == id {
			return &pfw, nil
		}
	}

	return nil, nil
}

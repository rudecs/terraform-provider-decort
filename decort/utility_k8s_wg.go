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

func utilityK8sWgCheckPresence(d *schema.ResourceData, m interface{}) (*K8sNodeRecord, error) {
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))

	resp, err := controller.decortAPICall("POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	var k8s K8sRecord
	if err := json.Unmarshal([]byte(resp), &k8s); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, err
	}

	for _, wg := range k8s.Groups.Workers {
		if wg.ID == id {
			return &wg, nil
		}
	}

	return nil, nil
}

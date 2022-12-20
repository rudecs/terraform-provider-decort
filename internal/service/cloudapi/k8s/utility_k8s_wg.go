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

package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityK8sWgCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*K8SGroup, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))

	resp, err := c.DecortAPICall(ctx, "POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, err
	}

	var k8s K8SRecord
	if err := json.Unmarshal([]byte(resp), &k8s); err != nil {
		return nil, err
	}

	var id int
	if d.Id() != "" {
		id, err = strconv.Atoi(d.Id())
		if err != nil {
			return nil, err
		}
	} else {
		id = d.Get("wg_id").(int)
	}

	for _, wg := range k8s.K8SGroups.Workers {
		if wg.ID == uint64(id) {
			return &wg, nil
		}
	}

	return nil, fmt.Errorf("Not found wg with id: %v in k8s cluster: %v", id, k8s.ID)
}

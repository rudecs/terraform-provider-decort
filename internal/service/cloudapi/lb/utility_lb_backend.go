/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>

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

package lb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityLBBackendCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*Backend, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	bName := d.Get("name").(string)

	if (d.Get("lb_id").(int)) != 0 {
		urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	} else {
		parameters := strings.Split(d.Id(), "-")
		urlValues.Add("lbId", parameters[0])
		bName = parameters[1]
	}

	resp, err := c.DecortAPICall(ctx, "POST", lbGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	lb := &LoadBalancer{}
	if err := json.Unmarshal([]byte(resp), lb); err != nil {
		return nil, fmt.Errorf("can not unmarshall data to lb: %s %+v", resp, lb)
	}

	backends := lb.Backends
	for _, b := range backends {
		if b.Name == bName {
			return &b, nil
		}
	}

	return nil, fmt.Errorf("can not find backend with name: %s for lb: %d", bName, lb.ID)
}

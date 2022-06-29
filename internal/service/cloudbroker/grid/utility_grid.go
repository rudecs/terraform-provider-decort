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

package grid

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityGridCheckPresence(d *schema.ResourceData, m interface{}) (*Grid, error) {
	grid := &Grid{}
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if gridId, ok := d.GetOk("grid_id"); ok {
		urlValues.Add("gridId", strconv.Itoa(gridId.(int)))
	} else {
		return nil, errors.New("grid_id is required")
	}

	log.Debugf("utilityGridCheckPresence: load grid")
	gridRaw, err := c.DecortAPICall("POST", GridGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(gridRaw), grid)
	if err != nil {
		return nil, err
	}

	return grid, nil
}

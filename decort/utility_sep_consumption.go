/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>, <svs1370@gmail.com>

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

func utilitySepConsumptionCheckPresence(d *schema.ResourceData, m interface{}) (*SepConsumption, error) {
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	sepCons := &SepConsumption{}

	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))

	sepConsRaw, err := controller.decortAPICall("POST", sepConsumptionAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(sepConsRaw), sepCons)
	if err != nil {
		return nil, err
	}

	return sepCons, nil
}

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
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func utilityAccountConsumedUnitsByTypeCheckPresence(d *schema.ResourceData, m interface{}) (float64, error) {
	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
	urlValues.Add("cutype", strings.ToUpper(d.Get("cu_type").(string)))

	log.Debugf("utilityAccountConsumedUnitsByTypeCheckPresence")
	resultRaw, err := controller.decortAPICall("POST", accountGetConsumedUnitsByTypeAPI, urlValues)
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseFloat(resultRaw, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}

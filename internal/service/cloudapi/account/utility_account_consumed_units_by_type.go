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

package account

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityAccountConsumedUnitsByTypeCheckPresence(d *schema.ResourceData, m interface{}) (float64, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
	urlValues.Add("cutype", strings.ToUpper(d.Get("cu_type").(string)))

	log.Debugf("utilityAccountConsumedUnitsByTypeCheckPresence")
	resultRaw, err := c.DecortAPICall("POST", accountGetConsumedUnitsByTypeAPI, urlValues)
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseFloat(resultRaw, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}

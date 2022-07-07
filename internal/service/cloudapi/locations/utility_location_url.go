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

package locations

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityLocationUrlCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	c := m.(*controller.ControllerCfg)

	log.Debugf("utilityLocationUrlCheckPresence: load locations list")
	locationUrl, err := c.DecortAPICall(ctx, "POST", locationURLAPI, &url.Values{})
	if err != nil {
		return "", err
	}

	location := new(string)
	err = json.Unmarshal([]byte(locationUrl), location)
	if err != nil {
		return "", nil
	}

	return *location, nil
}

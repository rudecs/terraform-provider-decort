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

package bservice

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityBasicServiceSnapshotListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (BasicServiceSnapshots, error) {
	basicServiceSnapshotList := BasicServiceSnapshots{}
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if serviceId, ok := d.GetOk("service_id"); ok {
		urlValues.Add("serviceId", strconv.Itoa(serviceId.(int)))
	}

	log.Debugf("utilityBasicServiceSnapshotListCheckPresence")
	basicServiceSnapshotListRaw, err := c.DecortAPICall(ctx, "POST", bserviceSnapshotListAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(basicServiceSnapshotListRaw), &basicServiceSnapshotList)
	if err != nil {
		return nil, err
	}

	return basicServiceSnapshotList, nil
}

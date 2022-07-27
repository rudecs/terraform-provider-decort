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

package disks

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityDiskCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*Disk, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	disk := &Disk{}

	if d.Get("disk_id").(int) == 0 {
		urlValues.Add("diskId", d.Id())
	} else {
		urlValues.Add("diskId", strconv.Itoa(d.Get("disk_id").(int)))
	}

	log.Debugf("utilityDiskCheckPresence: load disk")
	diskRaw, err := c.DecortAPICall(ctx, "POST", disksGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(diskRaw), disk)
	if err != nil {
		return nil, err
	}

	return disk, nil
}

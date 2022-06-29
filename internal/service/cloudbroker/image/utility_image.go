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

package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityImageCheckPresence(d *schema.ResourceData, m interface{}) (*Image, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if (strconv.Itoa(d.Get("image_id").(int))) != "0" {
		urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	} else {
		urlValues.Add("imageId", d.Id())
	}

	resp, err := c.DecortAPICall("POST", imageGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	image := &Image{}
	if err := json.Unmarshal([]byte(resp), image); err != nil {
		return nil, errors.New(fmt.Sprint("Can not unmarshall data to image: ", resp, " ", image))
	}

	return image, nil
}

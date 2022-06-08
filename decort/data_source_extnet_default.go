/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceExtnetDefaultRead(d *schema.ResourceData, m interface{}) error {
	extnetId, err := utilityExtnetDefaultCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	extnetIdInt, err := strconv.ParseInt(extnetId, 10, 32)
	if err != nil {
		return err
	}
	d.Set("net_id", extnetIdInt)

	return nil
}

func dataSourceExtnetDefaultSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"net_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
	return res
}

func dataSourceExtnetDefault() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceExtnetDefaultRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceExtnetDefaultSchemaMake(),
	}
}

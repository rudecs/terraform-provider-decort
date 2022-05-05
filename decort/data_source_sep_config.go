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

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSepConfigRead(d *schema.ResourceData, m interface{}) error {
	sepConfig, err := utilitySepConfigCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())

	data, _ := json.Marshal(sepConfig)
	d.Set("config", string(data))

	return nil
}

func dataSourceSepConfigSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sep_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "storage endpoint provider ID",
		},
		"config": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "sep config json string",
		},
	}
}

func dataSourceSepConfig() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepConfigRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepConfigSchemaMake(),
	}
}

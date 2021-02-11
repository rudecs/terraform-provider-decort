/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

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
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceImageRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	// rg_id, rgid_set := d.GetOk("rg_id")
	account_id, account_set := d.GetOk("account_id")

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	if account_set {
		url_values.Add("accountId", fmt.Sprintf("%d", account_id.(int)))
	}
	body_string, err := controller.decortAPICall("POST", ImagesListAPI, url_values)
	if err != nil {
		return err
	}

	log.Debugf("dataSourceImageRead: ready to decode response body from %q", ImagesListAPI)
	model := ImagesListResp{}
	err = json.Unmarshal([]byte(body_string), &model)
	if err != nil {
		return err
	}

	// log.Printf("%#v", model)
	log.Debugf("dataSourceImageRead: traversing decoded JSON of length %d", len(model))
	for index, item := range model {
		// need to match Image by name
		if item.Name == name {
			log.Printf("dataSourceImageRead: index %d, matched name %q", index, item.Name)
			d.SetId(fmt.Sprintf("%d", item.ID))
			d.Set("account_id", item.AccountID)
			d.Set("arch", item.Arch)
			d.Set("sep_id", item.SepID)
			d.Set("pool", item.Pool)
			d.Set("status", item.Status)
			d.Set("size", item.Size)
			// d.Set("field_name", value)
			return nil
		}
	}

	return fmt.Errorf("Cannot find OS Image name %q", name)
}

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceImageRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the OS image to locate. This parameter is case sensitive.",
			},

			"account_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Optional ID of the account to limit image search to.",
			},

			"arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Binary architecture this image is created for.",
			},

			"sep_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage end-point provider serving this image.",
			},

			/*
				"sep_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Type of the storage end-point provider serving this image.",
				},
			*/

			"pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pool where this image is located.",
			},

			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the image in GB.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current model status of this disk.",
			},
		},
	}
}

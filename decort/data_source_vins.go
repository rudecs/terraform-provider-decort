/*
Copyright (c) 2020-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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

	log "github.com/sirupsen/logrus"

	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// vins_facts is a response string from API vins/get
func flattenVins(d *schema.ResourceData, vins_facts string) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceVinsExists(...) method
	log.Debugf("flattenVins: ready to decode response body from API")
	details := VinsGetResp{}
	err := json.Unmarshal([]byte(vins_facts), &details)
	if err != nil {
		return err
	}

	log.Debugf("flattenVins: decoded ViNS name:ID %s:%d, account ID %d, RG ID %d",
		details.Name, details.ID, details.AccountID, details.RgID)

	d.SetId(fmt.Sprintf("%d", details.ID))
	d.Set("rg_id", details.ID)
	d.Set("name", details.Name)
	d.Set("account_name", details.AccountName)
	d.Set("account_id", details.AccountID)
	d.Set("grid_id", details.GridID)
	d.Set("description", details.Desc)
	d.Set("status", details.Status)
	d.Set("def_net_type", details.DefaultNetType)
	d.Set("def_net_id", details.DefaultNetID)

	return nil
}

func dataSourceVinsRead(d *schema.ResourceData, m interface{}) error {
	vinsFacts, err := utilityVinsCheckPresence(d, m)
	if vinsFacts == "" {
		// if empty string is returned from utilityVinsCheckPresence then there is no
		// such ViNS and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty in this case
		return err
	}

	return flattenVins(d, vinsFacts)
}

func dataSourceVins() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceVinsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the ViNS. Names are case sensitive and unique within the context of an account or resource group.",
			},

			/*
			"vins_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the ViNS. If ViNS ID is specified, then ViNS name, rg_id and account_id are ignored.",
			},
			*/

			"rg_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the resource group, where this ViNS is belongs to (for ViNS created at resource group level, 0 otherwise).",
			},

			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the account, which this ViNS belongs to.",
			},

			// the rest of attributes are computed
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the account, which this ViNS belongs to.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-defined text description of this ViNS.",
			},

			"ext_ip_addr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the external connection (valid for ViNS connected to external network, empty string otherwise).",
			},

			"ext_net_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the external network this ViNS is connected to (-1 means no external connection).",
			},

			"ipcidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network address used by this ViNS."
			},
		},
	}
}

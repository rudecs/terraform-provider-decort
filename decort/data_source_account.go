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
	// "net/url"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func flattenAccount(d *schema.ResourceData, acc_facts string) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceAccountExists(...) method

	// log.Debugf("flattenAccount: ready to decode response body from %q", CloudspacesGetAPI)
	details := AccountRecord{}
	err := json.Unmarshal([]byte(acc_facts), &details)
	if err != nil {
		return err
	}

	log.Debugf("flattenAccount: decoded Account name %q / ID %d, status %q", details.Name, details.ID, details.Status)

	d.SetId(fmt.Sprintf("%d", details.ID))
	d.Set("name", details.Name)
	d.Set("status", details.Status)

	return nil
}

func dataSourceAccountRead(d *schema.ResourceData, m interface{}) error {
	acc_facts, err := utilityAccountCheckPresence(d, m)
	if acc_facts == "" {
		// if empty string is returned from utilityAccountCheckPresence then there is no
		// such account and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty in this case
		return err
	}

	return flattenAccount(d, acc_facts)
}

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the account. Names are case sensitive and unique.",
			},

			"account_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the account. If account ID is specified, then account name is ignored.",
			},

			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the account.",
			},

			/* We keep the following attributes commented out, as we are not implementing account
			   management with Terraform plugin, so we do not need this extra info.

			"quota": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: quotaRgSubresourceSchema(), // this is a dictionary
				},
				Description: "Quotas on the resources for this account and all its resource groups.",
			},

			"resource_groups": {
				Type:         schema.TypeList,
				Computed:     true,
				Elem: &schema.Schema {
					Type:  schema.TypeInt,
				},
				Description:  "IDs of resource groups in this account."
			},

			"vins": {
				Type:         schema.TypeList,
				Computed:     true,
				Elem: &schema.Schema {
					Type:  schema.TypeInt,
				},
				Description:  "IDs of VINSes created at the account level."
			},
			*/
		},
	}
}

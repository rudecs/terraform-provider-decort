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
	// "net/url"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"

)

func flattenResgroup(d *schema.ResourceData, rg_facts string) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceRsgroupExists(...) method
	log.Debugf("%s", rg_facts)
	log.Debugf("flattenResgroup: ready to decode response body from %q", CloudspacesGetAPI)
	details := ResgroupGetResp{}
	err := json.Unmarshal([]byte(rg_facts), &details)
	if err != nil {
		return err
	}

	log.Debugf("flattenResgroup: decoded ResGroup name %q / ID %d, account ID %d, public IP %q", 
				details.Name, details.ID, details.AccountID, details.PublicIP)

	d.SetId(fmt.Sprintf("%d", details.ID))
	d.Set("name", details.Name)
	d.Set("account_id", details.AccountID)
	d.Set("grid_id", details.GridID)

	log.Debugf("flattenResgroup: calling flattenQuota()")
	if err = d.Set("quotas", flattenQuota(details.Quotas)); err != nil {
		return err
	}

	return nil
}

func dataSourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	rg_facts, err := utilityResgroupCheckPresence(d, m)
	if rg_facts == "" {
		// if empty string is returned from utilityResgroupCheckPresence then there is no
		// such resource group and err tells so - just return it to the calling party 
		d.SetId("") // ensure ID is empty in this case
		return err
	}

	return flattenResgroup(d, rg_facts)
}


func dataSourceResgroup() *schema.Resource {
	return &schema.Resource {
		SchemaVersion: 1,

		Read:   dataSourceResgroupRead,

		Timeouts: &schema.ResourceTimeout {
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema {
			"name": {
				Type:          schema.TypeString,
				Required:      true,
				Description:  "Name of this resource group. Names are case sensitive and unique within the context of a tenant.",
			},

			"account": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the tenant, which this resource group belongs to.",
			},

			"account_id": &schema.Schema {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the tenant, which this resource group belongs to.",
			},

			"grid_id": &schema.Schema {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the grid, where this resource group is deployed.",
			},

			"public_ip": {  // this may be obsoleted as new network segments and true resource groups are implemented
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Public IP address of this resource group (if any).",
			},

			"quotas": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        &schema.Resource {
					Schema:  quotasSubresourceSchema(),
				},
				Description: "Quotas on the resources for this resource group.",
			},
		},
	}
}

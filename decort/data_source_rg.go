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

	log "github.com/sirupsen/logrus"

	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func flattenResgroup(d *schema.ResourceData, rg_facts string) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceRsgroupExists(...) method
	// log.Debugf("%s", rg_facts)
	log.Debugf("flattenResgroup: ready to decode response body from API")
	details := ResgroupGetResp{}
	err := json.Unmarshal([]byte(rg_facts), &details)
	if err != nil {
		return err
	}

	log.Debugf("flattenResgroup: decoded RG name %q / ID %d, account ID %d",
		details.Name, details.ID, details.AccountID)

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
	d.Set("vins", details.Vins)
	d.Set("computes", details.Computes)

	log.Debugf("flattenResgroup: calling flattenQuota()")
	if err = d.Set("quota", parseQuota(details.Quota)); err != nil {
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
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceResgroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the resource group. Names are case sensitive and unique within the context of an account.",
			},

			"rg_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the resource group. If this ID is specified, then resource group name is ignored.",
			},

			"account_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the account, which this resource group belongs to.",
			},

			"account_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the account, which this resource group belongs to. If account ID is specified, then account name is ignored.",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-defined text description of this resource group.",
			},

			"grid_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the grid, where this resource group is deployed.",
			},

			"quota": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: quotaRgSubresourceSchemaMake(), // this is a dictionary
				},
				Description: "Quota settings for this resource group.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of this resource group.",
			},

			"def_net_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the default network for this resource group.",
			},

			"def_net_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the default network for this resource group (if any).",
			},

			"vins": {
				Type:     schema.TypeList, // this is a list of ints
				Computed: true,
				MaxItems: LimitMaxVinsPerResgroup,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of VINs deployed in this resource group.",
			},

			"computes": {
				Type:     schema.TypeList, //t his is a list of ints
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of computes deployed in this resource group.",
			},
		},
	}
}

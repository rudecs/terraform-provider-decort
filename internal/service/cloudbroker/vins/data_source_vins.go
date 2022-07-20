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

package vins

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rudecs/terraform-provider-decort/internal/constants"
	log "github.com/sirupsen/logrus"

	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// vins_facts is a response string from API vins/get
func flattenVins(d *schema.ResourceData, vins_facts string) diag.Diagnostics {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceVinsExists(...) method
	// log.Debugf("flattenVins: ready to decode response body from API %s", vins_facts)
	vinsRecord := VinsRecord{}
	err := json.Unmarshal([]byte(vins_facts), &vinsRecord)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("flattenVins: decoded ViNS name:ID %s:%d, account ID %d, RG ID %d",
		vinsRecord.Name, vinsRecord.ID, vinsRecord.AccountID, vinsRecord.RgID)

	d.SetId(fmt.Sprintf("%d", vinsRecord.ID))
	d.Set("name", vinsRecord.Name)
	d.Set("account_id", vinsRecord.AccountID)
	d.Set("account_name", vinsRecord.AccountName)
	d.Set("rg_id", vinsRecord.RgID)
	d.Set("description", vinsRecord.Desc)
	d.Set("ipcidr", vinsRecord.IPCidr)

	noExtNetConnection := true
	for _, value := range vinsRecord.VNFs {
		if value.Type == "GW" {
			log.Debugf("flattenVins: discovered GW VNF ID %d in ViNS ID %d", value.ID, vinsRecord.ID)
			extNetID, idOk := value.Config["ext_net_id"] // NOTE: unknown numbers are unmarshalled to float64. This is by design!
			extNetIP, ipOk := value.Config["ext_net_ip"]
			if idOk && ipOk {
				log.Debugf("flattenVins: ViNS ext_net_id=%d, ext_net_ip=%s", int(extNetID.(float64)), extNetIP.(string))
				d.Set("ext_ip_addr", extNetIP.(string))
				d.Set("ext_net_id", int(extNetID.(float64)))
			} else {
				return diag.Errorf("Failed to unmarshal VNF GW Config - structure is invalid.")
			}
			noExtNetConnection = false
			break
		}
	}

	if noExtNetConnection {
		d.Set("ext_ip_addr", "")
		d.Set("ext_net_id", 0)
	}

	log.Debugf("flattenVins: EXTRA CHECK - schema rg_id=%d, ext_net_id=%d", d.Get("rg_id").(int), d.Get("ext_net_id").(int))

	return nil
}

func dataSourceVinsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vinsFacts, err := utilityVinsCheckPresence(ctx, d, m)
	if vinsFacts == "" {
		// if empty string is returned from utilityVinsCheckPresence then there is no
		// such ViNS and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty in this case
		return diag.FromErr(err)
	}

	return flattenVins(d, vinsFacts)
}

func DataSourceVins() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceVinsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
				Description: "Network address used by this ViNS.",
			},
		},
	}
}

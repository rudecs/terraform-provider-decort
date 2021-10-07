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

package decort

import (

	"encoding/json"
	"fmt"
	// "hash/fnv"
	log "github.com/sirupsen/logrus" 
	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)


func flattenPfw(d *schema.ResourceData, pfwFacts string) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourcePfwExists(...) method
	// log.Debugf("flattenPfw: ready to decode response body from API %s", pfwFacts)
	pfwRecord := ComputePfwListResp{}
	err := json.Unmarshal([]byte(pfwFacts), &pfwRecord)
	if err != nil {
		return err
	}

	log.Debugf("flattenPfw: decoded %d PFW rules for compute ID %s on ViNS ID %d",
	           len(pfwRecord.Rules), pfwRecord.Header.VinsID, pfwRecord.Header.VinsID)

	/*
	combo := fmt.Sprintf("%d:%d", compId.(int), pfwRecord.ViNS.VinsID)
	hasher := fnv.New32a()
	hasher.Write([]byte(combo))
	d.SetId(fmt.Sprintf("%d", hasher.Sum32()))
	*/
	// set ID of this PFW rule set as "compute_id:vins_id"
	d.SetId(fmt.Sprintf("%d:%d", pfwRecord.Header.ComputeID, pfwRecord.Header.VinsID))
	log.Debugf("flattenPfw: PFW rule set ID %s", d.Id())
	d.Set("compute_id", pfwRecord.Header.ComputeID)
	d.Set("vins_id", pfwRecord.Header.VinsID)

	pfwRulesList := []interface{}{}
	for _, runner := range pfwRecord.Rules {
		rule := map[string]interface{}{
			"pub_port_start": runner.PublicPortStart,
			"pub_port_end":   runner.PublicPortEnd,
			"local_port":     runner.LocalPort,
			"proto":          runner.Protocol,
			"rule_id":        runner.ID,
		}
		pfwRulesList = append(pfwRulesList, rule) 
	}
	if err = d.Set("rule", pfwRulesList); err != nil {
		return err
	}

	return nil
}

func dataSourcePfwRead(d *schema.ResourceData, m interface{}) error {
	pfwFacts, err := utilityPfwCheckPresence(d, m)
	if pfwFacts == "" {
		// if empty string is returned from dataSourcePfwRead then we got no
		// PFW rules. It could also be because there was some error, which
		// is indicated by non-nil err value
		d.SetId("") // ensure ID is empty in this case anyway
		return err
	}

	return flattenPfw(d, pfwFacts)
}

func dataSourcePfw() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourcePfwRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},
		
		Schema: map[string]*schema.Schema{
			"compute_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the compute instance to configure port forwarding rules for.",
			},

			"vins_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the ViNS to configure port forwarding rules on. Compute must be already plugged into this ViNS and ViNS must have external network connection.",
			},

			"rule": {
				Type:         schema.TypeSet,
				Optional:     true,
				Elem:         &schema.Resource{
					Schema: rulesSubresourceSchemaMake(),
				},
				Description:  "Port forwarding rule. You may specify several rules, one in each such block.",
			},
		},
	}
}

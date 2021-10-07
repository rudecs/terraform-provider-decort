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
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func utilityPfwCheckPresence(d *schema.ResourceData, m interface{}) (string, error) {
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for the Terraform resource Exists method.
	//

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}

	// NOTE on importing PFW into TF state resource:
	//
	// Port forward rules are NOT represented by any "individual" resource in the platform.
	// Consequently, there is no unique ID reported by the platform that could be used to
	// identify PFW rule set. 
	// However, we need some ID to identify PFW resource in TF state, and compute ID is the most
	// convenient way, as it is:
	// 1) unique;
	// 2) compute may have only one PFW rule set.
	//

	var compId, vinsId int
	
	if d.Id() != "" {
		log.Debugf("utilityPfwCheckPresence: setting context from d.Id() %s", d.Id())
		idParts := strings.SplitN(d.Id(), ":", 2)
		compId, _ = strconv.Atoi(idParts[0])
		vinsId, _ = strconv.Atoi(idParts[1])
		log.Debugf("utilityPfwCheckPresence: extracted Compute ID %d, ViNS %d", compId, vinsId)
		if compId <= 0 || vinsId <= 0 {
			return "", fmt.Errorf("Ivalid context from d.Id %s", d.Id())
		}
	} else {
		scId, cSet := d.GetOk("compute_id")
		svId, vSet := d.GetOk("vins_id")
		if cSet || vSet {
			log.Debugf("utilityPfwCheckPresence: setting Compute ID from schema")
			compId = scId.(int)
			vinsId = svId.(int)
			log.Debugf("utilityPfwCheckPresence: extractted Compute ID %d, ViNS %d", compId, vinsId)
		} else {
			return "", fmt.Errorf("Cannot get context to check PFW rules neither from d.Id() nor from schema")
		}
	}
	
	log.Debugf("utilityPfwCheckPresence: preparing to get PFW rules for Compute ID %d on ViNS ID %d", compId, vinsId)

	urlValues.Add("computeId", fmt.Sprintf("%d", compId))
	apiResp, err, respCode := controller.decortAPICall("POST", ComputePfwListAPI, urlValues)
	if respCode == 500 {
		// this is workaround for API 3.7.0 "feature" - will be removed in one of the future versions
		log.Errorf("utilityPfwCheckPresence: Compute ID %d has no PFW and no connection to PFW-ready ViNS", compId)
		return "", nil
	}
	if err != nil {
		
		return "", err
	}
	
	pfwListResp := ComputePfwListResp{}

	// Note the specifics of compute/pfwList response in API 3.7.x (this may be changed in the future):
	// 1) if there are no PFW rules and compute is not connected to any PFW-able ViNS 
	//    the response will be empty string (or HTTP error code 500)
	// 2) if there are no PFW rules but compute is connected to a PFW-able ViNS 
	//    the response will contain a list with a single element - prefix (see PfwPrefixRecord) 
	// 3) if there are port forwarding rules, the response will contain a list which starts 
	//    with prefix (see PfwPrefixRecord) and then followed by one or more rule records 
	//    (see PfwRuleRecord)
	//
	// EXTRA NOTE: in API 3.7.0 and the likes pfwList returns HTTP response code 500 for a compute
	// that is not connected to any PFW-able ViNS - need to implement temporary workaround 
	
	if apiResp == "" {
		// No port forward rules defined for this compute
		return "", nil
	}

	log.Debugf("utilityPfwCheckPresence: ready to split API response string %s", apiResp)

	twoParts := strings.SplitN(apiResp, "},", 2)
	if len(twoParts) < 1 || len(twoParts) > 2 {
		// Case: invalid format of API response
		log.Errorf("utilityPfwCheckPresence: non-empty pfwList response for compute ID %d failed to split properly", compId)
		return "", fmt.Errorf("Non-empty pfwList response failed to split properly")
	}

	if len(twoParts) == 1 {
		// Case: compute is connected to a PWF-ready ViNS but has no PFW rules defined
		log.Debugf("utilityPfwCheckPresence: compute ID %d is connected to PFW-ready ViNS but has no PFW rules", compId)
		return "", nil
	}

	// Case: compute is connected to a PFW ready ViNS and has some PFW rule
	prefixResp := strings.TrimSuffix(strings.TrimPrefix(twoParts[0], "["), ",") + "}"
	log.Debugf("utilityPfwCheckPresence: ready to unmarshal prefix part %s", prefixResp)
	err = json.Unmarshal([]byte(prefixResp), &pfwListResp.Header)
	if err != nil {
		log.Errorf("utilityPfwCheckPresence: failed to unmarshal prefix part of API response: %s", err)
		return "", err
	}

	rulesResp := "[" + twoParts[1]
	log.Debugf("utilityPfwCheckPresence: ready to unmarshal rules part %s", rulesResp)
	err = json.Unmarshal([]byte(rulesResp), &pfwListResp.Rules)
	if err != nil {
		log.Errorf("utilityPfwCheckPresence: failed to unmarshal rules part of API response: %s", err)
		return "", err
	}

	log.Debugf("utilityPfwCheckPresence: successfully read %d port forward rules for Compute ID %d, ViNS ID %d",
               len(pfwListResp.Rules), compId, pfwListResp.Header.VinsID)
	
	if pfwListResp.Header.VinsID != vinsId {
		log.Errorf("utilityPfwCheckPresence: ViNS ID mismatch for PFW rules on compute ID %d: actual %d, required %d",
				   compId, pfwListResp.Header.VinsID, vinsId)
		return "", fmt.Errorf("ViNS ID mismatch for PFW rules on compute ID %d: actual %d, required %d",
							  compId, pfwListResp.Header.VinsID, vinsId)
	}

	// reconstruct API response string for return
	pfwListResp.Header.ComputeID = compId
	reencodedItem, err := json.Marshal(pfwListResp)
	if err != nil {
		return "", err
	}

	return string(reencodedItem[:]), nil 
}

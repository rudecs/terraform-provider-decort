/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

package rg

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

func flattenAccountSeps(seps map[string]map[string]ResourceSep) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for sepKey, sepVal := range seps {
		for dataKey, dataVal := range sepVal {
			temp := map[string]interface{}{
				"sep_id":        sepKey,
				"data_name":     dataKey,
				"disk_size":     dataVal.DiskSize,
				"disk_size_max": dataVal.DiskSizeMax,
			}
			res = append(res, temp)
		}
	}
	return res
}

func flattenAccResource(r Resource) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cpu":        r.CPU,
		"disksize":   r.Disksize,
		"extips":     r.Extips,
		"exttraffic": r.Exttraffic,
		"gpu":        r.GPU,
		"ram":        r.RAM,
		"seps":       flattenAccountSeps(r.SEPs),
	}
	res = append(res, temp)
	return res
}

func flattenRgResources(r Resources) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"current":  flattenAccResource(r.Current),
		"reserved": flattenAccResource(r.Reserved),
	}
	res = append(res, temp)
	return res
}

func flattenDataResgroup(d *schema.ResourceData, details ResgroupGetResp) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceRsgroupExists(...) method
	// log.Debugf("%s", rg_facts)

	log.Debugf("flattenResgroup: decoded RG name %q / ID %d, account ID %d",
		details.Name, details.ID, details.AccountID)

	d.SetId(fmt.Sprintf("%d", details.ID))
	d.Set("rg_id", details.ID)
	d.Set("name", details.Name)
	d.Set("account_name", details.AccountName)
	d.Set("account_id", details.AccountID)
	d.Set("gid", details.GridID)
	d.Set("description", details.Desc)
	d.Set("status", details.Status)
	d.Set("def_net_type", details.DefaultNetType)
	d.Set("def_net_id", details.DefaultNetID)
	d.Set("resources", flattenRgResources(details.Resources))
	d.Set("vins", details.Vins)
	d.Set("vms", details.Computes)
	log.Debugf("flattenResgroup: calling flattenQuota()")
	if err := d.Set("quota", parseQuota(details.Quota)); err != nil {
		return err
	}

	return nil
}

func flattenResgroup(d *schema.ResourceData, details ResgroupGetResp) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceRsgroupExists(...) method
	// log.Debugf("%s", rg_facts)
	//log.Debugf("flattenResgroup: ready to decode response body from API")
	//details := ResgroupGetResp{}
	//err := json.Unmarshal([]byte(rg_facts), &details)
	//if err != nil {
	//return err
	//}

	log.Debugf("flattenResgroup: decoded RG name %q / ID %d, account ID %d",
		details.Name, details.ID, details.AccountID)

	d.SetId(fmt.Sprintf("%d", details.ID))
	d.Set("rg_id", details.ID)
	d.Set("name", details.Name)
	d.Set("account_name", details.AccountName)
	d.Set("account_id", details.AccountID)
	d.Set("gid", details.GridID)
	d.Set("description", details.Desc)
	d.Set("status", details.Status)
	d.Set("def_net_type", details.DefaultNetType)
	d.Set("def_net_id", details.DefaultNetID)
	d.Set("resources", flattenRgResources(details.Resources))
	d.Set("vins", details.Vins)
	d.Set("vms", details.Computes)
	log.Debugf("flattenResgroup: calling flattenQuota()")
	if err := d.Set("quota", parseQuota(details.Quota)); err != nil {
		return err
	}

	return nil
}

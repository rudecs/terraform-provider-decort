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

	// "encoding/json"
	// "fmt"
	// "log"
	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func makeQuotaRecord(arg_list []interface{}) QuotaRecord {
	quota := QuotaRecord{
		Cpu:        -1,
		Ram:        -1., // this is float64, but may change in the future
		Disk:       -1,
		ExtTraffic: -1,
		ExtIPs:     -1,
		GpuUnits:   -1,
	}
	subres_data := arg_list[0].(map[string]interface{})

	if subres_data["cpu"].(int) > 0 {
		quota.Cpu = subres_data["cpu"].(int)
	}

	if subres_data["disk"].(int) > 0 {
		quota.Disk = subres_data["disk"].(int) // Disk capacity ib GB
	}

	if subres_data["ram"].(float64) > 0 {
		quota.Ram = subres_data["ram"].(float64) // RAM volume in MB, as float64!
	}

	if subres_data["ext_traffic"].(int) > 0 {
		quota.ExtTraffic = subres_data["ext_traffic"].(int)
	}

	if subres_data["ext_ips"].(int) > 0 {
		quota.ExtIPs = subres_data["ext_ips"].(int)
	}

	if subres_data["gpu_units"].(int) > 0 {
		quota.GpuUnits = subres_data["gpu_units"].(int)
	}

	return quota
}

func parseQuota(quota QuotaRecord) []interface{} {
	quota_map := make(map[string]interface{})

	quota_map["cpu"] = quota.Cpu
	quota_map["ram"] = quota.Ram // NB: this is float64, unlike the rest of values
	quota_map["disk"] = quota.Disk
	quota_map["ext_traffic"] = quota.ExtTraffic
	quota_map["ext_ips"] = quota.ExtIPs
	quota_map["gpu_units"] = quota.GpuUnits

	result := make([]interface{}, 1)
	result[0] = quota_map

	return result // this result will be used to d.Set("quota,") of dataSourceResgroup schema
}

func quotaRgSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"cpu": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "Limit on the total number of CPUs in this resource group.",
		},

		"ram": {
			Type:        schema.TypeFloat, // NB: API expects and returns this as float in units of MB! This may be changed in the future.
			Optional:    true,
			Default:     -1.,
			Description: "Limit on the total amount of RAM in this resource group, specified in MB.",
		},

		"disk": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "Limit on the total volume of storage resources in this resource group, specified in GB.",
		},

		"ext_traffic": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "Limit on the total ingress network traffic for this resource group, specified in GB.",
		},

		"ext_ips": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "Limit on the total number of external IP addresses this resource group can use.",
		},

		"gpu_units": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "Limit on the total number of virtual GPUs this resource group can use.",
		},
	}
	return rets
}

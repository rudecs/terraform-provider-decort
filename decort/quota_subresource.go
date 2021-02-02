/*
Copyright (c) 2019 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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

package decs

import (

	// "encoding/json"
	// "fmt"
	// "log"
	// "net/url"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"

)

func makeQuotaConfig(arg_list []interface{}) (ResgroupQuotaConfig, int) {
	quota := ResgroupQuotaConfig{
		Cpu:        -1,
		Ram:        -1,
		Disk:       -1,
		NetTraffic: -1,
		ExtIPs:     -1,
	}
	subres_data := arg_list[0].(map[string]interface{})

	if subres_data["cpu"].(int) > 0 {
		quota.Cpu = subres_data["cpu"].(int)
	}

	if subres_data["disk"].(int) > 0 {
		quota.Disk = subres_data["disk"].(int)
	}

	if subres_data["ram"].(int) > 0 {
		ram_limit := subres_data["ram"].(int)
		quota.Ram = float32(ram_limit) // /1024 // legacy fix - this can be obsoleted once redmine FR #1465 is implemented
	}

	if subres_data["net_traffic"].(int) > 0 {
		quota.NetTraffic = subres_data["net_traffic"].(int)
	}

	if subres_data["ext_ips"].(int) > 0 {
		quota.ExtIPs = subres_data["ext_ips"].(int)
	}

	return quota, 1
} 

func flattenQuota(quotas QuotaRecord) []interface{} {
	quotas_map :=  make(map[string]interface{})

	quotas_map["cpu"] = quotas.Cpu
	quotas_map["ram"] = int(quotas.Ram)
	quotas_map["disk"] = quotas.Disk
	quotas_map["net_traffic"] = quotas.NetTraffic
	quotas_map["ext_ips"] = quotas.ExtIPs

	result := make([]interface{}, 1)
	result[0] = quotas_map

	return result
}

func quotasSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"cpu": &schema.Schema {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total number of CPUs in this resource group.",
		},

		"ram": &schema.Schema {
			Type:        schema.TypeInt, // NB: API expects and returns this as float! This may be changed in the future.
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total amount of RAM in this resource group, specified in GB (Gigabytes!).",
			},

		"disk": &schema.Schema {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total volume of storage resources in this resource group, specified in GB.",
		},

		"net_traffic": &schema.Schema {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total ingress network traffic for this resource group, specified in GB.",
		},

		"ext_ips": &schema.Schema {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total number of external IP addresses this resource group can use.",
		},
	}
	return rets
}
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

	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func makeDisksConfig(arg_list []interface{}) (disks []DiskConfig, count int) {
	count = len(arg_list) 
	if count < 1 {
		return nil, 0
	}

	// allocate DataDisks list and fill it 
	disks = make([]DiskConfig, count)
	var subres_data map[string]interface{}
	for index, value := range arg_list {
		subres_data = value.(map[string]interface{})
		disks[index].Label = subres_data["label"].(string)
		disks[index].Size = subres_data["size"].(int)
		disks[index].Pool = subres_data["pool"].(string)
		disks[index].Provider = subres_data["provider"].(string)
	}

	return disks, count
}

func flattenDataDisks(disks []DataDiskRecord) []interface{} {
	var length = 0
	for _, value := range disks {
		if value.DiskType == "D" {
			length += 1
		}
	}
	log.Printf("flattenDataDisks: found %d disks with D type", length)

	result := make([]interface{}, length)
	if length == 0 {
		return result
	}

	elem := make(map[string]interface{})

	var subindex = 0
	for _, value := range disks {
		if value.DiskType == "D" {
			elem["label"] = value.Label
			elem["size"] = value.SizeMax
			elem["disk_id"] = value.ID
			elem["pool"] = "default"
			elem["provider"] = "default"
			result[subindex] = elem
			subindex += 1
		}
		
	}

	return result
}

/*
func makeDataDisksArgString(disks []DiskConfig) string {
	// Prepare a string with the sizes of data disks for the virtual machine.
	// It is designed to be passed as "datadisks" argument of virtual machine create API call.
	if len(disks) < 1 {
		return ""
	}
	return ""
}
*/

func diskSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Unique label to identify this disk among other disks connected to this VM.",
		},

		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntAtLeast(1),
			Description: "Size of the disk in GB.",
		},

		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			Description: "Pool from which this disk should be provisioned.",
		},

		"provider": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			Description: "Storage provider (storage technology type) by which this disk should be served.",
		},

		"disk_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of this disk resource.",
		},
		
	}

	return rets
}

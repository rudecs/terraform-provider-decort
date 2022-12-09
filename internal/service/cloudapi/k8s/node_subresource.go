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

package k8s

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func nodeMasterDefault() K8sNodeRecord {
	return K8sNodeRecord{
		Num:  1,
		Cpu:  2,
		Ram:  2048,
		Disk: 0,
	}
}

func nodeWorkerDefault() K8sNodeRecord {
	return K8sNodeRecord{
		Num:  1,
		Cpu:  1,
		Ram:  1024,
		Disk: 0,
	}
}

func parseNode(nodeList []interface{}) K8sNodeRecord {
	node := nodeList[0].(map[string]interface{})

	return K8sNodeRecord{
		Num:  node["num"].(int),
		Cpu:  node["cpu"].(int),
		Ram:  node["ram"].(int),
		Disk: node["disk"].(int),
	}
}

func nodeToResource(node K8sNodeRecord) []interface{} {
	mp := make(map[string]interface{})

	mp["num"] = node.Num
	mp["cpu"] = node.Cpu
	mp["ram"] = node.Ram
	mp["disk"] = node.Disk

	return []interface{}{mp}
}

func nodeK8sSubresourceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"num": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of nodes to create.",
		},

		"cpu": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Node CPU count.",
		},

		"ram": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Node RAM in MB.",
		},

		"disk": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Node boot disk size in GB.",
		},
	}
}

func mastersSchemaMake() map[string]*schema.Schema {
	masters := masterGroupSchemaMake()
	masters["num"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Number of nodes to create.",
	}
	masters["cpu"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "Node CPU count.",
	}
	masters["ram"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "Node RAM in MB.",
	}
	masters["disk"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "Node boot disk size in GB.",
	}
	return masters
}

func workersSchemaMake() map[string]*schema.Schema {
	workers := k8sGroupListSchemaMake()
	workers["num"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Number of nodes to create.",
	}
	workers["cpu"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "Node CPU count.",
	}
	workers["ram"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "Node RAM in MB.",
	}
	workers["disk"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "Node boot disk size in GB.",
	}
	return workers
}

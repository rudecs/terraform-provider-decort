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
 This file contains definitions and code for handling Interface component of Compute schema
*/

package decort

import (
	/*
		"log"
		"strconv"
		"strings"
	*/

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func interfaceSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"net_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the network entity this interface is connected to.",
		},

		"net_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of the network entity this interface is connected to.",
		},

		"ip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP addresses assigned to this interface.",
		},

		"netmask": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Network mask to be used with this interface.",
		},

		"mac": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MAC address of this interface.",
		},

		"default_gw": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Default gateway associated with this interface.",
		},

		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Interface name.",
		},

		"connection_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "VxLAN or VLAN ID this interface is connected to.",
		},

		"connection_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of the segment (VLAN or VxLAN) this interface is connected to.",
		},

		"qos": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: interfaceQosSubresourceSchemaMake(),
			},
			Description: "QoS settings for this interface.",
		},
	}

	return rets
}

func interfaceQosSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"egress_rate": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Egress rate limit on this interface.",
		},

		"ingress_burst": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Ingress burst limit on this interface.",
		},

		"ingress_rate": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Ingress rate limit on this interface.",
		},

		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "GUID of this QoS record.",
		},
	}

	return rets
}

/*
func flattenNetworks(nets []NicRecord) []interface{} {
	// this function expects an array of NicRecord as returned by machines/get API call
	// NOTE: it does NOT expect a strucutre as returned by externalnetwork/list
	var length = 0
	var strarray []string

	for _, value := range nets {
		if value.NicType == "PUBLIC" {
			length += 1
		}
	}
	log.Debugf("flattenNetworks: found %d NICs with PUBLIC type", length)

	result := make([]interface{}, length)
	if length == 0 {
		return result
	}

	elem := make(map[string]interface{})

	var subindex = 0
	for index, value := range nets {
		if value.NicType == "PUBLIC" {
			// this will be changed as network segments entity
			// value.Params for ext net comes in a form "gateway:176.118.165.1 externalnetworkId:6"
			// for network_id we need to extract from this string
			strarray = strings.Split(value.Params, " ")
			substr := strings.Split(strarray[1], ":")
			elem["network_id"], _ = strconv.Atoi(substr[1])
			elem["ip_range"] = value.IPAddress
			// elem["label"] = ... - should be uncommented for the future release
			log.Debugf("flattenNetworks: parsed element %d - network_id %d, ip_range %s",
				index, elem["network_id"].(int), value.IPAddress)
			result[subindex] = elem
			subindex += 1
		}
	}

	return result
}

func makePortforwardsConfig(arg_list []interface{}) (pfws []PortforwardConfig, count int) {
	count = len(arg_list)
	if count < 1 {
		return nil, 0
	}

	pfws = make([]PortforwardConfig, count)
	var subres_data map[string]interface{}
	for index, value := range arg_list {
		subres_data = value.(map[string]interface{})
		// pfws[index].Label = subres_data["label"].(string) - should be uncommented for future release
		pfws[index].ExtPort = subres_data["ext_port"].(int)
		pfws[index].IntPort = subres_data["int_port"].(int)
		pfws[index].Proto = subres_data["proto"].(string)
	}

	return pfws, count
}

func flattenPortforwards(pfws []PortforwardRecord) []interface{} {
	result := make([]interface{}, len(pfws))
	elem := make(map[string]interface{})
	var port_num int

	for index, value := range pfws {
		// elem["label"] = ... - should be uncommented for the future release

		// external port field is of TypeInt in the portforwardSubresourceSchema, but string is returned
		// by portforwards/list API, so we need conversion here
		port_num, _ = strconv.Atoi(value.ExtPort)
		elem["ext_port"] = port_num
		// internal port field is of TypeInt in the portforwardSubresourceSchema, but string is returned
		// by portforwards/list API, so we need conversion here
		port_num, _ = strconv.Atoi(value.IntPort)
		elem["int_port"] = port_num
		elem["proto"] = value.Proto
		elem["ext_ip"] = value.ExtIP
		elem["int_ip"] = value.IntIP
		result[index] = elem
	}

	return result
}

func portforwardSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Unique label of this network connection to identify it amnong other connections for this VM.",
		},

		"ext_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "External port number for this port forwarding rule.",
		},

		"int_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "Internal port number for this port forwarding rule.",
		},

		"proto": {
			Type:     schema.TypeString,
			Required: true,
			// ValidateFunc: validation.IntBetween(1, ),
			Description: "Protocol type for this port forwarding rule. Should be either 'tcp' or 'udp'.",
		},

		"ext_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ".",
		},

		"int_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ".",
		},
	}

	return rets
}

func flattenNICs(nics []NicRecord) []interface{} {
	var result = make([]interface{}, len(nics))
	elem := make(map[string]interface{})

	for index, value := range nics {
		elem["status"] = value.Status
		elem["type"] = value.NicType
		elem["mac"] = value.MacAddress
		elem["ip_address"] = value.IPAddress
		elem["parameters"] = value.Params
		elem["reference_id"] = value.ReferenceID
		elem["network_id"] = value.NetworkID
		result[index] = elem
	}

	return result
}

func nicSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current status of this NIC.",
		},

		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of this NIC.",
		},

		"mac": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MAC address assigned to this NIC.",
		},

		"ip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP address assigned to this NIC.",
		},

		"parameters": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Additional NIC parameters.",
		},

		"reference_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Reference ID of this NIC.",
		},

		"network_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Network ID which this NIC is connected to.",
		},
	}

	return rets
}

*/

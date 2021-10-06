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
	// "bytes"
	// log "github.com/sirupsen/logrus" 
	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// This is subresource of network subresource of compute resource used 
// when creating/managing port forwarding rules for a compute connected 
// to the corresponding network
// It only applies to a ViNS connection AND to a ViNS with external network connection

func pfwSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"pub_port_start": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description: "Port number on the external interface. For a ranged rule it set the starting port number.",
		},

		"pub_port_end": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description: "End port number on the external interface for a ranged rule. Set it equal to start port for a single port rule.",
		},

		"local_port": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description: "Port number on the local interface.",
		},

		"proto": {
			Type:        schema.TypeString,
			Required:    true,
			StateFunc:   stateFuncToLower,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false), 
			Description: "Protocol for this rule. Could be either tcp or udp.",
		},

	}
	return rets
}

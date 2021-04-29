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
	log "github.com/sirupsen/logrus" 
	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// This is subresource of compute resource used when creating/managing compute network connections

func networkSubresIPAddreDiffSupperss(key, oldVal, newVal string, d *schema.ResourceData) bool {
	if newVal != "" && newVal != oldVal {
		log.Debugf("networkSubresIPAddreDiffSupperss: key=%s, oldVal=%q, newVal=%q -> suppress=FALSE", key, oldVal, newVal)
		return false
	}
	log.Debugf("networkSubresIPAddreDiffSupperss: key=%s, oldVal=%q, newVal=%q -> suppress=TRUE", key, oldVal, newVal)
	return true // suppress difference
}

func networkSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"net_type": {
			Type:        schema.TypeString,
			Required:    true,
			StateFunc:   stateFuncToUpper,
			ValidateFunc: validation.StringInSlice([]string{"EXTNET", "VINS"}, false), // observe case while validating
			Description: "Type of the network for this connection, either EXTNET or VINS.",
		},

		"net_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the network for this connection.",
		},

		"ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			DiffSuppressFunc: networkSubresIPAddreDiffSupperss,
			Description: "Optional IP address to assign to this connection. This IP should belong to the selected network and free for use.",
		},

		"mac": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MAC address associated with this connection. MAC address is assigned automatically.",
		},

	}
	return rets
}

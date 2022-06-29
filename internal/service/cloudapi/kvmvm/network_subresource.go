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

package kvmvm

import (
	"bytes"
	"hash/fnv"

	"github.com/rudecs/terraform-provider-decort/internal/provider"
	log "github.com/sirupsen/logrus"

	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

// This function is based on the original Terraform SerializeResourceForHash found
// in helper/schema/serialize.go
// It skips network subresource attributes, which are irrelevant for identification
// of unique network blocks
func networkSubresourceSerialize(output *bytes.Buffer, val interface{}, resource *schema.Resource) {
	if val == nil {
		return
	}

	rs := resource.Schema
	m := val.(map[string]interface{})

	keys := make([]string, 0, len(rs))
	allComputed := true

	for k, val := range rs {
		if val.Optional || val.Required {
			allComputed = false
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		// explicitly ignore "ip_address" when hashing
		if k == "ip_address" {
			continue
		}

		subSchema := rs[k]
		// Skip attributes that are not user-provided. Computed attributes
		// do not contribute to the hash since their ultimate value cannot
		// be known at plan/diff time.
		if !allComputed && !(subSchema.Required || subSchema.Optional) {
			continue
		}

		output.WriteString(k)
		output.WriteRune(':')
		value := m[k]
		schema.SerializeValueForHash(output, value, subSchema)
	}
}

// HashNetworkSubresource hashes network subresource of compute resource. It uses
// specially designed networkSubresourceSerialize (see above) to make sure hashing
// does not involve attributes that we deem irrelevant to the uniqueness of network
// subresource definitions.
// It is this function that should be specified as SchemaSetFunc when creating Set
// from network subresource (e.g. in flattenCompute)
//
// This function is based on the original Terraform function HashResource from
// helper/schema/set.go
func HashNetworkSubresource(resource *schema.Resource) schema.SchemaSetFunc {
	return func(v interface{}) int {
		var serialized bytes.Buffer
		networkSubresourceSerialize(&serialized, v, resource)

		hs := fnv.New32a()
		hs.Write(serialized.Bytes())
		return int(hs.Sum32())
	}
}

func networkSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"net_type": {
			Type:         schema.TypeString,
			Required:     true,
			StateFunc:    provider.StateFuncToUpper,
			ValidateFunc: validation.StringInSlice([]string{"EXTNET", "VINS"}, false), // observe case while validating
			Description:  "Type of the network for this connection, either EXTNET or VINS.",
		},

		"net_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the network for this connection.",
		},

		"ip_address": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			DiffSuppressFunc: networkSubresIPAddreDiffSupperss,
			Description:      "Optional IP address to assign to this connection. This IP should belong to the selected network and free for use.",
		},

		"mac": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MAC address associated with this connection. MAC address is assigned automatically.",
		},
	}
	return rets
}

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

package pfw

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourcePfwCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourcePfwCreate: called for compute %d", d.Get("compute_id").(int))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	urlValues.Add("publicPortStart", strconv.Itoa(d.Get("public_port_start").(int)))
	urlValues.Add("localBasePort", strconv.Itoa(d.Get("local_base_port").(int)))
	urlValues.Add("proto", d.Get("proto").(string))

	if portEnd, ok := d.GetOk("public_port_end"); ok {
		urlValues.Add("publicPortEnd", strconv.Itoa(portEnd.(int)))
	}

	pfwId, err := c.DecortAPICall(ctx, "POST", ComputePfwAddAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d-%s", d.Get("compute_id").(int), pfwId))

	pfw, err := utilityPfwCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("local_ip", pfw.LocalIP)
	if _, ok := d.GetOk("public_port_end"); !ok {
		d.Set("public_port_end", pfw.PublicPortEnd)
	}

	return nil
}

func resourcePfwRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourcePfwRead: called for compute %d, rule %s", d.Get("compute_id").(int), d.Id())

	pfw, err := utilityPfwCheckPresence(ctx, d, m)
	if pfw == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("compute_id", pfw.ComputeID)
	d.Set("public_port_start", pfw.PublicPortStart)
	d.Set("public_port_end", pfw.PublicPortEnd)
	d.Set("local_ip", pfw.LocalIP)
	d.Set("local_base_port", pfw.LocalPort)
	d.Set("proto", pfw.Protocol)

	return nil
}

func resourcePfwDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourcePfwDelete: called for compute %d, rule %s", d.Get("compute_id").(int), d.Id())

	pfw, err := utilityPfwCheckPresence(ctx, d, m)
	if pfw == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	urlValues.Add("ruleId", strconv.Itoa(pfw.ID))

	_, err = c.DecortAPICall(ctx, "POST", ComputePfwDelAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePfwExists(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourcePfwExists: called for compute %d, rule %s", d.Get("compute_id").(int), d.Id())

	pfw, err := utilityPfwCheckPresence(ctx, d, m)
	if pfw == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func resourcePfwSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of compute instance.",
		},

		"public_port_start": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "External start port number for the rule.",
		},

		"public_port_end": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "End port number (inclusive) for the ranged rule.",
		},

		"local_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP address of compute instance.",
		},

		"local_base_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "Internal base port number.",
		},

		"proto": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
			Description:  "Network protocol, either 'tcp' or 'udp'.",
		},
	}
}

func ResourcePfw() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourcePfwCreate,
		ReadContext:   resourcePfwRead,
		DeleteContext: resourcePfwDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout60s,
			Read:    &constants.Timeout30s,
			Update:  &constants.Timeout60s,
			Delete:  &constants.Timeout60s,
			Default: &constants.Timeout60s,
		},

		Schema: resourcePfwSchemaMake(),
	}
}

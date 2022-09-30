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

package lb

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceLBFrontendBindCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBFrontendBindCreate")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("frontendName", d.Get("frontend_name").(string))
	urlValues.Add("bindingName", d.Get("name").(string))
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("bindingAddress", d.Get("address").(string))
	urlValues.Add("bindingPort", strconv.Itoa(d.Get("port").(int)))

	_, err := c.DecortAPICall(ctx, "POST", lbFrontendBindAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("lb_id").(int)) + "#" + d.Get("frontend_name").(string) + "#" + d.Get("name").(string))

	_, err = utilityLBFrontendBindCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceLBFrontendBindRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceLBFrontendBindRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBFrontendBindRead")

	b, err := utilityLBFrontendBindCheckPresence(ctx, d, m)
	if b == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	lbId, _ := strconv.ParseInt(strings.Split(d.Id(), "#")[0], 10, 32)
	frontendName := strings.Split(d.Id(), "#")[1]

	d.Set("lb_id", lbId)
	d.Set("frontend_name", frontendName)
	d.Set("name", b.Name)
	d.Set("address", b.Address)
	d.Set("guid", b.GUID)
	d.Set("port", b.Port)

	return nil
}

func resourceLBFrontendBindDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBFrontendBindDelete")

	b, err := utilityLBFrontendBindCheckPresence(ctx, d, m)
	if b == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("bindingName", d.Get("name").(string))
	urlValues.Add("frontendName", d.Get("frontend_name").(string))

	_, err = c.DecortAPICall(ctx, "POST", lbFrontendBindDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceLBFrontendBindEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBFrontendBindEdit")
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("frontendName", d.Get("frontend_name").(string))
	urlValues.Add("bindingName", d.Get("name").(string))
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))

	if d.HasChange("address") {
		urlValues.Add("bindingAddress", d.Get("address").(string))
	}

	if d.HasChange("port") {
		urlValues.Add("bindingPort", strconv.Itoa(d.Get("port").(int)))
	}

	_, err := c.DecortAPICall(ctx, "POST", lbFrontendBindUpdateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLBFrontendBindRead(ctx, d, m)
}

func ResourceLBFrontendBind() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceLBFrontendBindCreate,
		ReadContext:   resourceLBFrontendBindRead,
		UpdateContext: resourceLBFrontendBindEdit,
		DeleteContext: resourceLBFrontendBindDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout60s,
			Read:    &constants.Timeout30s,
			Update:  &constants.Timeout60s,
			Delete:  &constants.Timeout60s,
			Default: &constants.Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the LB instance to backendCreate",
			},
			"frontend_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Must be unique among all backends of this LB - name of the new backend to create",
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

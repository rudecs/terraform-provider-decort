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

func resourceLBFrontendCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBFrontendCreate")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("backendName", d.Get("backend_name").(string))
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("frontendName", d.Get("name").(string))

	_, err := c.DecortAPICall(ctx, "POST", lbFrontendCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("lb_id").(int)) + "#" + d.Get("name").(string))

	_, err = utilityLBFrontendCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceLBFrontendRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceLBFrontendRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBFrontendRead")

	f, err := utilityLBFrontendCheckPresence(ctx, d, m)
	if f == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	lbId, _ := strconv.ParseInt(strings.Split(d.Id(), "#")[0], 10, 32)
	d.Set("lb_id", lbId)
	d.Set("backend_name", f.Backend)
	d.Set("name", f.Name)
	d.Set("guid", f.GUID)
	d.Set("bindings", flattendBindings(f.Bindings))

	return nil
}

func resourceLBFrontendDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBFrontendDelete")

	lb, err := utilityLBFrontendCheckPresence(ctx, d, m)
	if lb == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("frontendName", d.Get("name").(string))

	_, err = c.DecortAPICall(ctx, "POST", lbFrontendDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceLBFrontendEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	//TODO: перенести bindings сюда

	return nil
}

func ResourceLBFrontend() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceLBFrontendCreate,
		ReadContext:   resourceLBFrontendRead,
		UpdateContext: resourceLBFrontendEdit,
		DeleteContext: resourceLBFrontendDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout600s,
			Read:    &constants.Timeout300s,
			Update:  &constants.Timeout300s,
			Delete:  &constants.Timeout300s,
			Default: &constants.Timeout300s,
		},

		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the LB instance to backendCreate",
			},
			"backend_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bindings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

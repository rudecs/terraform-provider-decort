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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceLBBackendServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendServerCreate")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("backendName", d.Get("backend_name").(string))
	urlValues.Add("serverName", d.Get("name").(string))
	urlValues.Add("address", d.Get("address").(string))
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("port", strconv.Itoa(d.Get("port").(int)))

	if check, ok := d.GetOk("check"); ok {
		urlValues.Add("check", check.(string))
	}

	if inter, ok := d.GetOk("inter"); ok {
		urlValues.Add("inter", strconv.Itoa(inter.(int)))
	}
	if downinter, ok := d.GetOk("downinter"); ok {
		urlValues.Add("downinter", strconv.Itoa(downinter.(int)))
	}
	if rise, ok := d.GetOk("rise"); ok {
		urlValues.Add("rise", strconv.Itoa(rise.(int)))
	}
	if fall, ok := d.GetOk("fall"); ok {
		urlValues.Add("fall", strconv.Itoa(fall.(int)))
	}
	if slowstart, ok := d.GetOk("slowstart"); ok {
		urlValues.Add("slowstart", strconv.Itoa(slowstart.(int)))
	}
	if maxconn, ok := d.GetOk("maxconn"); ok {
		urlValues.Add("maxconn", strconv.Itoa(maxconn.(int)))
	}
	if maxqueue, ok := d.GetOk("maxqueue"); ok {
		urlValues.Add("maxqueue", strconv.Itoa(maxqueue.(int)))
	}
	if weight, ok := d.GetOk("weight"); ok {
		urlValues.Add("weight", strconv.Itoa(weight.(int)))
	}

	_, err := c.DecortAPICall(ctx, "POST", lbBackendServerAddAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("lb_id").(int)) + "-" + d.Get("backend_name").(string) + "-" + d.Get("name").(string))

	_, err = utilityLBBackendServerCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceLBBackendServerRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceLBBackendServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendServerRead")

	s, err := utilityLBBackendServerCheckPresence(ctx, d, m)
	if s == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	lbId, _ := strconv.ParseInt(strings.Split(d.Id(), "-")[0], 10, 32)
	backendName := strings.Split(d.Id(), "-")[1]

	d.Set("lb_id", lbId)
	d.Set("backend_name", backendName)
	d.Set("name", s.Name)
	d.Set("port", s.Port)
	d.Set("address", s.Address)
	d.Set("check", s.Check)
	d.Set("guid", s.GUID)
	d.Set("downinter", s.ServerSettings.DownInter)
	d.Set("fall", s.ServerSettings.Fall)
	d.Set("inter", s.ServerSettings.Inter)
	d.Set("maxconn", s.ServerSettings.MaxConn)
	d.Set("maxqueue", s.ServerSettings.MaxQueue)
	d.Set("rise", s.ServerSettings.Rise)
	d.Set("slowstart", s.ServerSettings.SlowStart)
	d.Set("weight", s.ServerSettings.Weight)

	return nil
}

func resourceLBBackendServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendServerDelete")

	lb, err := utilityLBBackendServerCheckPresence(ctx, d, m)
	if lb == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("serverName", d.Get("name").(string))
	urlValues.Add("backendName", d.Get("backend_name").(string))

	_, err = c.DecortAPICall(ctx, "POST", lbBackendServerDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceLBBackendServerEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendServerEdit")
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("backendName", d.Get("backend_name").(string))
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("serverName", d.Get("name").(string))
	urlValues.Add("address", d.Get("address").(string))
	urlValues.Add("port", strconv.Itoa(d.Get("port").(int)))

	if d.HasChange("check") {
		urlValues.Add("check", d.Get("check").(string))
	}
	if d.HasChange("inter") {
		urlValues.Add("inter", strconv.Itoa(d.Get("inter").(int)))
	}
	if d.HasChange("downinter") {
		urlValues.Add("downinter", strconv.Itoa(d.Get("downinter").(int)))
	}
	if d.HasChange("rise") {
		urlValues.Add("rise", strconv.Itoa(d.Get("rise").(int)))
	}
	if d.HasChange("fall") {
		urlValues.Add("fall", strconv.Itoa(d.Get("fall").(int)))
	}
	if d.HasChange("slowstart") {
		urlValues.Add("slowstart", strconv.Itoa(d.Get("slowstart").(int)))
	}
	if d.HasChange("maxconn") {
		urlValues.Add("maxconn", strconv.Itoa(d.Get("maxconn").(int)))
	}
	if d.HasChange("maxqueue") {
		urlValues.Add("maxqueue", strconv.Itoa(d.Get("maxqueue").(int)))
	}
	if d.HasChange("weight") {
		urlValues.Add("weight", strconv.Itoa(d.Get("weight").(int)))
	}

	_, err := c.DecortAPICall(ctx, "POST", lbBackendServerUpdateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	//TODO: перенести servers сюда

	return resourceLBBackendServerRead(ctx, d, m)
}

func ResourceLBBackendServer() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceLBBackendServerCreate,
		ReadContext:   resourceLBBackendServerRead,
		UpdateContext: resourceLBBackendServerEdit,
		DeleteContext: resourceLBBackendServerDelete,

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
			"backend_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Must be unique among all backends of this LB - name of the new backend to create",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Must be unique among all servers defined for this backend - name of the server definition to add.",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address of the server.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port number on the server",
			},
			"check": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
				Description:  "set to disabled if this server should be used regardless of its state.",
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"downinter": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"fall": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"inter": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"maxconn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"maxqueue": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"rise": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"slowstart": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

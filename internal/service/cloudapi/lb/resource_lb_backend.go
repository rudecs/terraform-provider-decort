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

func resourceLBBackendCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendCreate")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("backendName", d.Get("name").(string))
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))

	if algorithm, ok := d.GetOk("algorithm"); ok {
		urlValues.Add("algorithm", algorithm.(string))
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

	_, err := c.DecortAPICall(ctx, "POST", lbBackendCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("lb_id").(int)) + "#" + d.Get("name").(string))

	_, err = utilityLBBackendCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceLBBackendRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceLBBackendRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendRead")

	b, err := utilityLBBackendCheckPresence(ctx, d, m)
	if b == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	lbId, _ := strconv.ParseInt(strings.Split(d.Id(), "#")[0], 10, 32)

	d.Set("lb_id", lbId)
	d.Set("name", b.Name)
	d.Set("algorithm", b.Algorithm)
	d.Set("guid", b.GUID)
	d.Set("downinter", b.ServerDefaultSettings.DownInter)
	d.Set("fall", b.ServerDefaultSettings.Fall)
	d.Set("inter", b.ServerDefaultSettings.Inter)
	d.Set("maxconn", b.ServerDefaultSettings.MaxConn)
	d.Set("maxqueue", b.ServerDefaultSettings.MaxQueue)
	d.Set("rise", b.ServerDefaultSettings.Rise)
	d.Set("slowstart", b.ServerDefaultSettings.SlowStart)
	d.Set("weight", b.ServerDefaultSettings.Weight)
	d.Set("servers", flattenServers(b.Servers))

	return nil
}

func resourceLBBackendDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendDelete")

	lb, err := utilityLBBackendCheckPresence(ctx, d, m)
	if lb == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
	urlValues.Add("backendName", d.Get("name").(string))

	_, err = c.DecortAPICall(ctx, "POST", lbBackendDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceLBBackendEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBBackendEdit")
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("backendName", d.Get("name").(string))
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))

	if d.HasChange("algorithm") {
		urlValues.Add("algorithm", d.Get("algorithm").(string))
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

	_, err := c.DecortAPICall(ctx, "POST", lbBackendUpdateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	//TODO: перенести servers сюда

	return resourceLBBackendRead(ctx, d, m)
}

func ResourceLBBackend() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceLBBackendCreate,
		ReadContext:   resourceLBBackendRead,
		UpdateContext: resourceLBBackendEdit,
		DeleteContext: resourceLBBackendDelete,

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Must be unique among all backends of this LB - name of the new backend to create",
			},
			"algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"roundrobin", "static-rr", "leastconn"}, false),
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
			"servers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"check": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"server_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"guid": {
										Type:     schema.TypeString,
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
							},
						},
					},
				},
			},
		},
	}
}

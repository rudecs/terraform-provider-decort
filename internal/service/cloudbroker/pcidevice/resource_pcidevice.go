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

package pcidevice

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/flattens"
	log "github.com/sirupsen/logrus"
)

func resourcePcideviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourcePcideviceCreate: called for pcidevice %s", d.Get("name").(string))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("hwPath", d.Get("hw_path").(string))
	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("stackId", strconv.Itoa(d.Get("stack_id").(int)))

	if description, ok := d.GetOk("description"); ok {
		urlValues.Add("description", description.(string))
	}

	pcideviceId, err := c.DecortAPICall(ctx, "POST", pcideviceCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pcideviceId)
	d.Set("device_id", pcideviceId)

	diagnostics := resourcePcideviceRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourcePcideviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pcidevice, err := utilityPcideviceCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(pcidevice.ID))
	d.Set("ckey", pcidevice.CKey)
	d.Set("meta", flattens.FlattenMeta(pcidevice.Meta))
	d.Set("compute_id", pcidevice.Computeid)
	d.Set("description", pcidevice.Description)
	d.Set("guid", pcidevice.Guid)
	d.Set("hw_path", pcidevice.HwPath)
	d.Set("device_id", pcidevice.ID)
	d.Set("rg_id", pcidevice.RgID)
	d.Set("name", pcidevice.Name)
	d.Set("stack_id", pcidevice.StackID)
	d.Set("status", pcidevice.Status)
	d.Set("system_name", pcidevice.SystemName)

	return nil
}

func resourcePcideviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourcePcideviceDelete: called for %s, id: %s", d.Get("name").(string), d.Id())

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("deviceId", d.Id())
	urlValues.Add("force", strconv.FormatBool(d.Get("force").(bool)))

	_, err := c.DecortAPICall(ctx, "POST", pcideviceDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourcePcideviceExists(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	pcidevice, err := utilityPcideviceCheckPresence(ctx, d, m)
	if err != nil {
		return false, err
	}
	if pcidevice == nil {
		return false, nil
	}

	return true, nil
}

func resourcePcideviceEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChange("enable") {
		state := d.Get("enable").(bool)
		c := m.(*controller.ControllerCfg)
		urlValues := &url.Values{}
		api := ""

		urlValues.Add("deviceId", strconv.Itoa(d.Get("device_id").(int)))

		if state {
			api = pcideviceEnableAPI
		} else {
			api = pcideviceDisableAPI
		}

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	diagnostics := resourcePcideviceRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourcePcideviceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"meta": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"compute_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "description, just for information",
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"hw_path": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "PCI address of the device",
		},
		"device_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of Device",
		},
		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Resource GROUP",
		},
		"stack_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "stackId",
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"force": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force delete",
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable pci device",
		},
	}
}

func ResourcePcidevice() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourcePcideviceCreate,
		ReadContext:   resourcePcideviceRead,
		UpdateContext: resourcePcideviceEdit,
		DeleteContext: resourcePcideviceDelete,

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

		Schema: resourcePcideviceSchemaMake(),
	}
}

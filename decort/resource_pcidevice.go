/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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
This file is part of Terraform (by Hashicorp) provider for Digital Energy Cloud Orchestration
Technology platfom.

Visit https://github.com/rudecs/terraform-provider-decort for full source code package and updates.
*/

package decort

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourcePcideviceCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourcePcideviceCreate: called for pcidevice %s", d.Get("name").(string))

	if deviceId, ok := d.GetOk("device_id"); ok {
		if exists, err := resourcePcideviceExists(d, m); exists {
			if err != nil {
				return err
			}
			d.SetId(strconv.Itoa(deviceId.(int)))
			err = resourcePcideviceRead(d, m)
			if err != nil {
				return err
			}

			return nil
		}
		return errors.New("provided device id does not exist")
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("hwPath", d.Get("hw_path").(string))
	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("stackId", strconv.Itoa(d.Get("stack_id").(int)))

	if description, ok := d.GetOk("description"); ok {
		urlValues.Add("description", description.(string))
	}

	pcideviceId, err := controller.decortAPICall("POST", pcideviceCreateAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId(pcideviceId)
	d.Set("device_id", pcideviceId)

	err = resourcePcideviceRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourcePcideviceRead(d *schema.ResourceData, m interface{}) error {
	pcidevice, err := utilityPcideviceCheckPresence(d, m)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(pcidevice.ID))
	d.Set("ckey", pcidevice.CKey)
	d.Set("meta", flattenMeta(pcidevice.Meta))
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

func resourcePcideviceDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourcePcideviceDelete: called for %s, id: %s", d.Get("name").(string), d.Id())

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("deviceId", d.Id())
	urlValues.Add("force", strconv.FormatBool(d.Get("force").(bool)))

	_, err := controller.decortAPICall("POST", pcideviceDeleteAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourcePcideviceExists(d *schema.ResourceData, m interface{}) (bool, error) {
	pcidevice, err := utilityPcideviceCheckPresence(d, m)
	if err != nil {
		return false, err
	}
	if pcidevice == nil {
		return false, nil
	}

	return true, nil
}

func resourcePcideviceEdit(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("enable") {
		state := d.Get("enable").(bool)
		c := m.(*ControllerCfg)
		urlValues := &url.Values{}
		api := ""

		urlValues.Add("deviceId", strconv.Itoa(d.Get("device_id").(int)))

		if state {
			api = pcideviceEnableAPI
		} else {
			api = pcideviceDisableAPI
		}

		_, err := c.decortAPICall("POST", api, urlValues)
		if err != nil {
			return err
		}
	}

	err := resourcePcideviceRead(d, m)
	if err != nil {
		return err
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

func resourcePcidevice() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourcePcideviceCreate,
		Read:   resourcePcideviceRead,
		Update: resourcePcideviceEdit,
		Delete: resourcePcideviceDelete,
		Exists: resourcePcideviceExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &Timeout60s,
			Read:    &Timeout30s,
			Update:  &Timeout60s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

		Schema: resourcePcideviceSchemaMake(),
	}
}

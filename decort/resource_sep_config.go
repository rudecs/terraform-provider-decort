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
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceSepConfigCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepConfigCreate: called for sep id %d", d.Get("sep_id").(int))

	if _, ok := d.GetOk("sep_id"); ok {
		if exists, err := resourceSepConfigExists(d, m); exists {
			if err != nil {
				return err
			}
			id := uuid.New()
			d.SetId(id.String())
			err = resourceSepConfigRead(d, m)
			if err != nil {
				return err
			}

			return nil
		}
		return errors.New("provided sep id config does not exist")
	}

	return resourceSepConfigRead(d, m)
}

func resourceSepConfigRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepConfigRead: called for sep id: %d", d.Get("sep_id").(int))

	sepConfig, err := utilitySepConfigCheckPresence(d, m)
	if sepConfig == nil {
		d.SetId("")
		return err
	}
	data, _ := json.Marshal(sepConfig)
	d.Set("config", string(data))
	return nil
}

func resourceSepConfigDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func resourceSepConfigExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourceSepConfigExists: called for sep id: %d", d.Get("sep_id").(int))

	sepDesConfig, err := utilitySepConfigCheckPresence(d, m)
	if sepDesConfig == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func resourceSepConfigEdit(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepConfigEdit: called for sep id: %d", d.Get("sep_id").(int))
	c := m.(*ControllerCfg)

	urlValues := &url.Values{}
	if d.HasChange("config") {
		urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
		urlValues.Add("config", d.Get("config").(string))
		_, err := c.decortAPICall("POST", sepConfigValidateAPI, urlValues)
		if err != nil {
			return err
		}
		_, err = c.decortAPICall("POST", sepConfigInsertAPI, urlValues)
		if err != nil {
			return err
		}

	}

	urlValues = &url.Values{}
	if d.HasChange("field_edit") {
		fieldConfig := d.Get("field_edit").([]interface{})
		field := fieldConfig[0].(map[string]interface{})
		urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
		urlValues.Add("field_name", field["field_name"].(string))
		urlValues.Add("field_value", field["field_value"].(string))
		urlValues.Add("field_type", field["field_type"].(string))

		_, err := c.decortAPICall("POST", sepConfigFieldEditAPI, urlValues)
		if err != nil {
			return err
		}
	}

	err := resourceSepConfigRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceSepConfigSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sep_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"config": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"field_edit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"field_value": {
						Type:     schema.TypeString,
						Required: true,
					},
					"field_type": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func resourceSepConfig() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceSepConfigCreate,
		Read:   resourceSepConfigRead,
		Update: resourceSepConfigEdit,
		Delete: resourceSepConfigDelete,
		Exists: resourceSepConfigExists,

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

		Schema: resourceSepConfigSchemaMake(),
	}
}

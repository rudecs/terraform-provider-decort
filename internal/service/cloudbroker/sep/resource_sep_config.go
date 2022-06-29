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

package sep

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceSepConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSepConfigCreate: called for sep id %d", d.Get("sep_id").(int))

	if _, ok := d.GetOk("sep_id"); ok {
		if exists, err := resourceSepConfigExists(d, m); exists {
			if err != nil {
				return diag.FromErr(err)
			}
			id := uuid.New()
			d.SetId(id.String())
			diagnostics := resourceSepConfigRead(ctx, d, m)
			if diagnostics != nil {
				return diagnostics
			}

			return nil
		}
		return diag.Errorf("provided sep id config does not exist")
	}

	return resourceSepConfigRead(ctx, d, m)
}

func resourceSepConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSepConfigRead: called for sep id: %d", d.Get("sep_id").(int))

	sepConfig, err := utilitySepConfigCheckPresence(d, m)
	if sepConfig == nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	data, _ := json.Marshal(sepConfig)
	d.Set("config", string(data))
	return nil
}

func resourceSepConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceSepConfigEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSepConfigEdit: called for sep id: %d", d.Get("sep_id").(int))
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	if d.HasChange("config") {
		urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
		urlValues.Add("config", d.Get("config").(string))
		_, err := c.DecortAPICall("POST", sepConfigValidateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall("POST", sepConfigInsertAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
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

		_, err := c.DecortAPICall("POST", sepConfigFieldEditAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	diagnostics := resourceSepConfigRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
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

func ResourceSepConfig() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceSepConfigCreate,
		ReadContext:   resourceSepConfigRead,
		UpdateContext: resourceSepConfigEdit,
		DeleteContext: resourceSepConfigDelete,
		Exists:        resourceSepConfigExists,

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

		Schema: resourceSepConfigSchemaMake(),
	}
}

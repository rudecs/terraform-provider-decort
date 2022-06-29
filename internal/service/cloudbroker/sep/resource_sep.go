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
	"errors"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/flattens"
	log "github.com/sirupsen/logrus"
)

func resourceSepCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSepCreate: called for sep %s", d.Get("name").(string))

	if sepId, ok := d.GetOk("sep_id"); ok {
		if exists, err := resourceSepExists(ctx, d, m); exists {
			if err != nil {
				return diag.FromErr(err)
			}
			d.SetId(strconv.Itoa(sepId.(int)))
			diagnostics := resourceSepRead(ctx, d, m)
			if diagnostics != nil {
				return diagnostics
			}

			return nil
		}
		return diag.Errorf("provided sep id does not exist")
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))
	urlValues.Add("sep_type", d.Get("type").(string))

	if desc, ok := d.GetOk("desc"); ok {
		urlValues.Add("description", desc.(string))
	}
	if configString, ok := d.GetOk("config"); ok {
		urlValues.Add("config", configString.(string))
	}
	if enable, ok := d.GetOk("enable"); ok {
		urlValues.Add("enable", strconv.FormatBool(enable.(bool)))
	}

	tstr := d.Get("consumed_by").([]interface{})
	temp := ""
	l := len(tstr)
	for i, str := range tstr {
		s := "\"" + str.(string) + "\""
		if i != (l - 1) {
			s += ","
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("consumer_nids", temp)

	tstr = d.Get("provided_by").([]interface{})
	temp = ""
	l = len(tstr)
	for i, str := range tstr {
		s := "\"" + str.(string) + "\""
		if i != (l - 1) {
			s += ","
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("provider_nids", temp)

	sepId, err := c.DecortAPICall(ctx, "POST", sepCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(sepId)
	d.Set("sep_id", sepId)

	diagnostics := resourceSepRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	d.SetId(id.String())

	return nil
}

func resourceSepRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSepRead: called for %s id: %d", d.Get("name").(string), d.Get("sep_id").(int))

	sep, err := utilitySepCheckPresence(ctx, d, m)
	if sep == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("ckey", sep.Ckey)
	d.Set("meta", flattens.FlattenMeta(sep.Meta))
	d.Set("consumed_by", sep.ConsumedBy)
	d.Set("desc", sep.Desc)
	d.Set("gid", sep.Gid)
	d.Set("guid", sep.Guid)
	d.Set("sep_id", sep.Id)
	d.Set("milestones", sep.Milestones)
	d.Set("name", sep.Name)
	d.Set("obj_status", sep.ObjStatus)
	d.Set("provided_by", sep.ProvidedBy)
	d.Set("tech_status", sep.TechStatus)
	d.Set("type", sep.Type)
	data, _ := json.Marshal(sep.Config)
	d.Set("config", string(data))

	return nil
}

func resourceSepDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSepDelete: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))

	sepDes, err := utilitySepCheckPresence(ctx, d, m)
	if sepDes == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))

	_, err = c.DecortAPICall(ctx, "POST", sepDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceSepExists(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourceSepExists: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))

	sepDes, err := utilitySepCheckPresence(ctx, d, m)
	if sepDes == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func resourceSepEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSepEdit: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	if d.HasChange("decommission") {
		decommission := d.Get("decommission").(bool)
		if decommission {
			urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
			urlValues.Add("clear_physically", strconv.FormatBool(d.Get("clear_physically").(bool)))
			_, err := c.DecortAPICall(ctx, "POST", sepDecommissionAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("upd_capacity_limit") {
		updCapacityLimit := d.Get("upd_capacity_limit").(bool)
		if updCapacityLimit {
			urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
			_, err := c.DecortAPICall(ctx, "POST", sepUpdateCapacityLimitAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("config") {
		urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
		urlValues.Add("config", d.Get("config").(string))
		_, err := c.DecortAPICall(ctx, "POST", sepConfigValidateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall(ctx, "POST", sepConfigInsertAPI, urlValues)
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

		_, err := c.DecortAPICall(ctx, "POST", sepConfigFieldEditAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("enable") {
		err := resourceSepChangeEnabled(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("consumed_by") {
		err := resourceSepUpdateNodes(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("provided_by") {
		err := resourceSepUpdateProviders(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	urlValues = &url.Values{}
	if diagnostics := resourceSepRead(ctx, d, m); diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceSepChangeEnabled(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	var api string

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
	if d.Get("enable").(bool) {
		api = sepEnableAPI
	} else {
		api = sepDisableAPI
	}
	resp, err := c.DecortAPICall(ctx, "POST", api, urlValues)
	if err != nil {
		return err
	}
	res, err := strconv.ParseBool(resp)
	if err != nil {
		return err
	}
	if !res {
		return errors.New("Cannot enable/disable")
	}
	return nil
}

func resourceSepUpdateNodes(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepUpdateNodes: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	t1, t2 := d.GetChange("consumed_by")

	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))

	consumedIds := make([]interface{}, 0)
	temp := ""
	api := ""

	if d1, d2 := t1.([]interface{}), t2.([]interface{}); len(d1) > len(d2) {
		for _, n := range d2 {
			if !findElInt(d1, n) {
				consumedIds = append(consumedIds, n)
			}
		}
		api = sepDelConsumerNodesAPI
	} else {
		consumedIds = d.Get("consumed_by").([]interface{})
		api = sepAddConsumerNodesAPI
	}

	l := len(consumedIds)
	for i, consumedId := range consumedIds {
		s := strconv.Itoa(consumedId.(int))
		if i != (l - 1) {
			s += ","
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("consumer_nids", temp)
	_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func findElInt(sl []interface{}, el interface{}) bool {
	for _, e := range sl {
		if e.(int) == el.(int) {
			return true
		}
	}
	return false
}

func resourceSepUpdateProviders(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepUpdateProviders: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
	providerIds := d.Get("provided_by").([]interface{})
	temp := ""
	l := len(providerIds)
	for i, providerId := range providerIds {
		s := strconv.Itoa(providerId.(int))
		if i != (l - 1) {
			s += ","
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("provider_nids", temp)
	_, err := c.DecortAPICall(ctx, "POST", sepAddProviderNodesAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceSepSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "sep type des id",
		},
		"upd_capacity_limit": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Update SEP capacity limit",
		},
		"decommission": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "unlink everything that exists from SEP",
		},
		"clear_physically": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "clear disks and images physically",
		},
		"config": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "sep config string",
		},
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
		"consumed_by": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "list of consumer nodes IDs",
		},
		"desc": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "sep description",
		},
		"gid": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "grid (platform) ID",
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SEP name",
		},
		"obj_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provided_by": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "list of provider nodes IDs",
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "type of storage",
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "enable SEP after creation",
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

func ResourceSep() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceSepCreate,
		ReadContext:   resourceSepRead,
		UpdateContext: resourceSepEdit,
		DeleteContext: resourceSepDelete,

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

		Schema: resourceSepSchemaMake(),
	}
}

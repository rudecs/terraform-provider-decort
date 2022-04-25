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
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceSepDesCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepDesCreate: called for sep %s type \"des\"", d.Get("name").(string))

	if sepId, ok := d.GetOk("sep_id"); ok {
		if exists, err := resourceSepDesExists(d, m); exists {
			if err != nil {
				return err
			}
			d.SetId(strconv.Itoa(sepId.(int)))
			err = resourceSepDesRead(d, m)
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
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))
	urlValues.Add("sep_type", "des")

	if desc, ok := d.GetOk("desc"); ok {
		urlValues.Add("description", desc.(string))
	}
	if configString, ok := d.GetOk("config_string"); ok {
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

	sepDesId, err := controller.decortAPICall("POST", sepCreateAPI, urlValues)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(sepDesId)
	d.Set("sep_id", sepDesId)

	err = resourceSepDesRead(d, m)
	if err != nil {
		return err
	}

	d.SetId(id.String())

	return nil
}

func resourceSepDesRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepDesRead: called for %s id: %d", d.Get("name").(string), d.Get("sep_id").(int))

	sepDes, err := utilitySepDesCheckPresence(d, m)
	if sepDes == nil {
		d.SetId("")
		return err
	}

	d.Set("ckey", sepDes.Ckey)
	d.Set("meta", flattenMeta(sepDes.Meta))
	d.Set("consumed_by", sepDes.ConsumedBy)
	d.Set("desc", sepDes.Desc)
	d.Set("gid", sepDes.Gid)
	d.Set("guid", sepDes.Guid)
	d.Set("sep_id", sepDes.Id)
	d.Set("milestones", sepDes.Milestones)
	d.Set("name", sepDes.Name)
	d.Set("obj_status", sepDes.ObjStatus)
	d.Set("provided_by", sepDes.ProvidedBy)
	d.Set("tech_status", sepDes.TechStatus)
	d.Set("type", sepDes.Type)
	data, _ := json.Marshal(sepDes.Config)
	d.Set("config_string", string(data))
	d.Set("config", flattenSepDesConfig(sepDes.Config))

	return nil
}

func resourceSepDesDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepDesDelete: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))

	sepDes, err := utilitySepDesCheckPresence(d, m)
	if sepDes == nil {
		if err != nil {
			return err
		}
		return nil
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))

	_, err = controller.decortAPICall("POST", sepDeleteAPI, urlValues)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceSepDesExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourceSepDesExists: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))

	sepDes, err := utilitySepDesCheckPresence(d, m)
	if sepDes == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func resourceSepDesEdit(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSepDesEdit: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}

	if d.HasChange("decommision") {
		decommision := d.Get("decommision").(bool)
		if decommision {
			urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
			urlValues.Add("clear_physically", strconv.FormatBool(d.Get("clear_physically").(bool)))
			_, err := c.decortAPICall("POST", sepDecommissionAPI, urlValues)
			if err != nil {
				return err
			}
		}
	}
	urlValues = &url.Values{}
	if d.HasChange("upd_capacity_limit") {
		updCapacityLimit := d.Get("upd_capacity_limit").(bool)
		if updCapacityLimit {
			urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
			_, err := c.decortAPICall("POST", sepUpdateCapacityLimitAPI, urlValues)
			if err != nil {
				return err
			}
		}
	}

	urlValues = &url.Values{}

	err := resourceSepDesRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceSepChangeEnabled(d *schema.ResourceDiff, m interface{}) error {
	var api string

	c := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
	if d.Get("enable").(bool) {
		api = sepEnableAPI
	} else {
		api = sepDisableAPI
	}
	resp, err := c.decortAPICall("POST", api, urlValues)
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

func resourceSepUpdateNodes(d *schema.ResourceDiff, m interface{}) error {
	log.Debugf("resourceSepUpdateNodes: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}

	t1, t2 := d.GetChange("consumed_by")
	d1 := t1.([]interface{})
	d2 := t2.([]interface{})

	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))

	consumedIds := make([]interface{}, 0)
	temp := ""
	api := ""

	if len(d1) > len(d2) {
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
	_, err := c.decortAPICall("POST", api, urlValues)
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

func resourceSepUpdateProviders(d *schema.ResourceDiff, m interface{}) error {
	log.Debugf("resourceSepUpdateProviders: called for %s, id: %d", d.Get("name").(string), d.Get("sep_id").(int))
	c := m.(*ControllerCfg)
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
	_, err := c.decortAPICall("POST", sepAddProviderNodesAPI, urlValues)
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
		"decommision": {
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
		"config_string": {
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
			Computed:    true,
			Description: "type of storage",
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "enable SEP after creation",
		},
	}
}

func resourceSepDesSchemaMake(sh map[string]*schema.Schema) map[string]*schema.Schema {
	sh["config"] = &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"api_ips": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"capacity_limit": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"protocol": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"decs3o_app_secret": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"format": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"edgeuser_password": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"edgeuser_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"decs3o_app_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"transport": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"pools": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"types": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"reference_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"pagecache_ratio": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"uris": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ip": {
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
						},
					},
				},
				"housekeeping_settings": {
					Type:     schema.TypeList,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"disk_del_queue": {
								Type:     schema.TypeList,
								Computed: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"chunk_max_size": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"disk_count_max": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"enabled": {
											Type:     schema.TypeBool,
											Computed: true,
										},
										"normal_time_to_sleep": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"one_minute_la_threshold": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"oversize_time_to_sleep": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"purge_attempts_threshold": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"purgatory_id": {
											Type:     schema.TypeInt,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return sh
}

func resourceSepDes() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceSepDesCreate,
		Read:   resourceSepDesRead,
		Update: resourceSepDesEdit,
		Delete: resourceSepDesDelete,
		Exists: resourceSepDesExists,

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

		Schema: resourceSepDesSchemaMake(resourceSepSchemaMake()),

		CustomizeDiff: customdiff.All(
			customdiff.IfValueChange("enable", func(old, new, meta interface{}) bool {
				if old.(bool) != new.(bool) {
					return true
				}
				return false
			}, resourceSepChangeEnabled),
			customdiff.IfValueChange("consumed_by", func(old, new, meta interface{}) bool {
				o := old.([]interface{})
				n := new.([]interface{})

				if len(o) != len(n) {
					return true
				} else if len(o) == 0 {
					return false
				}
				count := 0
				for i, v := range n {
					if v.(int) == o[i].(int) {
						count++
					}
				}
				if count == 0 {
					return true
				}
				return false
			}, resourceSepUpdateNodes),
			customdiff.IfValueChange("provided_by", func(old, new, meta interface{}) bool {
				o := old.([]interface{})
				n := new.([]interface{})

				if len(o) != len(n) {
					return true
				} else if len(o) == 0 {
					return false
				}
				count := 0
				for i, v := range n {
					if v.(int) == o[i].(int) {
						count++
					}
				}
				if count == 0 {
					return true
				}
				return false
			}, resourceSepUpdateProviders),
		),
	}
}

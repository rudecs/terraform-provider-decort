/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>, <svs1370@gmail.com>

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

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSepTatlinConfigRead(d *schema.ResourceData, m interface{}) error {
	tatlinSepConfig, err := utilitySepTatlinConfigCheckPresence(d, m)
	if err != nil {
		return err
	}
	id := uuid.New()
	d.SetId(id.String())

	data, _ := json.Marshal(tatlinSepConfig)
	d.Set("config_string", string(data))
	d.Set("config", flattenSepTatlinConfig(*tatlinSepConfig))

	return nil
}

func dataSourceSepTatlinConfigSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sep_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "storage endpoint provider ID",
		},
		"config_string": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"api_urls": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"disk_max_size": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"ovs_settings": {
						Type:     schema.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"vpool_data_metadatacache": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"vpool_vmstor_metadatacache": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"edgeuser_password": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"edgeuser_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"format": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"host_group_name": {
						Type:     schema.TypeString,
						Computed: true,
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
											"purge_attempts_threshold": {
												Type:     schema.TypeInt,
												Computed: true,
											},
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
					"mgmt_password": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"mgmt_user": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name_prefix": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pools": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"types": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"usage_limit": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"ports": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ips": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"iqn": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"protocol": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tech_disk": {
						Type:     schema.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"pool": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"size": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"wwid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSepTatlinConfig() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepTatlinConfigRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepTatlinConfigSchemaMake(),
	}
}

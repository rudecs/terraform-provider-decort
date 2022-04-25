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

func dataSourceSepDesRead(d *schema.ResourceData, m interface{}) error {
	desSep, err := utilitySepDesCheckPresence(d, m)
	if err != nil {
		return err
	}
	id := uuid.New()
	d.SetId(id.String())

	d.Set("ckey", desSep.Ckey)

	d.Set("meta", flattenMeta(desSep.Meta))
	d.Set("consumed_by", desSep.ConsumedBy)
	d.Set("desc", desSep.Desc)
	d.Set("gid", desSep.Gid)
	d.Set("guid", desSep.Guid)
	d.Set("sep_id", desSep.Id)
	d.Set("milestones", desSep.Milestones)
	d.Set("name", desSep.Name)
	d.Set("obj_status", desSep.ObjStatus)
	d.Set("provided_by", desSep.ProvidedBy)
	d.Set("tech_status", desSep.TechStatus)
	d.Set("type", desSep.Type)
	data, _ := json.Marshal(desSep.Config)
	d.Set("config_string", string(data))
	err = d.Set("config", flattenSepDesConfig(desSep.Config))
	if err != nil {
		return err
	}

	return nil
}

func flattenSepDesConfig(sc DesConfigSep) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"api_ips":           sc.ApiIps,
			"protocol":          sc.Protocol,
			"capacity_limit":    sc.CapacityLimit,
			"decs3o_app_secret": sc.Decs3oAppSecret,
			"decs3o_app_id":     sc.Decs3oAppId,
			"format":            sc.Format,
			"edgeuser_password": sc.EdgeuserPassword,
			"edgeuser_name":     sc.EdgeuserName,
			"transport":         sc.Transport,
			"ovs_settings": []interface{}{
				map[string]interface{}{
					"vpool_data_metadatacache":   sc.OVSSettings.VPoolDataMetadataCache,
					"vpool_vmstor_metadatacache": sc.OVSSettings.VPoolVMstorMetadataCache,
				},
			},
			"housekeeping_settings": []interface{}{
				map[string]interface{}{
					"disk_del_queue": []interface{}{
						map[string]interface{}{
							"purgatory_id":             sc.HousekeepingSettings.DiskDelQueue.PurgatoryId,
							"chunk_max_size":           sc.HousekeepingSettings.DiskDelQueue.ChunkMaxSize,
							"disk_count_max":           sc.HousekeepingSettings.DiskDelQueue.DiskCountMax,
							"enabled":                  sc.HousekeepingSettings.DiskDelQueue.Enabled,
							"normal_time_to_sleep":     sc.HousekeepingSettings.DiskDelQueue.NormalTimeToSleep,
							"one_minute_la_threshold":  sc.HousekeepingSettings.DiskDelQueue.OneMinuteLaThreshold,
							"oversize_time_to_sleep":   sc.HousekeepingSettings.DiskDelQueue.OversizeTimeToSleep,
							"purge_attempts_threshold": sc.HousekeepingSettings.DiskDelQueue.PurgeAttemptsThreshold,
						},
					},
				},
			},
			"pools": flattenSepDesPools(sc.Pools),
		},
	}
}

func flattenSepDesPools(pools DesConfigPoolList) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, pool := range pools {
		t := map[string]interface{}{
			"types":           pool.Types,
			"reference_id":    pool.ReferenceId,
			"name":            pool.Name,
			"pagecache_ratio": pool.PagecacheRatio,
			"uris":            flattenDesSepPoolUris(pool.URIS),
		}
		temp = append(temp, t)
	}
	return temp
}

func flattenDesSepPoolUris(uris URIList) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, uri := range uris {
		t := map[string]interface{}{
			"ip":   uri.IP,
			"port": uri.Port,
		}
		temp = append(temp, t)
	}
	return temp
}

func dataSourceSepDesSchemaMake(sh map[string]*schema.Schema) map[string]*schema.Schema {
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

func dataSourceSepSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sep_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "sep type des id",
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
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
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
			Type:     schema.TypeString,
			Computed: true,
		},
		"obj_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provided_by": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"config_string": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceSepDes() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepDesRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepDesSchemaMake(dataSourceSepSchemaMake()),
	}
}

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

func dataSourceSepTatlinRead(d *schema.ResourceData, m interface{}) error {
	tatlinSep, err := utilitySepTatlinCheckPresence(d, m)
	if err != nil {
		return err
	}
	id := uuid.New()
	d.SetId(id.String())

	d.Set("ckey", tatlinSep.Ckey)
	d.Set("meta", flattenMeta(tatlinSep.Meta))
	d.Set("consumed_by", tatlinSep.ConsumedBy)
	d.Set("desc", tatlinSep.Desc)
	d.Set("gid", tatlinSep.Gid)
	d.Set("guid", tatlinSep.Guid)
	d.Set("sep_id", tatlinSep.Id)
	d.Set("milestones", tatlinSep.Milestones)
	d.Set("name", tatlinSep.Name)
	d.Set("obj_status", tatlinSep.ObjStatus)
	d.Set("provided_by", tatlinSep.ProvidedBy)
	d.Set("tech_status", tatlinSep.TechStatus)
	d.Set("type", tatlinSep.Type)
	data, _ := json.Marshal(tatlinSep.Config)
	d.Set("config_string", string(data))
	d.Set("config", flattenSepTatlinConfig(tatlinSep.Config))

	return nil
}

func flattenSepTatlinConfig(tc TatlinConfigSep) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"api_urls":          tc.ApiUrls,
			"disk_max_size":     tc.DiskMaxSize,
			"edgeuser_password": tc.EdgeuserPassword,
			"edgeuser_name":     tc.EdgeuserName,
			"format":            tc.Format,
			"host_group_name":   tc.HostGroupName,
			"housekeeping_settings": []interface{}{
				map[string]interface{}{
					"disk_del_queue": []interface{}{
						map[string]interface{}{
							"purgatory_id":             tc.HousekeepingSettings.DiskDelQueue.PurgatoryId,
							"chunk_max_size":           tc.HousekeepingSettings.DiskDelQueue.ChunkMaxSize,
							"disk_count_max":           tc.HousekeepingSettings.DiskDelQueue.DiskCountMax,
							"enabled":                  tc.HousekeepingSettings.DiskDelQueue.Enabled,
							"normal_time_to_sleep":     tc.HousekeepingSettings.DiskDelQueue.NormalTimeToSleep,
							"one_minute_la_threshold":  tc.HousekeepingSettings.DiskDelQueue.OneMinuteLaThreshold,
							"oversize_time_to_sleep":   tc.HousekeepingSettings.DiskDelQueue.OversizeTimeToSleep,
							"purge_attempts_threshold": tc.HousekeepingSettings.DiskDelQueue.PurgeAttemptsThreshold,
						},
					},
				},
			},
			"ovs_settings": []interface{}{
				map[string]interface{}{
					"vpool_data_metadatacache":   tc.OVSSettings.VPoolDataMetadataCache,
					"vpool_vmstor_metadatacache": tc.OVSSettings.VPoolVMstorMetadataCache,
				},
			},
			"mgmt_password": tc.MGMTPassword,
			"mgmt_user":     tc.MGMTUser,
			"model":         tc.Model,
			"name_prefix":   tc.NamePrefix,
			"pools":         flattenSepTatlinPools(tc.Pools),
			"ports":         flattenSepTatlinPorts(tc.Ports),
			"protocol":      tc.Protocol,
			"tech_disk": []interface{}{
				map[string]interface{}{
					"name": tc.TechDisk.Name,
					"size": tc.TechDisk.Size,
					"pool": tc.TechDisk.Pool,
					"wwid": tc.TechDisk.WWID,
				},
			},
		},
	}
}

func flattenSepTatlinPools(tp PoolList) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, pool := range tp {
		t := map[string]interface{}{
			"name":        pool.Name,
			"types":       pool.Types,
			"usage_limit": pool.UsageLimit,
		}
		temp = append(temp, t)
	}

	return temp
}

func flattenSepTatlinPorts(tp TatlinPortList) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, port := range tp {
		t := map[string]interface{}{
			"ips":  port.IPS,
			"iqn":  port.IQN,
			"name": port.Name,
		}
		temp = append(temp, t)
	}

	return temp
}

func dataSourceSepTatlinSchemaMake(sh map[string]*schema.Schema) map[string]*schema.Schema {
	sh["config"] = &schema.Schema{
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
	}

	return sh
}

func dataSourceSepTatlin() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepTatlinRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepTatlinSchemaMake(dataSourceSepSchemaMake()),
	}
}

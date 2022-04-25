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

func dataSourceSepDoradoRead(d *schema.ResourceData, m interface{}) error {
	doradoSep, err := utilitySepDoradoCheckPresence(d, m)
	if err != nil {
		return err
	}
	id := uuid.New()
	d.SetId(id.String())

	d.Set("ckey", doradoSep.Ckey)
	d.Set("meta", flattenMeta(doradoSep.Meta))
	d.Set("consumed_by", doradoSep.ConsumedBy)
	d.Set("desc", doradoSep.Desc)
	d.Set("gid", doradoSep.Gid)
	d.Set("guid", doradoSep.Guid)
	d.Set("sep_id", doradoSep.Id)
	d.Set("milestones", doradoSep.Milestones)
	d.Set("name", doradoSep.Name)
	d.Set("obj_status", doradoSep.ObjStatus)
	d.Set("provided_by", doradoSep.ProvidedBy)
	d.Set("tech_status", doradoSep.TechStatus)
	d.Set("type", doradoSep.Type)
	data, _ := json.Marshal(doradoSep.Config)
	d.Set("config_string", string(data))
	d.Set("config", flattenSepDoradoConfig(doradoSep.Config))

	return nil
}

func flattenSepDoradoConfig(dc DoradoConfigSep) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"api_urls":          dc.ApiUrls,
			"disk_max_size":     dc.DiskMaxSize,
			"edgeuser_password": dc.EdgeuserPassword,
			"edgeuser_name":     dc.EdgeuserName,
			"format":            dc.Format,
			"groups": []interface{}{
				map[string]interface{}{
					"host_group": dc.Groups.HostGroup,
					"lung_group": dc.Groups.LungGroup,
					"port_group": dc.Groups.PortGroup,
				},
			},
			"ovs_settings": []interface{}{
				map[string]interface{}{
					"vpool_data_metadatacache":   dc.OVSSettings.VPoolDataMetadataCache,
					"vpool_vmstor_metadatacache": dc.OVSSettings.VPoolVMstorMetadataCache,
				},
			},
			"host_group_name": dc.HostGroupName,
			"housekeeping_settings": []interface{}{
				map[string]interface{}{
					"disk_del_queue": []interface{}{
						map[string]interface{}{
							"purgatory_id":             dc.HousekeepingSettings.DiskDelQueue.PurgatoryId,
							"chunk_max_size":           dc.HousekeepingSettings.DiskDelQueue.ChunkMaxSize,
							"disk_count_max":           dc.HousekeepingSettings.DiskDelQueue.DiskCountMax,
							"enabled":                  dc.HousekeepingSettings.DiskDelQueue.Enabled,
							"normal_time_to_sleep":     dc.HousekeepingSettings.DiskDelQueue.NormalTimeToSleep,
							"one_minute_la_threshold":  dc.HousekeepingSettings.DiskDelQueue.OneMinuteLaThreshold,
							"oversize_time_to_sleep":   dc.HousekeepingSettings.DiskDelQueue.OversizeTimeToSleep,
							"purge_attempts_threshold": dc.HousekeepingSettings.DiskDelQueue.PurgeAttemptsThreshold,
						},
					},
				},
			},
			"mgmt_password": dc.MGMTPassword,
			"mgmt_user":     dc.MGMTUser,
			"model":         dc.Model,
			"name_prefix":   dc.NamePrefix,
			"pools":         flattenSepDoradoPools(dc.Pools),
			"ports":         flattenSepDoradoPorts(dc.Ports),
			"protocol":      dc.Protocol,
		},
	}
}

func flattenSepDoradoPorts(dp DoradoPortList) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, port := range dp {
		t := map[string]interface{}{
			"name": port.Name,
			"ip":   port.IP,
		}
		temp = append(temp, t)
	}

	return temp
}

func flattenSepDoradoPools(dp PoolList) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, pool := range dp {
		t := map[string]interface{}{
			"name":        pool.Name,
			"types":       pool.Types,
			"usage_limit": pool.UsageLimit,
		}
		temp = append(temp, t)
	}

	return temp
}

func dataSourceSepDoradoSchemaMake(sh map[string]*schema.Schema) map[string]*schema.Schema {
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
				"groups": {
					Type:     schema.TypeList,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"host_group": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"lung_group": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"port_group": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
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
							"ip": {
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
			},
		},
	}

	return sh
}

func dataSourceSepDorado() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceSepDoradoRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceSepDoradoSchemaMake(dataSourceSepSchemaMake()),
	}
}

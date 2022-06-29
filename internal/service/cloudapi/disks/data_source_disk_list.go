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

package disks

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func flattenDiskList(dl DisksListResp) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, disk := range dl {
		diskAcl, _ := json.Marshal(disk.Acl)
		diskIotune, _ := json.Marshal(disk.IOTune)
		temp := map[string]interface{}{
			"account_id":            disk.AccountID,
			"account_name":          disk.AccountName,
			"acl":                   string(diskAcl),
			"boot_partition":        disk.BootPartition,
			"compute_id":            disk.ComputeID,
			"compute_name":          disk.ComputeName,
			"created_time":          disk.CreatedTime,
			"deleted_time":          disk.DeletedTime,
			"desc":                  disk.Desc,
			"destruction_time":      disk.DestructionTime,
			"devicename":            disk.DeviceName,
			"disk_path":             disk.DiskPath,
			"gid":                   disk.GridID,
			"guid":                  disk.GUID,
			"disk_id":               disk.ID,
			"image_id":              disk.ImageID,
			"images":                disk.Images,
			"iotune":                string(diskIotune),
			"iqn":                   disk.IQN,
			"login":                 disk.Login,
			"machine_id":            disk.MachineId,
			"machine_name":          disk.MachineName,
			"milestones":            disk.Milestones,
			"name":                  disk.Name,
			"order":                 disk.Order,
			"params":                disk.Params,
			"parent_id":             disk.ParentId,
			"passwd":                disk.Passwd,
			"pci_slot":              disk.PciSlot,
			"pool":                  disk.Pool,
			"purge_attempts":        disk.PurgeAttempts,
			"purge_time":            disk.PurgeTime,
			"reality_device_number": disk.RealityDeviceNumber,
			"reference_id":          disk.ReferenceId,
			"res_id":                disk.ResID,
			"res_name":              disk.ResName,
			"role":                  disk.Role,
			"sep_id":                disk.SepID,
			"sep_type":              disk.SepType,
			"size_max":              disk.SizeMax,
			"size_used":             disk.SizeUsed,
			"snapshots":             flattendDiskSnapshotList(disk.Snapshots),
			"status":                disk.Status,
			"tech_status":           disk.TechStatus,
			"type":                  disk.Type,
			"vmid":                  disk.VMID,
			"update_by":             disk.UpdateBy,
		}
		res = append(res, temp)
	}
	return res

}

func flattendDiskSnapshotList(sl SnapshotRecordList) []interface{} {
	res := make([]interface{}, 0)
	for _, snapshot := range sl {
		temp := map[string]interface{}{
			"guid":          snapshot.Guid,
			"label":         snapshot.Label,
			"res_id":        snapshot.ResId,
			"snap_set_guid": snapshot.SnapSetGuid,
			"snap_set_time": snapshot.SnapSetTime,
			"timestamp":     snapshot.TimeStamp,
		}
		res = append(res, temp)
	}

	return res

}

func dataSourceDiskListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diskList, err := utilityDiskListCheckPresence(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenDiskList(diskList))

	return nil
}

func dataSourceDiskListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID of the account the disks belong to",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "type of the disks",
		},
		"page": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Page number",
		},
		"size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Page size",
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"acl": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"boot_partition": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"compute_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"compute_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"deleted_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"desc": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"destruction_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"devicename": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"disk_path": {
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
					"disk_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"image_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"images": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"iotune": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"iqn": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"login": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"machine_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"machine_name": {
						Type:     schema.TypeString,
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
					"order": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"params": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"parent_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"passwd": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pci_slot": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"pool": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"purge_attempts": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"purge_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"reality_device_number": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"reference_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"res_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"res_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"role": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sep_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"sep_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"size_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"size_used": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"snapshots": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"guid": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"label": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"res_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"snap_set_guid": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"snap_set_time": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"timestamp": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tech_status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vmid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"update_by": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func DataSourceDiskList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceDiskListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceDiskListSchemaMake(),
	}
}

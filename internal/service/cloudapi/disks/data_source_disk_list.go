/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

func flattenDiskComputes(computes map[string]string) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for computeKey, computeVal := range computes {
		temp := map[string]interface{}{
			"compute_id":   computeKey,
			"compute_name": computeVal,
		}
		res = append(res, temp)
	}
	return res
}

func flattenIOTune(iot IOTune) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"read_bytes_sec":      iot.ReadBytesSec,
		"read_bytes_sec_max":  iot.ReadBytesSecMax,
		"read_iops_sec":       iot.ReadIopsSec,
		"read_iops_sec_max":   iot.ReadIopsSecMax,
		"size_iops_sec":       iot.SizeIopsSec,
		"total_bytes_sec":     iot.TotalBytesSec,
		"total_bytes_sec_max": iot.TotalBytesSecMax,
		"total_iops_sec":      iot.TotalIopsSec,
		"total_iops_sec_max":  iot.TotalIopsSecMax,
		"write_bytes_sec":     iot.WriteBytesSec,
		"write_bytes_sec_max": iot.WriteBytesSecMax,
		"write_iops_sec":      iot.WriteIopsSec,
		"write_iops_sec_max":  iot.WriteIopsSecMax,
	}

	res = append(res, temp)
	return res
}

func flattenDiskList(dl DisksList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, disk := range dl {
		diskAcl, _ := json.Marshal(disk.Acl)
		temp := map[string]interface{}{
			"account_id":            disk.AccountID,
			"account_name":          disk.AccountName,
			"acl":                   string(diskAcl),
			"computes":              flattenDiskComputes(disk.Computes),
			"boot_partition":        disk.BootPartition,
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
			"iotune":                flattenIOTune(disk.IOTune),
			"iqn":                   disk.IQN,
			"login":                 disk.Login,
			"machine_id":            disk.MachineId,
			"machine_name":          disk.MachineName,
			"milestones":            disk.Milestones,
			"disk_name":             disk.Name,
			"order":                 disk.Order,
			"params":                disk.Params,
			"parent_id":             disk.ParentId,
			"passwd":                disk.Passwd,
			"pci_slot":              disk.PciSlot,
			"pool":                  disk.Pool,
			"present_to":            disk.PresentTo,
			"purge_attempts":        disk.PurgeAttempts,
			"purge_time":            disk.PurgeTime,
			"reality_device_number": disk.RealityDeviceNumber,
			"reference_id":          disk.ReferenceId,
			"res_id":                disk.ResID,
			"res_name":              disk.ResName,
			"role":                  disk.Role,
			"sep_id":                disk.SepID,
			"sep_type":              disk.SepType,
			"shareable":             disk.Shareable,
			"size_max":              disk.SizeMax,
			"size_used":             disk.SizeUsed,
			"snapshots":             flattenDiskSnapshotList(disk.Snapshots),
			"status":                disk.Status,
			"tech_status":           disk.TechStatus,
			"type":                  disk.Type,
			"vmid":                  disk.VMID,
		}
		res = append(res, temp)
	}
	return res

}

func flattenDiskSnapshotList(sl SnapshotList) []interface{} {
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
	diskList, err := utilityDiskListCheckPresence(ctx, d, m, disksListAPI)
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
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The unique ID of the subscriber-owner of the disk",
					},
					"account_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the subscriber '(account') to whom this disk belongs",
					},
					"acl": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"boot_partition": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Number of disk partitions",
					},
					"computes": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"compute_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"compute_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"created_time": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Created time",
					},
					"deleted_time": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Deleted time",
					},
					"desc": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Description of disk",
					},
					"destruction_time": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Time of final deletion",
					},
					"devicename": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the device",
					},
					"disk_path": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Disk path",
					},
					"gid": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "ID of the grid (platform)",
					},
					"guid": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Disk ID on the storage side",
					},
					"disk_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The unique ID of the subscriber-owner of the disk",
					},
					"image_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Image ID",
					},
					"images": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Description: "IDs of images using the disk",
					},
					"iotune": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"read_bytes_sec": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Number of bytes to read per second",
								},
								"read_bytes_sec_max": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Maximum number of bytes to read",
								},
								"read_iops_sec": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Number of io read operations per second",
								},
								"read_iops_sec_max": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Maximum number of io read operations",
								},
								"size_iops_sec": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Size of io operations",
								},
								"total_bytes_sec": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Total size bytes per second",
								},
								"total_bytes_sec_max": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Maximum total size of bytes per second",
								},
								"total_iops_sec": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Total number of io operations per second",
								},
								"total_iops_sec_max": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Maximum total number of io operations per second",
								},
								"write_bytes_sec": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Number of bytes to write per second",
								},
								"write_bytes_sec_max": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Maximum number of bytes to write per second",
								},
								"write_iops_sec": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Number of write operations per second",
								},
								"write_iops_sec_max": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Maximum number of write operations per second",
								},
							},
						},
					},
					"iqn": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Disk IQN",
					},
					"login": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Login to access the disk",
					},
					"machine_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Machine ID",
					},
					"machine_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Machine name",
					},
					"milestones": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Milestones",
					},
					"disk_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of disk",
					},
					"order": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Disk order",
					},
					"params": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Disk params",
					},
					"parent_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "ID of the parent disk",
					},
					"passwd": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Password to access the disk",
					},
					"pci_slot": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "ID of the pci slot to which the disk is connected",
					},
					"pool": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Pool for disk location",
					},
					"present_to": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"purge_attempts": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Number of deletion attempts",
					},
					"purge_time": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Time of the last deletion attempt",
					},
					"reality_device_number": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Reality device number",
					},
					"reference_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the reference to the disk",
					},
					"res_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Resource ID",
					},
					"res_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the resource",
					},
					"role": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Disk role",
					},
					"sep_id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Storage endpoint provider ID to create disk",
					},
					"sep_type": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Type SEP. Defines the type of storage system and contains one of the values set in the cloud platform",
					},
					"shareable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"size_max": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Size in GB",
					},
					"size_used": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Number of used space, in GB",
					},
					"snapshots": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"guid": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "ID of the snapshot",
								},
								"label": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Name of the snapshot",
								},
								"res_id": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Reference to the snapshot",
								},
								"snap_set_guid": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "The set snapshot ID",
								},
								"snap_set_time": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "The set time of the snapshot",
								},
								"timestamp": {
									Type:        schema.TypeInt,
									Computed:    true,
									Description: "Snapshot time",
								},
							},
						},
					},
					"status": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Disk status",
					},
					"tech_status": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Technical status of the disk",
					},
					"type": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The type of disk in terms of its role in compute: 'B=Boot, D=Data, T=Temp'",
					},
					"vmid": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Virtual Machine ID (Deprecated)",
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

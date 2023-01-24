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

	// "net/url"

	"github.com/google/uuid"
	"github.com/rudecs/terraform-provider-decort/internal/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())

	diskAcl, _ := json.Marshal(disk.Acl)

	d.Set("account_id", disk.AccountID)
	d.Set("account_name", disk.AccountName)
	d.Set("acl", string(diskAcl))
	d.Set("boot_partition", disk.BootPartition)
	d.Set("computes", flattenDiskComputes(disk.Computes))
	d.Set("created_time", disk.CreatedTime)
	d.Set("deleted_time", disk.DeletedTime)
	d.Set("desc", disk.Desc)
	d.Set("destruction_time", disk.DestructionTime)
	d.Set("devicename", disk.DeviceName)
	d.Set("disk_path", disk.DiskPath)
	d.Set("gid", disk.GridID)
	d.Set("guid", disk.GUID)
	d.Set("disk_id", disk.ID)
	d.Set("image_id", disk.ImageID)
	d.Set("images", disk.Images)
	d.Set("iotune", flattenIOTune(disk.IOTune))
	d.Set("iqn", disk.IQN)
	d.Set("login", disk.Login)
	d.Set("milestones", disk.Milestones)
	d.Set("disk_name", disk.Name)
	d.Set("order", disk.Order)
	d.Set("params", disk.Params)
	d.Set("parent_id", disk.ParentId)
	d.Set("passwd", disk.Passwd)
	d.Set("pci_slot", disk.PciSlot)
	d.Set("pool", disk.Pool)
	d.Set("present_to", disk.PresentTo)
	d.Set("purge_attempts", disk.PurgeAttempts)
	d.Set("purge_time", disk.PurgeTime)
	d.Set("reality_device_number", disk.RealityDeviceNumber)
	d.Set("reference_id", disk.ReferenceId)
	d.Set("res_id", disk.ResID)
	d.Set("res_name", disk.ResName)
	d.Set("role", disk.Role)
	d.Set("sep_id", disk.SepID)
	d.Set("sep_type", disk.SepType)
	d.Set("shareable", disk.Shareable)
	d.Set("size_max", disk.SizeMax)
	d.Set("size_used", disk.SizeUsed)
	d.Set("snapshots", flattenDiskSnapshotList(disk.Snapshots))
	d.Set("status", disk.Status)
	d.Set("tech_status", disk.TechStatus)
	d.Set("type", disk.Type)
	d.Set("vmid", disk.VMID)

	return nil
}

func dataSourceDiskSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"disk_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The unique ID of the subscriber-owner of the disk",
		},
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
	}

	return rets
}

func DataSourceDisk() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceDiskRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceDiskSchemaMake(),
	}
}

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
	d.Set("compute_id", disk.ComputeID)
	d.Set("compute_name", disk.ComputeName)
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
	d.Set("purge_attempts", disk.PurgeAttempts)
	d.Set("purge_time", disk.PurgeTime)
	d.Set("reality_device_number", disk.RealityDeviceNumber)
	d.Set("reference_id", disk.ReferenceId)
	d.Set("res_id", disk.ResID)
	d.Set("res_name", disk.ResName)
	d.Set("role", disk.Role)
	d.Set("sep_id", disk.SepID)
	d.Set("sep_type", disk.SepType)
	d.Set("size_max", disk.SizeMax)
	d.Set("size_used", disk.SizeUsed)
	d.Set("snapshots", flattendDiskSnapshotList(disk.Snapshots))
	d.Set("status", disk.Status)
	d.Set("tech_status", disk.TechStatus)
	d.Set("type", disk.Type)
	d.Set("vmid", disk.VMID)

	return nil
}

func dataSourceDiskSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"disk_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
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
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"read_bytes_sec": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"read_bytes_sec_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"read_iops_sec": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"read_iops_sec_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"size_iops_sec": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"total_bytes_sec": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"total_bytes_sec_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"total_iops_sec": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"total_iops_sec_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"write_bytes_sec": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"write_bytes_sec_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"write_iops_sec": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"write_iops_sec_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"iqn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"login": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk_name": {
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
			Type:     schema.TypeFloat,
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

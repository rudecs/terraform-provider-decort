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
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDiskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("accountId", fmt.Sprintf("%d", d.Get("account_id").(int)))
	urlValues.Add("gid", fmt.Sprintf("%d", d.Get("gid").(int)))
	urlValues.Add("name", d.Get("disk_name").(string))
	urlValues.Add("size", fmt.Sprintf("%d", d.Get("size_max").(int)))
	if typeRaw, ok := d.GetOk("type"); ok {
		urlValues.Add("type", strings.ToUpper(typeRaw.(string)))
	} else {
		urlValues.Add("type", "D")
	}

	if sepId, ok := d.GetOk("sep_id"); ok {
		urlValues.Add("sep_id", strconv.Itoa(sepId.(int)))
	}

	if poolName, ok := d.GetOk("pool"); ok {
		urlValues.Add("pool", poolName.(string))
	}

	argVal, argSet := d.GetOk("desc")
	if argSet {
		urlValues.Add("description", argVal.(string))
	}

	diskId, err := c.DecortAPICall(ctx, "POST", disksCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	urlValues = &url.Values{}

	d.SetId(diskId) // update ID of the resource to tell Terraform that the disk resource exists

	if iotuneRaw, ok := d.GetOk("iotune"); ok {
		iot := iotuneRaw.([]interface{})[0]
		iotune := iot.(map[string]interface{})
		urlValues.Add("diskId", diskId)
		urlValues.Add("iops", strconv.Itoa(iotune["total_iops_sec"].(int)))
		urlValues.Add("read_bytes_sec", strconv.Itoa(iotune["read_bytes_sec"].(int)))
		urlValues.Add("read_bytes_sec_max", strconv.Itoa(iotune["read_bytes_sec_max"].(int)))
		urlValues.Add("read_iops_sec", strconv.Itoa(iotune["read_iops_sec"].(int)))
		urlValues.Add("read_iops_sec_max", strconv.Itoa(iotune["read_iops_sec_max"].(int)))
		urlValues.Add("size_iops_sec", strconv.Itoa(iotune["size_iops_sec"].(int)))
		urlValues.Add("total_bytes_sec", strconv.Itoa(iotune["total_bytes_sec"].(int)))
		urlValues.Add("total_bytes_sec_max", strconv.Itoa(iotune["total_bytes_sec_max"].(int)))
		urlValues.Add("total_iops_sec_max", strconv.Itoa(iotune["total_iops_sec_max"].(int)))
		urlValues.Add("write_bytes_sec", strconv.Itoa(iotune["write_bytes_sec"].(int)))
		urlValues.Add("write_bytes_sec_max", strconv.Itoa(iotune["write_bytes_sec_max"].(int)))
		urlValues.Add("write_iops_sec", strconv.Itoa(iotune["write_iops_sec"].(int)))
		urlValues.Add("write_iops_sec_max", strconv.Itoa(iotune["write_iops_sec_max"].(int)))

		_, err := c.DecortAPICall(ctx, "POST", disksIOLimitAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	dgn := resourceDiskRead(ctx, d, m)
	if dgn != nil {
		return dgn
	}

	return nil
}

func resourceDiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if disk == nil {
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

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

func resourceDiskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if d.HasChange("size_max") {
		oldSize, newSize := d.GetChange("size_max")
		if oldSize.(int) < newSize.(int) {
			log.Debugf("resourceDiskUpdate: resizing disk ID %s - %d GB -> %d GB",
				d.Id(), oldSize.(int), newSize.(int))
			urlValues.Add("diskId", d.Id())
			urlValues.Add("size", fmt.Sprintf("%d", newSize.(int)))
			_, err := c.DecortAPICall(ctx, "POST", disksResizeAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set("size_max", newSize)
		} else if oldSize.(int) > newSize.(int) {
			return diag.FromErr(fmt.Errorf("resourceDiskUpdate: Disk ID %s - reducing disk size is not allowed", d.Id()))
		}
		urlValues = &url.Values{}
	}

	if d.HasChange("disk_name") {
		urlValues.Add("diskId", d.Id())
		urlValues.Add("name", d.Get("disk_name").(string))
		_, err := c.DecortAPICall(ctx, "POST", disksRenameAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("iotune") {
		iot := d.Get("iotune").([]interface{})[0]
		iotune := iot.(map[string]interface{})
		urlValues.Add("diskId", d.Id())
		urlValues.Add("iops", strconv.Itoa(iotune["total_iops_sec"].(int)))
		urlValues.Add("read_bytes_sec", strconv.Itoa(iotune["read_bytes_sec"].(int)))
		urlValues.Add("read_bytes_sec_max", strconv.Itoa(iotune["read_bytes_sec_max"].(int)))
		urlValues.Add("read_iops_sec", strconv.Itoa(iotune["read_iops_sec"].(int)))
		urlValues.Add("read_iops_sec_max", strconv.Itoa(iotune["read_iops_sec_max"].(int)))
		urlValues.Add("size_iops_sec", strconv.Itoa(iotune["size_iops_sec"].(int)))
		urlValues.Add("total_bytes_sec", strconv.Itoa(iotune["total_bytes_sec"].(int)))
		urlValues.Add("total_bytes_sec_max", strconv.Itoa(iotune["total_bytes_sec_max"].(int)))
		urlValues.Add("total_iops_sec_max", strconv.Itoa(iotune["total_iops_sec_max"].(int)))
		urlValues.Add("write_bytes_sec", strconv.Itoa(iotune["write_bytes_sec"].(int)))
		urlValues.Add("write_bytes_sec_max", strconv.Itoa(iotune["write_bytes_sec_max"].(int)))
		urlValues.Add("write_iops_sec", strconv.Itoa(iotune["write_iops_sec"].(int)))
		urlValues.Add("write_iops_sec_max", strconv.Itoa(iotune["write_iops_sec_max"].(int)))

		_, err := c.DecortAPICall(ctx, "POST", disksIOLimitAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("restore") {
		if d.Get("restore").(bool) {
			urlValues.Add("diskId", d.Id())
			urlValues.Add("reason", d.Get("reason").(string))

			_, err := c.DecortAPICall(ctx, "POST", disksRestoreAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}

			urlValues = &url.Values{}
		}

	}

	return resourceDiskRead(ctx, d, m)
}

func resourceDiskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if disk == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	params := &url.Values{}
	params.Add("diskId", d.Id())
	params.Add("detach", strconv.FormatBool(d.Get("detach").(bool)))
	params.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))
	params.Add("reason", d.Get("reason").(string))

	c := m.(*controller.ControllerCfg)
	_, err = c.DecortAPICall(ctx, "POST", disksDeleteAPI, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDiskExists(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {

	diskFacts, err := utilityDiskCheckPresence(ctx, d, m)
	if diskFacts == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func resourceDiskSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"account_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"disk_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"size_max": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"gid": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"pool": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"sep_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"desc": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"D", "B", "T"}, false),
		},

		"detach": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "detach disk from machine first",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "whether to completely delete the disk, works only with non attached disks",
		},
		"reason": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "reason for an action",
		},
		"restore": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "restore deleting disk",
		},

		"disk_id": {
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
			Optional: true,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"read_bytes_sec": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"read_bytes_sec_max": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"read_iops_sec": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"read_iops_sec_max": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"size_iops_sec": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"total_bytes_sec": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"total_bytes_sec_max": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"total_iops_sec": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"total_iops_sec_max": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"write_bytes_sec": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"write_bytes_sec_max": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"write_iops_sec": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"write_iops_sec_max": {
						Type:     schema.TypeInt,
						Optional: true,
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

		"sep_type": {
			Type:     schema.TypeString,
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
		"vmid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}

	return rets
}

func ResourceDisk() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceDiskCreate,
		ReadContext:   resourceDiskRead,
		UpdateContext: resourceDiskUpdate,
		DeleteContext: resourceDiskDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout180s,
			Read:    &constants.Timeout30s,
			Update:  &constants.Timeout180s,
			Delete:  &constants.Timeout60s,
			Default: &constants.Timeout60s,
		},

		Schema: resourceDiskSchemaMake(),
	}
}

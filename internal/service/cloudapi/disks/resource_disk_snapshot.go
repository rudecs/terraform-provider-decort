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
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceDiskSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	urlValues := &url.Values{}
	c := m.(*controller.ControllerCfg)
	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if disk == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
	snapshots := disk.Snapshots
	snapshot := Snapshot{}
	label := d.Get("label").(string)
	for _, sn := range snapshots {
		if label == sn.Label {
			snapshot = sn
			break
		}
	}
	if label != snapshot.Label {
		return diag.Errorf("Snapshot with label \"%v\" not found", label)
	}
	if rollback := d.Get("rollback").(bool); rollback {
		urlValues.Add("diskId", strconv.Itoa(d.Get("disk_id").(int)))
		urlValues.Add("label", label)
		urlValues.Add("timestamp", strconv.Itoa(d.Get("timestamp").(int)))
		log.Debugf("resourceDiskCreate: Snapshot rollback with label", label)
		_, err := c.DecortAPICall(ctx, "POST", disksSnapshotRollbackAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		urlValues = &url.Values{}
	}
	return resourceDiskSnapshotRead(ctx, d, m)
}

func resourceDiskSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if disk == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
	snapshots := disk.Snapshots
	snapshot := Snapshot{}
	label := d.Get("label").(string)
	for _, sn := range snapshots {
		if label == sn.Label {
			snapshot = sn
			break
		}
	}
	if label != snapshot.Label {
		return diag.Errorf("Snapshot with label \"%v\" not found", label)
	}

	d.SetId(d.Get("label").(string))
	d.Set("timestamp", snapshot.TimeStamp)
	d.Set("guid", snapshot.Guid)
	d.Set("res_id", snapshot.ResId)
	d.Set("snap_set_guid", snapshot.SnapSetGuid)
	d.Set("snap_set_time", snapshot.SnapSetTime)
	return nil
}

func resourceDiskSnapshotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	urlValues := &url.Values{}
	c := m.(*controller.ControllerCfg)
	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if disk == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
	snapshots := disk.Snapshots
	snapshot := Snapshot{}
	label := d.Get("label").(string)
	for _, sn := range snapshots {
		if label == sn.Label {
			snapshot = sn
			break
		}
	}
	if label != snapshot.Label {
		return diag.Errorf("Snapshot with label \"%v\" not found", label)
	}
	if d.HasChange("rollback") && d.Get("rollback").(bool) == true {
		urlValues.Add("diskId", strconv.Itoa(d.Get("disk_id").(int)))
		urlValues.Add("label", label)
		urlValues.Add("timestamp", strconv.Itoa(d.Get("timestamp").(int)))
		log.Debugf("resourceDiskUpdtae: Snapshot rollback with label", label)
		_, err := c.DecortAPICall(ctx, "POST", disksSnapshotRollbackAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		urlValues = &url.Values{}
	}

	return resourceDiskSnapshotRead(ctx, d, m)
}

func resourceDiskSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*controller.ControllerCfg)

	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if disk == nil { //if disk not exits, can't call snapshotDelete
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	params := &url.Values{}
	params.Add("diskId", strconv.Itoa(d.Get("disk_id").(int)))
	params.Add("label", d.Get("label").(string))

	_, err = c.DecortAPICall(ctx, "POST", disksSnapshotDeleteAPI, params)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDiskSnapshotSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"disk_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The unique ID of the subscriber-owner of the disk",
		},
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the snapshot",
		},
		"rollback": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Needed in order to make a snapshot rollback",
		},
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the snapshot",
		},
		"timestamp": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Snapshot time",
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
	}
	return rets
}

func ResourceDiskSnapshot() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceDiskSnapshotCreate,
		ReadContext:   resourceDiskSnapshotRead,
		UpdateContext: resourceDiskSnapshotUpdate,
		DeleteContext: resourceDiskSnapshotDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout600s,
			Read:    &constants.Timeout300s,
			Update:  &constants.Timeout300s,
			Delete:  &constants.Timeout300s,
			Default: &constants.Timeout300s,
		},

		Schema: resourceDiskSnapshotSchemaMake(),
	}
}

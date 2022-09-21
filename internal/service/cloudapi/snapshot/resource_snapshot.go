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

package snapshot

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSnapshotCreate: called for snapshot %s", d.Get("label").(string))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("label", d.Get("label").(string))
	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))

	snapshotId, err := c.DecortAPICall(ctx, "POST", snapshotCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshotId = strings.ReplaceAll(snapshotId, "\"", "")

	d.SetId(snapshotId)
	d.Set("guid", snapshotId)

	diagnostics := resourceSnapshotRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	snapshot, err := utilitySnapshotCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("timestamp", snapshot.Timestamp)
	d.Set("guid", snapshot.Guid)
	d.Set("disks", snapshot.Disks)
	d.Set("label", snapshot.Label)

	return nil
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceSnapshotDelete: called for %s, id: %s", d.Get("label").(string), d.Id())

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	urlValues.Add("label", d.Get("label").(string))

	_, err := c.DecortAPICall(ctx, "POST", snapshotDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceSnapshotEdit(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChange("rollback") {
		if d.Get("rollback").(bool) {
			err := resourceSnapshotRollback(ctx, d, m)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceSnapshotRollback(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	urlValues.Add("label", d.Get("label").(string))

	_, err := c.DecortAPICall(ctx, "POST", snapshotRollbackAPI, urlValues)
	if err != nil {
		return err
	}
	return nil
}

func resourceSnapshotSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the compute instance to create snapshot for.",
		},
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "text label for snapshot. Must be unique among this compute snapshots.",
		},
		"rollback": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "is rollback the snapshot",
		},
		"disks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "guid of the snapshot",
		},
		"timestamp": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "timestamp",
		},
	}
}

func ResourceSnapshot() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceSnapshotCreate,
		ReadContext:   resourceSnapshotRead,
		UpdateContext: resourceSnapshotEdit,
		DeleteContext: resourceSnapshotDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout60s,
			Read:    &constants.Timeout30s,
			Update:  &constants.Timeout60s,
			Delete:  &constants.Timeout60s,
			Default: &constants.Timeout60s,
		},

		Schema: resourceSnapshotSchemaMake(),
	}
}

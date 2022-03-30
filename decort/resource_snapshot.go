/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceSnapshotCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSnapshotCreate: called for snapshot %s", d.Get("label").(string))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("label", d.Get("label").(string))
	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))

	snapshotId, err := controller.decortAPICall("POST", snapshotCreateAPI, urlValues)
	if err != nil {
		return err
	}

	snapshotId = strings.ReplaceAll(snapshotId, "\"", "")

	d.SetId(snapshotId)
	d.Set("guid", snapshotId)

	err = resourceSnapshotRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceSnapshotRead(d *schema.ResourceData, m interface{}) error {
	snapshot, err := utilitySnapshotCheckPresence(d, m)
	if err != nil {
		return err
	}

	d.Set("timestamp", snapshot.Timestamp)
	d.Set("guid", snapshot.Guid)
	d.Set("disks", snapshot.Disks)
	d.Set("label", snapshot.Label)

	return nil
}

func resourceSnapshotDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceSnapshotDelete: called for %s, id: %s", d.Get("label").(string), d.Id())

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	urlValues.Add("label", d.Get("label").(string))

	_, err := controller.decortAPICall("POST", snapshotDeleteAPI, urlValues)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceSnapshotExists(d *schema.ResourceData, m interface{}) (bool, error) {
	snapshot, err := utilitySnapshotCheckPresence(d, m)
	if err != nil {
		return false, err
	}
	if snapshot == nil {
		return false, nil
	}

	return true, nil
}

func resourceSnapshotEdit(d *schema.ResourceData, m interface{}) error {
	err := resourceSnapshotRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceSnapshotRollback(d *schema.ResourceDiff, m interface{}) error {
	c := m.(*ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	urlValues.Add("label", d.Get("label").(string))

	_, err := c.decortAPICall("POST", snapshotRollbackAPI, urlValues)
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

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceSnapshotCreate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotEdit,
		Delete: resourceSnapshotDelete,
		Exists: resourceSnapshotExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &Timeout60s,
			Read:    &Timeout30s,
			Update:  &Timeout60s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

		CustomizeDiff: customdiff.All(
			customdiff.IfValueChange("rollback", func(old, new, meta interface{}) bool {
				o := old.(bool)
				if o != new.(bool) && o == false {
					return true
				}
				return false
			}, resourceSnapshotRollback),
		),

		Schema: resourceSnapshotSchemaMake(),
	}
}

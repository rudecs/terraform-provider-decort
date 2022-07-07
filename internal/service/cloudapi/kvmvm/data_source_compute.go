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

package kvmvm

import (
	"context"
	"encoding/json"
	"fmt"

	// "net/url"

	"github.com/rudecs/terraform-provider-decort/internal/constants"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Parse list of all disks from API compute/get into a list of "extra disks" attached to this compute
// Extra disks are all compute disks but a boot disk.
func parseComputeDisksToExtraDisks(disks []DiskRecord) []interface{} {
	// this return value will be used to d.Set("extra_disks",) item of dataSourceCompute schema,
	// which is a simple list of integer disk IDs excluding boot disk ID
	length := len(disks)
	log.Debugf("parseComputeDisksToExtraDisks: called for %d disks", length)

	if length == 0 || (length == 1 && disks[0].Type == "B") {
		// the disk list is empty (which is kind of strange - diskless compute?), or
		// there is only one disk in the list and it is a boot disk;
		// as we skip boot disks, the result will be of 0 length anyway
		return make([]interface{}, 0)
	}

	result := make([]interface{}, length-1)
	idx := 0
	for _, value := range disks {
		if value.Type == "B" {
			// skip boot disk when iterating over the list of disks
			continue
		}

		result[idx] = value.ID
		idx++
	}

	return result
}

func parseBootDiskSize(disks []DiskRecord) int {
	// this return value will be used to d.Set("boot_disk_size",) item of dataSourceCompute schema
	if len(disks) == 0 {
		return 0
	}

	for _, value := range disks {
		if value.Type == "B" {
			return value.SizeMax
		}
	}

	return 0
}

func parseBootDiskId(disks []DiskRecord) uint {
	// this return value will be used to d.Set("boot_disk_id",) item of dataSourceCompute schema
	if len(disks) == 0 {
		return 0
	}

	for _, value := range disks {
		if value.Type == "B" {
			return value.ID
		}
	}

	return 0
}

// Parse the list of interfaces from compute/get response into a list of networks
// attached to this compute
func parseComputeInterfacesToNetworks(ifaces []InterfaceRecord) []interface{} {
	// return value will be used to d.Set("network") item of dataSourceCompute schema
	length := len(ifaces)
	log.Debugf("parseComputeInterfacesToNetworks: called for %d ifaces", length)

	result := []interface{}{}

	for _, value := range ifaces {
		elem := make(map[string]interface{})
		// Keys in this map should correspond to the Schema definition
		// as returned by networkSubresourceSchemaMake()
		elem["net_id"] = value.NetID
		elem["net_type"] = value.NetType
		elem["ip_address"] = value.IPAddress
		elem["mac"] = value.MAC

		// log.Debugf("   element %d: net_id=%d, net_type=%s", i, value.NetID, value.NetType)

		result = append(result, elem)
	}

	return result
}

func flattenCompute(d *schema.ResourceData, compFacts string) error {
	// This function expects that compFacts string contains response from API compute/get,
	// i.e. detailed information about compute instance.
	//
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceComputeExists(...) method
	model := ComputeGetResp{}
	log.Debugf("flattenCompute: ready to unmarshal string %s", compFacts)
	err := json.Unmarshal([]byte(compFacts), &model)
	if err != nil {
		return err
	}

	log.Debugf("flattenCompute: ID %d, RG ID %d", model.ID, model.RgID)

	d.SetId(fmt.Sprintf("%d", model.ID))
	// d.Set("compute_id", model.ID) - we should NOT set compute_id in the schema here: if it was set - it is already set, if it wasn't - we shouldn't
	d.Set("name", model.Name)
	d.Set("rg_id", model.RgID)
	d.Set("rg_name", model.RgName)
	d.Set("account_id", model.AccountID)
	d.Set("account_name", model.AccountName)
	d.Set("driver", model.Driver)
	d.Set("cpu", model.Cpu)
	d.Set("ram", model.Ram)
	// d.Set("boot_disk_size", model.BootDiskSize) - bootdiskSize key in API compute/get is always zero, so we set boot_disk_size in another way
	d.Set("boot_disk_size", parseBootDiskSize(model.Disks))
	d.Set("boot_disk_id", parseBootDiskId(model.Disks)) // we may need boot disk ID in resize operations
	d.Set("image_id", model.ImageID)
	d.Set("description", model.Desc)
	d.Set("cloud_init", "applied") // NOTE: for existing compute we hard-code this value as an indicator for DiffSuppress fucntion
	// d.Set("status", model.Status)
	// d.Set("tech_status", model.TechStatus)

	if model.TechStatus == "STARTED" {
		d.Set("started", true)
	} else {
		d.Set("started", false)
	}

	if len(model.Disks) > 0 {
		log.Debugf("flattenCompute: calling parseComputeDisksToExtraDisks for %d disks", len(model.Disks))
		if err = d.Set("extra_disks", parseComputeDisksToExtraDisks(model.Disks)); err != nil {
			return err
		}
	}

	if len(model.Interfaces) > 0 {
		log.Debugf("flattenCompute: calling parseComputeInterfacesToNetworks for %d interfaces", len(model.Interfaces))
		if err = d.Set("network", parseComputeInterfacesToNetworks(model.Interfaces)); err != nil {
			return err
		}
	}

	if len(model.OsUsers) > 0 {
		log.Debugf("flattenCompute: calling parseOsUsers for %d logins", len(model.OsUsers))
		if err = d.Set("os_users", parseOsUsers(model.OsUsers)); err != nil {
			return err
		}
	}

	return nil
}

func dataSourceComputeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	compFacts, err := utilityComputeCheckPresence(ctx, d, m)
	if compFacts == "" {
		// if empty string is returned from utilityComputeCheckPresence then there is no
		// such Compute and err tells so - just return it to the calling party
		d.SetId("") // ensure ID is empty
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenCompute(d, compFacts))
}

func DataSourceCompute() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of this compute instance. NOTE: this parameter is case sensitive.",
			},

			// TODO: consider removing compute_id from the schema, as it not practical to call this data provider if
			// corresponding compute ID is already known
			"compute_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the compute instance. If ID is specified, name and resource group ID are ignored.",
			},

			"rg_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the resource group where this compute instance is located.",
			},

			"rg_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the resource group where this compute instance is located.",
			},

			"account_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the account this compute instance belongs to.",
			},

			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the account this compute instance belongs to.",
			},

			"driver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hardware architecture of this compute instance.",
			},

			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of CPUs allocated for this compute instance.",
			},

			"ram": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Amount of RAM in MB allocated for this compute instance.",
			},

			"image_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the OS image this compute instance is based on.",
			},

			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the OS image this compute instance is based on.",
			},

			"boot_disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This compute instance boot disk size in GB.",
			},

			"boot_disk_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This compute instance boot disk ID.",
			},

			"extra_disks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "IDs of the extra disk(s) attached to this compute.",
			},

			/*
				"disks": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: dataSourceDiskSchemaMake(), // ID, type,  name, size, account ID, SEP ID, SEP type, pool, status, tech status, compute ID, image ID
					},
					Description: "Detailed specification for all disks attached to this compute instance (including bood disk).",
				},
			*/

			"network": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: constants.MaxNetworksPerCompute,
				Elem: &schema.Resource{
					Schema: networkSubresourceSchemaMake(),
				},
				Description: "Network connection(s) for this compute.",
			},

			"os_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: osUsersSubresourceSchemaMake(),
				},
				Description: "Guest OS users provisioned on this compute instance.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-defined text description of this compute instance.",
			},

			"cloud_init": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Placeholder for cloud_init parameters.",
			},

			"started": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is compute started.",
			},
		},
	}
}

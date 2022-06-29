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

package extnet

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/flattens"
)

func dataSourceExtnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	e, err := utilityExtnetCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("ckey", e.CKey)
	d.Set("meta", flattens.FlattenMeta(e.Meta))
	d.Set("check__ips", e.CheckIPs)
	d.Set("check_ips", e.CheckIps)
	d.Set("default", e.Default)
	d.Set("default_qos", flattenExtnetDefaultQos(e.DefaultQos))
	d.Set("desc", e.Desc)
	d.Set("dns", e.Dns)
	d.Set("excluded", e.Excluded)
	d.Set("free_ips", e.FreeIps)
	d.Set("gateway", e.Gateway)
	d.Set("gid", e.GID)
	d.Set("guid", e.GUID)
	d.Set("ipcidr", e.IPCidr)
	d.Set("milestones", e.Milestones)
	d.Set("net_name", e.Name)
	d.Set("network", e.Network)
	d.Set("network_id", e.NetworkId)
	d.Set("pre_reservations_num", e.PreReservationsNum)
	d.Set("prefix", e.Prefix)
	d.Set("pri_vnf_dev_id", e.PriVnfDevId)
	d.Set("reservations", flattenExtnetReservations(e.Reservations))
	d.Set("shared_with", e.SharedWith)
	d.Set("status", e.Status)
	d.Set("vlan_id", e.VlanID)
	d.Set("vnfs", flattenExtnetVNFS(e.VNFS))
	return nil
}

func flattenExtnetReservations(ers ExtnetReservations) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, er := range ers {
		temp := map[string]interface{}{
			"client_type": er.ClientType,
			"domainname":  er.DomainName,
			"hostname":    er.HostName,
			"desc":        er.Desc,
			"ip":          er.IP,
			"mac":         er.MAC,
			"type":        er.Type,
			"vm_id":       er.VMID,
		}
		res = append(res, temp)
	}

	return res
}

func flattenExtnetDefaultQos(edqos ExtnetQos) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"e_rate":   edqos.ERate,
		"guid":     edqos.GUID,
		"in_burst": edqos.InBurst,
		"in_rate":  edqos.InRate,
	}
	res = append(res, temp)
	return res
}

func flattenExtnetVNFS(evnfs ExtnetVNFS) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"dhcp": evnfs.DHCP,
	}
	res = append(res, temp)
	return res
}

func dataSourceExtnetSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"net_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"meta": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "meta",
		},
		"check__ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"check_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"default": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"default_qos": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"e_rate": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"in_burst": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"in_rate": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dns": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"excluded": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"free_ips": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gateway": {
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
		"ipcidr": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"net_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pre_reservations_num": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"prefix": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pri_vnf_dev_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"reservations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"client_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"domainname": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"hostname": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"desc": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"mac": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vm_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"shared_with": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vlan_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vnfs": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"dhcp": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func DataSourceExtnet() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceExtnetRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceExtnetSchemaMake(),
	}
}

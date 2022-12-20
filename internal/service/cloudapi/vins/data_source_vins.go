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

package vins

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceVinsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vins, err := utilityDataVinsCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatUint(vins.ID, 10))
	flattenVinsData(d, *vins)
	return nil
}

func vnfConfigMGMTSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_addr": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ssh_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"user": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func vnfConfigResourcesSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"stack_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"uuid": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func qosSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"e_rate": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"in_brust": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"in_rate": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func vnfInterfaceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"conn_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"conn_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"def_gw": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"flipgroup_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"listen_ssh": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"net_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"net_mask": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"net_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pci_slot": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"qos": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: qosSchemaMake(),
			},
		},
		"target": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vnfs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}
}

func vnfConfigSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mgmt": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: vnfConfigMGMTSchemaMake(),
			},
		},
		"resources": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: vnfConfigResourcesSchemaMake(),
			},
		},
	}
}

func vnfDevSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"_ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Unique ID of the account, which this ViNS belongs to.",
		},
		"capabilities": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"config": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: vnfConfigSchemaMake(),
			},
		},
		"config_saved": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"custom_pre_cfg": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"desc": {
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
		"vnf_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"interfaces": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: vnfInterfaceSchemaMake(),
			},
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vnf_name": {
			Type:     schema.TypeString,
			Computed: true,
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
		"vins": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}
}

func vinsComputeSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"compute_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func reservationSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"desc": {
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
	}
}

func dhcpConfigSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_gw": {
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
		"ip_end": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_start": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"lease": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"netmask": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"reservations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: reservationSchemaMake(),
			},
		},
	}
}

func devicesPrimarySchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dev_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"iface01": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"iface02": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func devicesSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"primary": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: devicesPrimarySchemaMake(),
			},
		},
	}
}

func dhcpSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"_ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"config": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dhcpConfigSchemaMake(),
			},
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"devices": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: devicesSchemaMake(),
			},
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dhcp_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"owner_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"owner_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pure_virtual": {
			Type:     schema.TypeBool,
			Computed: true,
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
	}
}

func gwConfigSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_gw": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_net_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ext_net_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_netmask": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"qos": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: qosSchemaMake(),
			},
		},
	}
}

func gwSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"_ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"config": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: gwConfigSchemaMake(),
			},
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"devices": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: devicesSchemaMake(),
			},
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gw_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"owner_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"owner_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pure_virtual": {
			Type:     schema.TypeBool,
			Computed: true,
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
	}
}

func rulesSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"local_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"local_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_port_end": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"public_port_start": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vm_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vm_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func configSchrmaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"net_mask": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rulesSchemaMake(),
			},
		},
	}
}

func natSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"_ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"config": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: configSchrmaMake(),
			},
		},
		"devices": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: devicesSchemaMake(),
			},
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"nat_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"owner_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"owner_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pure_virtual": {
			Type:     schema.TypeBool,
			Computed: true,
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
	}
}

func vnfsSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dhcp": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dhcpSchemaMake(),
			},
		},
		"gw": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: gwSchemaMake(),
			},
		},
		"nat": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: natSchemaMake(),
			},
		},
	}
}

func dataSourceVinsSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"vins_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Unique ID of the ViNS. If ViNS ID is specified, then ViNS name, rg_id and account_id are ignored.",
		},

		"vnf_dev": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: vnfDevSchemaMake(),
			},
		},
		"_ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Unique ID of the account, which this ViNS belongs to.",
		},
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the account, which this ViNS belongs to.",
		},
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: vinsComputeSchemaMake(),
			},
		},
		"default_gw": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"default_qos": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: qosSchemaMake(),
			},
		},
		"desc": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "User-defined text description of this ViNS.",
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"manager_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"manager_type": {
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
		"net_mask": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pre_reservations_num": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"redundant": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"rg_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Unique ID of the resource group, where this ViNS is belongs to (for ViNS created at resource group level, 0 otherwise).",
		},
		"rg_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sec_vnf_dev_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"user_managed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"vnfs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: vnfsSchemaMake(),
			},
		},
		"vxlan_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
	return rets
}

func DataSourceVins() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceVinsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceVinsSchemaMake(),
	}
}

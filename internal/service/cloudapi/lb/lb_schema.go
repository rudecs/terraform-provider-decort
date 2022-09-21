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

package lb

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func createLBSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ha_mode": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"backends": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"algorithm": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"server_default_settings": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"downinter": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"fall": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"guid": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"inter": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"maxconn": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"maxqueue": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"rise": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"slowstart": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"weight": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"servers": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"check": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"guid": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"port": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"server_settings": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"downinter": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"fall": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"guid": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"inter": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"maxconn": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"maxqueue": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"rise": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"slowstart": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"weight": {
												Type:     schema.TypeInt,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"created_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deleted_by": {
			Type:     schema.TypeString,
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
		"dp_api_user": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"extnet_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"frontends": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"backend": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"bindings": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"guid": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"port": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
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
		"lb_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeInt,
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
		"primary_node": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"backend_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"compute_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"frontend_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"mgmt_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"rg_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rg_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"secondary_node": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"backend_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"compute_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"frontend_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"mgmt_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network_id": {
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
		"updated_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vins_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

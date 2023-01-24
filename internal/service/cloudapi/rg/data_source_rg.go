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

package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"

	// "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityDataResgroupCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*ResgroupGetResp, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	rgData := &ResgroupGetResp{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	rgRaw, err := c.DecortAPICall(ctx, "POST", ResgroupGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(rgRaw), rgData)
	if err != nil {
		return nil, err
	}
	return rgData, nil
}

func dataSourceResgroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rg, err := utilityDataResgroupCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("") // ensure ID is empty in this case
		return diag.FromErr(err)
	}
	return diag.FromErr(flattenDataResgroup(d, *rg))
}

func DataSourceResgroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceResgroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the resource group. Names are case sensitive and unique within the context of an account.",
			},

			"rg_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the resource group. If this ID is specified, then resource group name is ignored.",
			},

			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the account, which this resource group belongs to.",
			},

			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique ID of the account, which this resource group belongs to.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-defined text description of this resource group.",
			},
			"gid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the grid, where this resource group is deployed.",
			},
			"quota": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: quotaRgSubresourceSchemaMake(), // this is a dictionary
				},
				Description: "Quota settings for this resource group.",
			},

			"def_net_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the default network for this resource group.",
			},

			"def_net_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the default network for this resource group (if any).",
			},

			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"disksize": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"extips": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"exttraffic": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"gpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"seps": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sep_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"data_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"disk_size": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"disk_size_max": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"reserved": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"disksize": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"extips": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"exttraffic": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"gpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"seps": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sep_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"data_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"disk_size": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"disk_size_max": {
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

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of this resource group.",
			},

			"vins": {
				Type:     schema.TypeList, // this is a list of ints
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of VINs deployed in this resource group.",
			},

			"vms": {
				Type:     schema.TypeList, //t his is a list of ints
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of computes deployed in this resource group.",
			},
		},
	}
}

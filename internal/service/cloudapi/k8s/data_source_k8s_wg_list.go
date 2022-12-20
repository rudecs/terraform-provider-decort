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

package k8s

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
)

func utilityK8sWgListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (K8SGroupList, error) {

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))

	resp, err := c.DecortAPICall(ctx, "POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	var k8s K8SRecord
	if err := json.Unmarshal([]byte(resp), &k8s); err != nil {
		return nil, err
	}

	return k8s.K8SGroups.Workers, nil
}

func dataSourceK8sWgListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	wgList, err := utilityK8sWgListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("k8s_id").(int)))

	workersComputeList := make(map[uint64][]kvmvm.ComputeGetResp)
	for _, worker := range wgList {
		workersComputeList[worker.ID] = make([]kvmvm.ComputeGetResp, 0, len(worker.DetailedInfo))
		for _, info := range worker.DetailedInfo {
			compute, err := utilityComputeCheckPresence(ctx, d, m, info.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			workersComputeList[worker.ID] = append(workersComputeList[worker.ID], *compute)
		}
	}
	flattenItemsWg(d, wgList, workersComputeList)
	return nil
}

func wgSchemaMake() map[string]*schema.Schema {
	wgSchema := dataSourceK8sWgSchemaMake()
	delete(wgSchema, "k8s_id")
	wgSchema["wg_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ID of k8s worker Group.",
	}
	return wgSchema
}

func dataSourceK8sWgListSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"k8s_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: wgSchemaMake(),
			},
		},
	}
}

func DataSourceK8sWgList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceK8sWgListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceK8sWgListSchemaMake(),
	}
}

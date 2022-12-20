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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
)

func utilityK8sCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*K8SRecord, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", d.Id())

	resp, err := c.DecortAPICall(ctx, "POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	k8s := K8SRecord{}
	if err := json.Unmarshal([]byte(resp), &k8s); err != nil {
		return nil, err
	}

	return &k8s, nil
}

func utilityComputeCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}, computeID uint64) (*kvmvm.ComputeGetResp, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("computeId", strconv.FormatUint(computeID, 10))

	computeRaw, err := c.DecortAPICall(ctx, "POST", kvmvm.ComputeGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	compute := &kvmvm.ComputeGetResp{}
	err = json.Unmarshal([]byte(computeRaw), compute)
	if err != nil {
		return nil, err
	}

	return compute, nil
}

func utilityDataK8sCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*K8SRecord, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	if d.Get("k8s_id") != 0 && d.Get("k8s_id") != nil {
		urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))
	} else if id := d.Id(); id != "" {
		if strings.Contains(id, "#") {
			urlValues.Add("k8sId", strings.Split(d.Id(), "#")[1])
		} else {
			urlValues.Add("k8sId", d.Id())
		}
	}
	k8sRaw, err := c.DecortAPICall(ctx, "POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	k8s := &K8SRecord{}
	err = json.Unmarshal([]byte(k8sRaw), k8s)
	if err != nil {
		return nil, err
	}
	return k8s, nil
}

func utilityK8sListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}, api string) (K8SList, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("includedeleted", "false")
	urlValues.Add("page", "0")
	urlValues.Add("size", "0")

	k8sListRaw, err := c.DecortAPICall(ctx, "POST", api, urlValues)
	if err != nil {
		return nil, err
	}

	k8sList := K8SList{}

	err = json.Unmarshal([]byte(k8sListRaw), &k8sList)
	if err != nil {
		return nil, err
	}
	return k8sList, nil
}

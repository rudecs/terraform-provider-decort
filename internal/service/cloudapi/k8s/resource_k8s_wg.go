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
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
	log "github.com/sirupsen/logrus"
)

func resourceK8sWgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceK8sWgCreate: called with k8s id %d", d.Get("k8s_id").(int))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("workerNum", strconv.Itoa(d.Get("num").(int)))
	urlValues.Add("workerCpu", strconv.Itoa(d.Get("cpu").(int)))
	urlValues.Add("workerRam", strconv.Itoa(d.Get("ram").(int)))
	urlValues.Add("workerDisk", strconv.Itoa(d.Get("disk").(int)))

	resp, err := c.DecortAPICall(ctx, "POST", K8sWgCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp)

	// This code is the supposed flow, but at the time of writing it's not yet implemented by the platfom

	//urlValues = &url.Values{}
	//urlValues.Add("auditId", strings.Trim(resp, `"`))

	//for {
	//resp, err := controller.decortAPICall("POST", AsyncTaskGetAPI, urlValues)
	//if err != nil {
	//return err
	//}

	//task := AsyncTask{}
	//if err := json.Unmarshal([]byte(resp), &task); err != nil {
	//return err
	//}
	//log.Debugf("resourceK8sCreate: workers group creating - %s", task.Stage)

	//if task.Completed {
	//if task.Error != "" {
	//return fmt.Errorf("cannot create workers group: %v", task.Error)
	//}

	//d.SetId(strconv.Itoa(int(task.Result)))
	//break
	//}

	//time.Sleep(time.Second * 5)
	//}

	return resourceK8sWgRead(ctx, d, m)
}

func resourceK8sWgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceK8sWgRead: called with %v", d.Id())

	k8s, err := utilityDataK8sCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	var id int
	if d.Id() != "" {
		id, err = strconv.Atoi(strings.Split(d.Id(), "#")[0])
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		id = d.Get("wg_id").(int)
	}

	curWg := K8SGroup{}
	for _, wg := range k8s.K8SGroups.Workers {
		if wg.ID == uint64(id) {
			curWg = wg
			break
		}
	}
	if curWg.ID == 0 {
		return diag.Errorf("Not found wg with id: %v in k8s cluster: %v", id, k8s.ID)
	}

	workersComputeList := make([]kvmvm.ComputeGetResp, 0, 0)
	for _, info := range curWg.DetailedInfo {
		compute, err := utilityComputeCheckPresence(ctx, d, m, info.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		workersComputeList = append(workersComputeList, *compute)
	}

	d.SetId(strings.Split(d.Id(), "#")[0])
	d.Set("k8s_id", k8s.ID)
	d.Set("wg_id", curWg.ID)
	flattenWgData(d, curWg, workersComputeList)

	return nil
}

func resourceK8sWgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceK8sWgUpdate: called with k8s id %d", d.Get("k8s_id").(int))

	c := m.(*controller.ControllerCfg)

	wg, err := utilityK8sWgCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))
	urlValues.Add("workersGroupId", d.Id())

	if newNum := d.Get("num").(int); uint64(newNum) > wg.Num {
		urlValues.Add("num", strconv.FormatUint(uint64(newNum)-wg.Num, 10))
		_, err := c.DecortAPICall(ctx, "POST", K8sWorkerAddAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		for i := int(wg.Num) - 1; i >= newNum; i-- {
			urlValues.Set("workerId", strconv.FormatUint(wg.DetailedInfo[i].ID, 10))
			_, err := c.DecortAPICall(ctx, "POST", K8sWorkerDeleteAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceK8sWgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceK8sWgDelete: called with k8s id %d", d.Get("k8s_id").(int))

	wg, err := utilityK8sWgCheckPresence(ctx, d, m)
	if wg == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))
	urlValues.Add("workersGroupId", strconv.FormatUint(wg.ID, 10))

	_, err = c.DecortAPICall(ctx, "POST", K8sWgDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceK8sWgSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"k8s_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of k8s instance.",
		},

		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the worker group.",
		},

		"num": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Number of worker nodes to create.",
		},

		"cpu": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     1,
			Description: "Worker node CPU count.",
		},

		"ram": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     1024,
			Description: "Worker node RAM in MB.",
		},

		"disk": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Worker node boot disk size. If unspecified or 0, size is defined by OS image size.",
		},
		"wg_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of k8s worker Group.",
		},
		"detailed_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: detailedInfoSchemaMake(),
			},
		},
		"labels": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"annotations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"taints": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func ResourceK8sWg() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceK8sWgCreate,
		ReadContext:   resourceK8sWgRead,
		UpdateContext: resourceK8sWgUpdate,
		DeleteContext: resourceK8sWgDelete,

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

		Schema: resourceK8sWgSchemaMake(),
	}
}

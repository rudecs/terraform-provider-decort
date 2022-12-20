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
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
	log "github.com/sirupsen/logrus"
)

func resourceK8sCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceK8sCreate: called with name %s, rg %d", d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("k8ciId", strconv.Itoa(d.Get("k8sci_id").(int)))
	urlValues.Add("workerGroupName", d.Get("wg_name").(string))

	var masterNode K8sNodeRecord
	if masters, ok := d.GetOk("masters"); ok {
		masterNode = parseNode(masters.([]interface{}))
	} else {
		masterNode = nodeMasterDefault()
	}
	urlValues.Add("masterNum", strconv.Itoa(masterNode.Num))
	urlValues.Add("masterCpu", strconv.Itoa(masterNode.Cpu))
	urlValues.Add("masterRam", strconv.Itoa(masterNode.Ram))
	urlValues.Add("masterDisk", strconv.Itoa(masterNode.Disk))

	var workerNode K8sNodeRecord
	if workers, ok := d.GetOk("workers"); ok {
		workerNode = parseNode(workers.([]interface{}))
	} else {
		workerNode = nodeWorkerDefault()
	}
	urlValues.Add("workerNum", strconv.Itoa(workerNode.Num))
	urlValues.Add("workerCpu", strconv.Itoa(workerNode.Cpu))
	urlValues.Add("workerRam", strconv.Itoa(workerNode.Ram))
	urlValues.Add("workerDisk", strconv.Itoa(workerNode.Disk))

	if withLB, ok := d.GetOk("with_lb"); ok {
		urlValues.Add("withLB", strconv.FormatBool(withLB.(bool)))
	} else {
		urlValues.Add("withLB", strconv.FormatBool(true))
	}

	if extNet, ok := d.GetOk("extnet_id"); ok {
		urlValues.Add("extnetId", strconv.Itoa(extNet.(int)))
	} else {
		urlValues.Add("extnetId", "0")
	}

	if desc, ok := d.GetOk("desc"); ok {
		urlValues.Add("desc", desc.(string))
	}

	resp, err := c.DecortAPICall(ctx, "POST", K8sCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	urlValues = &url.Values{}
	urlValues.Add("auditId", strings.Trim(resp, `"`))

	for {
		resp, err := c.DecortAPICall(ctx, "POST", AsyncTaskGetAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		task := AsyncTask{}
		if err := json.Unmarshal([]byte(resp), &task); err != nil {
			return diag.FromErr(err)
		}
		log.Debugf("resourceK8sCreate: instance creating - %s", task.Stage)

		if task.Completed {
			if task.Error != "" {
				return diag.FromErr(fmt.Errorf("cannot create k8s instance: %v", task.Error))
			}

			d.SetId(strconv.Itoa(int(task.Result)))
			break
		}

		time.Sleep(time.Second * 10)
	}

	return resourceK8sRead(ctx, d, m)
}

func resourceK8sRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//log.Debugf("resourceK8sRead: called with id %s, rg %d", d.Id(), d.Get("rg_id").(int))

	k8s, err := utilityDataK8sCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	k8sList, err := utilityK8sListCheckPresence(ctx, d, m, K8sListAPI)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	curK8s := K8SItem{}
	for _, k8sCluster := range k8sList {
		if k8sCluster.ID == k8s.ID {
			curK8s = k8sCluster
		}
	}
	if curK8s.ID == 0 {
		return diag.Errorf("Cluster with id %d not found", k8s.ID)
	}
	d.Set("vins_id", curK8s.VINSID)

	masterComputeList := make([]kvmvm.ComputeGetResp, 0, len(k8s.K8SGroups.Masters.DetailedInfo))
	workersComputeList := make([]kvmvm.ComputeGetResp, 0, len(k8s.K8SGroups.Workers))
	for _, masterNode := range k8s.K8SGroups.Masters.DetailedInfo {
		compute, err := utilityComputeCheckPresence(ctx, d, m, masterNode.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		masterComputeList = append(masterComputeList, *compute)
	}
	for _, worker := range k8s.K8SGroups.Workers {
		for _, info := range worker.DetailedInfo {
			compute, err := utilityComputeCheckPresence(ctx, d, m, info.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			workersComputeList = append(workersComputeList, *compute)
		}
	}

	flattenResourceK8s(d, *k8s, masterComputeList, workersComputeList)

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("lbId", strconv.FormatUint(k8s.LBID, 10))

	resp, err := c.DecortAPICall(ctx, "POST", LbGetAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	var lb LbRecord
	if err := json.Unmarshal([]byte(resp), &lb); err != nil {
		return diag.FromErr(err)
	}
	d.Set("extnet_id", lb.ExtNetID)
	d.Set("lb_ip", lb.PrimaryNode.FrontendIP)

	urlValues = &url.Values{}
	urlValues.Add("k8sId", d.Id())
	kubeconfig, err := c.DecortAPICall(ctx, "POST", K8sGetConfigAPI, urlValues)
	if err != nil {
		log.Warnf("could not get kubeconfig: %v", err)
	}
	d.Set("kubeconfig", kubeconfig)

	return nil
}

func resourceK8sUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceK8sUpdate: called with id %s, rg %d", d.Id(), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)

	if d.HasChange("name") {
		urlValues := &url.Values{}
		urlValues.Add("k8sId", d.Id())
		urlValues.Add("name", d.Get("name").(string))

		_, err := c.DecortAPICall(ctx, "POST", K8sUpdateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("workers") {
		k8s, err := utilityK8sCheckPresence(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}

		wg := k8s.K8SGroups.Workers[0]
		urlValues := &url.Values{}
		urlValues.Add("k8sId", d.Id())
		urlValues.Add("workersGroupId", strconv.FormatUint(wg.ID, 10))

		newWorkers := parseNode(d.Get("workers").([]interface{}))

		if uint64(newWorkers.Num) > wg.Num {
			urlValues.Add("num", strconv.FormatUint(uint64(newWorkers.Num)-wg.Num, 10))
			if _, err := c.DecortAPICall(ctx, "POST", K8sWorkerAddAPI, urlValues); err != nil {
				return diag.FromErr(err)
			}
		} else {
			for i := int(wg.Num) - 1; i >= newWorkers.Num; i-- {
				urlValues.Set("workerId", strconv.FormatUint(wg.DetailedInfo[i].ID, 10))
				if _, err := c.DecortAPICall(ctx, "POST", K8sWorkerDeleteAPI, urlValues); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return nil
}

func resourceK8sDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceK8sDelete: called with id %s, rg %d", d.Id(), d.Get("rg_id").(int))

	k8s, err := utilityK8sCheckPresence(ctx, d, m)
	if k8s == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", d.Id())
	urlValues.Add("permanently", "true")

	_, err = c.DecortAPICall(ctx, "POST", K8sDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceK8sSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the cluster.",
		},

		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Resource group ID that this instance belongs to.",
		},

		"k8sci_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the k8s catalog item to base this instance on.",
		},

		"wg_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name for first worker group created with cluster.",
		},

		"masters": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: mastersSchemaMake(),
			},
			Description: "Master node(s) configuration.",
		},

		"workers": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: workersSchemaMake(),
			},
			Description: "Worker node(s) configuration.",
		},

		"with_lb": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
			Description: "Create k8s with load balancer if true.",
		},

		"extnet_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "ID of the external network to connect workers to. If omitted network will be chosen by the platfom.",
		},

		"desc": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Text description of this instance.",
		},

		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: aclGroupSchemaMake(),
			},
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"account_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bservice_id": {
			Type:     schema.TypeInt,
			Computed: true,
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
		"k8s_ci_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"lb_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"lb_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP address of default load balancer.",
		},
		"rg_name": {
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
		"updated_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"default_wg_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of default workers group for this instace.",
		},

		"kubeconfig": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Kubeconfig for cluster access.",
		},
		"vins_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of default vins for this instace.",
		},
	}
}

func ResourceK8s() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceK8sCreate,
		ReadContext:   resourceK8sRead,
		UpdateContext: resourceK8sUpdate,
		DeleteContext: resourceK8sDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout30m,
			Read:    &constants.Timeout300s,
			Update:  &constants.Timeout300s,
			Delete:  &constants.Timeout300s,
			Default: &constants.Timeout300s,
		},

		Schema: resourceK8sSchemaMake(),
	}
}

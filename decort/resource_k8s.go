/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Petr Krutov, <petr.krutov@digitalenergy.online>

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
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceK8sCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceK8sCreate: called with name %s, rg %d", d.Get("name").(string), d.Get("rg_id").(int))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("k8ciId", strconv.Itoa(d.Get("k8sci_id").(int)))

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

	//TODO find a way to avoid hardcoding these values
	//if withLB, ok := d.GetOk("with_lb"); ok {
	//urlValues.Add("withLB", strconv.FormatBool(withLB.(bool)))
	//}
	urlValues.Add("withLB", strconv.FormatBool(true))

	//if extNet, ok := d.GetOk("extnet_id"); ok {
	//urlValues.Add("extnetId", strconv.Itoa(extNet.(int)))
	//}
	urlValues.Add("extnetId", strconv.Itoa(0))

	//if desc, ok := d.GetOk("desc"); ok {
	//urlValues.Add("desc", desc.(string))
	//}

	resp, err := controller.decortAPICall("POST", K8sCreateAPI, urlValues)
	if err != nil {
		return err
	}

	urlValues = &url.Values{}
	urlValues.Add("auditId", strings.Trim(resp, `"`))

	for {
		resp, err := controller.decortAPICall("POST", AsyncTaskGetAPI, urlValues)
		if err != nil {
			return err
		}

		task := AsyncTask{}
		if err := json.Unmarshal([]byte(resp), &task); err != nil {
			return err
		}
		log.Debugf("resourceK8sCreate: instance creating - %s", task.Stage)

		if task.Completed {
			if task.Error != "" {
				return fmt.Errorf("cannot create k8s instance: %v", task.Error)
			}

			d.SetId(strconv.Itoa(int(task.Result)))
			break
		}

		time.Sleep(time.Second * 10)
	}

	k8s, err := utilityK8sCheckPresence(d, m)
	if err != nil {
		return err
	}

	d.Set("default_wg_id", k8s.Groups.Workers[0].ID)

	return nil
}

func resourceK8sRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceK8sRead: called with id %s, rg %d", d.Id(), d.Get("rg_id").(int))

	k8s, err := utilityK8sCheckPresence(d, m)
	if k8s == nil {
		d.SetId("")
		return err
	}

	d.Set("name", k8s.Name)
	d.Set("rg_id", k8s.RgID)
	d.Set("k8sci_id", k8s.CI)
	d.Set("masters", nodeToResource(k8s.Groups.Masters))
	d.Set("workers", nodeToResource(k8s.Groups.Workers[0]))
	d.Set("default_wg_id", k8s.Groups.Workers[0].ID)

	return nil
}

func resourceK8sUpdate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceK8sUpdate: called with id %s, rg %d", d.Id(), d.Get("rg_id").(int))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", d.Id())
	urlValues.Add("name", d.Get("name").(string))

	_, err := controller.decortAPICall("POST", K8sUpdateAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceK8sDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceK8sDelete: called with id %s, rg %d", d.Id(), d.Get("rg_id").(int))

	k8s, err := utilityK8sCheckPresence(d, m)
	if k8s == nil {
		if err != nil {
			return err
		}
		return nil
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", d.Id())
	urlValues.Add("permanently", "true")

	_, err = controller.decortAPICall("POST", K8sDeleteAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceK8sExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourceK8sExists: called with id %s, rg %d", d.Id(), d.Get("rg_id").(int))

	k8s, err := utilityK8sCheckPresence(d, m)
	if k8s == nil {
		return false, err
	}

	return true, nil
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

		"masters": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: nodeK8sSubresourceSchemaMake(),
			},
			Description: "Master node(s) configuration.",
		},

		"workers": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: nodeK8sSubresourceSchemaMake(),
			},
			Description: "Worker node(s) configuration.",
		},

		//"with_lb": {
		//Type:        schema.TypeBool,
		//Optional:    true,
		//ForceNew:    true,
		//Default:     true,
		//Description: "Create k8s with load balancer if true.",
		//},

		//"extnet_id": {
		//Type:        schema.TypeInt,
		//Optional:    true,
		//ForceNew:    true,
		//Default:     0,
		//Description: "ID of the external network to connect workers to.",
		//},

		//"desc": {
		//Type:        schema.TypeString,
		//Optional:    true,
		//Description: "Text description of this instance.",
		//},

		"default_wg_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of default workers group for this instace.",
		},
	}
}

func resourceK8s() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceK8sCreate,
		Read:   resourceK8sRead,
		Update: resourceK8sUpdate,
		Delete: resourceK8sDelete,
		Exists: resourceK8sExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		//TODO timeouts

		Schema: resourceK8sSchemaMake(),
	}
}

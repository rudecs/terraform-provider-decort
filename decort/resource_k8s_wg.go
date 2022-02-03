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
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceK8sWgCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceK8sWgCreate: called with k8s id %d", d.Get("k8s_id").(int))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))
	urlValues.Add("name", uuid.New().String())
	urlValues.Add("workerNum", strconv.Itoa(d.Get("num").(int)))
	urlValues.Add("workerCpu", strconv.Itoa(d.Get("cpu").(int)))
	urlValues.Add("workerRam", strconv.Itoa(d.Get("ram").(int)))
	urlValues.Add("workerDisk", strconv.Itoa(d.Get("disk").(int)))

	resp, err := controller.decortAPICall("POST", K8sWgCreateAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId(resp)
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

	return nil
}

func resourceK8sWgRead(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceK8sWgRead: called with k8s id %d", d.Get("k8s_id").(int))

	wg, err := utilityK8sWgCheckPresence(d, m)
	if wg == nil {
		d.SetId("")
		return err
	}

	d.Set("num", wg.Num)
	d.Set("cpu", wg.Cpu)
	d.Set("ram", wg.Ram)
	d.Set("disk", wg.Disk)

	return nil
}

func resourceK8sWgDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceK8sWgDelete: called with k8s id %d", d.Get("k8s_id").(int))

	wg, err := utilityK8sWgCheckPresence(d, m)
	if wg == nil {
		if err != nil {
			return err
		}
		return nil
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))
	urlValues.Add("workersGroupId", strconv.Itoa(wg.ID))

	_, err = controller.decortAPICall("POST", K8sWgDeleteAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceK8sWgExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Debugf("resourceK8sWgExists: called with k8s id %d", d.Get("k8s_id").(int))

	wg, err := utilityK8sWgCheckPresence(d, m)
	if wg == nil {
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func resourceK8sWgSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"k8s_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of k8s instance.",
		},

		//Unused but required by creation API. Sending generated UUID each time
		//"name": {
		//Type:        schema.TypeString,
		//Required:    true,
		//ForceNew:    true,
		//Description: "Name of the worker group.",
		//},

		"num": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
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
	}
}

func resourceK8sWg() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceK8sWgCreate,
		Read:   resourceK8sWgRead,
		Delete: resourceK8sWgDelete,
		Exists: resourceK8sWgExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		//TODO timeouts

		Schema: resourceK8sWgSchemaMake(),
	}
}

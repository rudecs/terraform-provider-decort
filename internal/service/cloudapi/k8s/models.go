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

package k8s

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type K8sNodeRecord struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Disk         int    `json:"disk"`
	Cpu          int    `json:"cpu"`
	Num          int    `json:"num"`
	Ram          int    `json:"ram"`
	DetailedInfo []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"detailedInfo"`
}

//K8sRecord represents k8s instance
type K8sRecord struct {
	AccountID   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	CI          int    `json:"ciId"`
	ID          int    `json:"id"`
	Groups      struct {
		Masters K8sNodeRecord   `json:"masters"`
		Workers []K8sNodeRecord `json:"workers"`
	} `json:"k8sGroups"`
	LbID   int    `json:"lbId"`
	Name   string `json:"name"`
	RgID   int    `json:"rgId"`
	RgName string `json:"rgName"`
}

//LbRecord represents load balancer instance
type LbRecord struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	RgID        int    `json:"rgId"`
	VinsID      int    `json:"vinsId"`
	ExtNetID    int    `json:"extnetId"`
	PrimaryNode struct {
		BackendIP  string `json:"backendIp"`
		ComputeID  int    `json:"computeId"`
		FrontendIP string `json:"frontendIp"`
		NetworkID  int    `json:"networkId"`
	} `json:"primaryNode"`
}

//Blasphemous workaround for parsing Result value
type TaskResult int

func (r *TaskResult) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		b := b[1 : len(b)-1]
		if len(b) == 0 {
			*r = 0
			return nil
		}
		n, err := strconv.Atoi(string(b))
		if err != nil {
			return err
		}
		*r = TaskResult(n)
	} else if b[0] == '[' {
		res := []interface{}{}
		if err := json.Unmarshal(b, &res); err != nil {
			return err
		}
		if n, ok := res[0].(float64); ok {
			*r = TaskResult(n)
		} else {
			return fmt.Errorf("could not unmarshal %v into int", res[0])
		}
	}

	return nil
}

//AsyncTask represents a long task completion status
type AsyncTask struct {
	AuditID     string     `json:"auditId"`
	Completed   bool       `json:"completed"`
	Error       string     `json:"error"`
	Log         []string   `json:"log"`
	Result      TaskResult `json:"result"`
	Stage       string     `json:"stage"`
	Status      string     `json:"status"`
	UpdateTime  uint64     `json:"updateTime"`
	UpdatedTime uint64     `json:"updatedTime"`
}

type SshKeyConfig struct {
	User      string
	SshKey    string
	UserShell string
}

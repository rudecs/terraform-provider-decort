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
	VinsID int    `json:"vinsId"`
}

type K8sRecordList []K8sRecord

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

//FromSDK
type K8SGroup struct {
	Annotations  []string         `json:"annotations"`
	CPU          uint64           `json:"cpu"`
	DetailedInfo DetailedInfoList `json:"detailedInfo"`
	Disk         uint64           `json:"disk"`
	GUID         string           `json:"guid"`
	ID           uint64           `json:"id"`
	Labels       []string         `json:"labels"`
	Name         string           `json:"name"`
	Num          uint64           `json:"num"`
	RAM          uint64           `json:"ram"`
	Taints       []string         `json:"taints"`
}

type K8SGroupList []K8SGroup

type DetailedInfo struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	TechStatus string `json:"techStatus"`
}

type DetailedInfoList []DetailedInfo

type K8SRecord struct {
	ACL         ACLGroup  `json:"ACL"`
	AccountID   uint64    `json:"accountId"`
	AccountName string    `json:"accountName"`
	BServiceID  uint64    `json:"bserviceId"`
	CIID        uint64    `json:"ciId"`
	CreatedBy   string    `json:"createdBy"`
	CreatedTime uint64    `json:"createdTime"`
	DeletedBy   string    `json:"deletedBy"`
	DeletedTime uint64    `json:"deletedTime"`
	ID          uint64    `json:"id"`
	K8CIName    string    `json:"k8ciName"`
	K8SGroups   K8SGroups `json:"k8sGroups"`
	LBID        uint64    `json:"lbId"`
	Name        string    `json:"name"`
	RGID        uint64    `json:"rgId"`
	RGName      string    `json:"rgName"`
	Status      string    `json:"status"`
	TechStatus  string    `json:"techStatus"`
	UpdatedBy   string    `json:"updatedBy"`
	UpdatedTime uint64    `json:"updatedTime"`
}

type K8SRecordList []K8SRecord

type K8SGroups struct {
	Masters MasterGroup  `json:"masters"`
	Workers K8SGroupList `json:"workers"`
}

type MasterGroup struct {
	CPU          uint64           `json:"cpu"`
	DetailedInfo DetailedInfoList `json:"detailedInfo"`
	Disk         uint64           `json:"disk"`
	ID           uint64           `json:"id"`
	Name         string           `json:"name"`
	Num          uint64           `json:"num"`
	RAM          uint64           `json:"ram"`
}

type ACLGroup struct {
	AccountACL ACLList `json:"accountAcl"`
	K8SACL     ACLList `json:"k8sAcl"`
	RGACL      ACLList `json:"rgAcl"`
}

type ACL struct {
	Explicit    bool   `json:"explicit"`
	GUID        string `json:"guid"`
	Right       string `json:"right"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	UserGroupID string `json:"userGroupId"`
}

type ACLList []ACL

type K8SItem struct {
	AccountID      uint64         `json:"accountId"`
	AccountName    string         `json:"accountName"`
	ACL            []interface{}  `json:"acl"`
	BServiceID     uint64         `json:"bserviceId"`
	CIID           uint64         `json:"ciId"`
	Config         interface{}    `json:"config"`
	CreatedBy      string         `json:"createdBy"`
	CreatedTime    uint64         `json:"createdTime"`
	DeletedBy      string         `json:"deletedBy"`
	DeletedTime    uint64         `json:"deletedTime"`
	Description    string         `json:"desc"`
	ExtNetID       uint64         `json:"extnetId"`
	GID            uint64         `json:"gid"`
	GUID           uint64         `json:"guid"`
	ID             uint64         `json:"id"`
	LBID           uint64         `json:"lbId"`
	Milestones     uint64         `json:"milestones"`
	Name           string         `json:"name"`
	RGID           uint64         `json:"rgId"`
	RGName         string         `json:"rgName"`
	ServiceAccount ServiceAccount `json:"serviceAccount"`
	Status         string         `json:"status"`
	TechStatus     string         `json:"techStatus"`
	UpdatedBy      string         `json:"updatedBy"`
	UpdatedTime    uint64         `json:"updatedTime"`
	VINSID         uint64         `json:"vinsId"`
	WorkersGroup   K8SGroupList   `json:"workersGroups"`
}

type ServiceAccount struct {
	GUID     string `json:"guid"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type K8SList []K8SItem

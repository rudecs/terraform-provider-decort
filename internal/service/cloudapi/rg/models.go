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

type ResourceLimits struct {
	CUC      float64 `json:"CU_C"`
	CUD      float64 `json:"CU_D"`
	CUI      float64 `json:"CU_I"`
	CUM      float64 `json:"CU_M"`
	CUNP     float64 `json:"CU_NP"`
	GpuUnits float64 `json:"gpu_units"`
}

type ResgroupRecord struct {
	ACLs             []AccountAclRecord `json:"acl"`
	AccountID        int                `json:"accountId"`
	AccountName      string             `json:"accountName"`
	CreatedBy        string             `json:"createdBy"`
	CreatedTime      uint64             `json:"createdTime"`
	DefaultNetID     int                `json:"def_net_id"`
	DefaultNetType   string             `json:"def_net_type"`
	DeletedBy        string             `json:"deletedBy"`
	DeletedTime      int                `json:"deletedTime"`
	Decsription      string             `json:"desc"`
	GridID           int                `json:"gid"`
	GUID             int                `json:"guid"`
	ID               uint               `json:"id"`
	LockStatus       string             `json:"lockStatus"`
	Milestones       int                `json:"milestones"`
	Name             string             `json:"name"`
	RegisterComputes bool               `json:"registerComputes"`
	ResourceLimits   ResourceLimits     `json:"resourceLimits"`
	Secret           string             `json:"secret"`
	Status           string             `json:"status"`
	UpdatedBy        string             `json:"updatedBy"`
	UpdatedTime      uint64             `json:"updatedTime"`
	Vins             []int              `json:"vins"`
	Computes         []int              `json:"vms"`
}

type ResgroupListResp []ResgroupRecord

type ResgroupUpdateParam struct {
	RgId       int    `json:"rgId"`
	Name       string `json:"name"`
	Desc       string `json:"decs"`
	Ram        int    `json:"maxMemoryCapacity"`
	Disk       int    `json:"maxVDiskCapacity"`
	Cpu        int    `json:"maxCPUCapacity"`
	NetTraffic int    `json:"maxNetworkPeerTransfer"`
	Reason     string `json:"reason"`
}

type AccountAclRecord struct {
	IsExplicit   bool   `json:"explicit"`
	Guid         string `json:"guid"`
	Rights       string `json:"right"`
	Status       string `json:"status"`
	Type         string `json:"type"`
	UgroupID     string `json:"userGroupId"`
	CanBeDeleted bool   `json:"canBeDeleted"`
}

type ResgroupGetResp struct {
	Resources Resources       `json:"Resources"`
	ACLs      []UserAclRecord `json:"ACLs"`
	//Usage          UsageRecord     `json:"Resources"`
	AccountID      int         `json:"accountId"`
	AccountName    string      `json:"accountName"`
	GridID         int         `json:"gid"`
	CreatedBy      string      `json:"createdBy"`
	CreatedTime    uint64      `json:"createdTime"`
	DefaultNetID   int         `json:"def_net_id"`
	DefaultNetType string      `json:"def_net_type"`
	DeletedBy      string      `json:"deletedBy"`
	DeletedTime    uint64      `json:"deletedTime"`
	Desc           string      `json:"desc"`
	ID             uint        `json:"id"`
	LockStatus     string      `json:"lockStatus"`
	Name           string      `json:"name"`
	Quota          QuotaRecord `json:"resourceLimits"`
	Status         string      `json:"status"`
	UpdatedBy      string      `json:"updatedBy"`
	UpdatedTime    uint64      `json:"updatedTime"`
	Vins           []int       `json:"vins"`
	Computes       []int       `json:"vms"`

	Ignored map[string]interface{} `json:"-"`
}

type UserAclRecord struct {
	IsExplicit bool   `json:"explicit"`
	Rights     string `json:"right"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	UgroupID   string `json:"userGroupId"`
	// CanBeDeleted bool      `json:"canBeDeleted"`
}

type QuotaRecord struct { // this is how quota is reported by /api/.../rg/get
	Cpu        int     `json:"CU_C"`      // CPU count in pcs
	Ram        float64 `json:"CU_M"`      // RAM volume in MB, it is STILL reported as FLOAT
	Disk       int     `json:"CU_D"`      // Disk capacity in GB
	ExtIPs     int     `json:"CU_I"`      // Ext IPs count
	ExtTraffic int     `json:"CU_NP"`     // Ext network traffic
	GpuUnits   int     `json:"gpu_units"` // GPU count
}

type ResourceRecord struct { // this is how actual usage is reported by /api/.../rg/get
	Cpu        int     `json:"cpu"`
	Disk       float64 `json:"disksize"`
	ExtIPs     int     `json:"extips"`
	ExtTraffic int     `json:"exttraffic"`
	Gpu        int     `json:"gpu"`
	Ram        int     `json:"ram"`
}

type UsageRecord struct {
	Current  ResourceRecord `json:"Current"`
	Reserved ResourceRecord `json:"Reserved"`
}

type ResourceSep struct {
	DiskSize    float64 `json:"disksize"`
	DiskSizeMax int     `json:"disksizemax"`
}

type Resource struct {
	CPU        int                               `json:"cpu"`
	Disksize   float64                           `json:"disksize"`
	Extips     int                               `json:"extips"`
	Exttraffic int                               `json:"exttraffic"`
	GPU        int                               `json:"gpu"`
	RAM        int                               `json:"ram"`
	SEPs       map[string]map[string]ResourceSep `json:"seps"`
}

type Resources struct {
	Current  Resource `json:"Current"`
	Reserved Resource `json:"Reserved"`
}

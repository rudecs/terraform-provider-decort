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

package bservice

///Structs

type BasicServiceCompute struct {
	CompGroupId   int    `json:"compgroupId"`
	CompGroupName string `json:"compgroupName"`
	CompGroupRole string `json:"compgroupRole"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
}

type BasicServiceComputes []BasicServiceCompute

type BasicServiceSnapshot struct {
	GUID      string `json:"guid"`
	Label     string `json:"label"`
	Timestamp int    `json:"timestamp"`
	Valid     bool   `json:"valid"`
}

type BasicServiceSnapshots []BasicServiceSnapshot

type BasicService struct {
	AccountId   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	BaseDomain  string `json:"baseDomain"`

	CreatedBy     string `json:"createdBy"`
	CreatedTime   int    `json:"createdTime"`
	DeletedBy     string `json:"deletedBy"`
	DeletedTime   int    `json:"deletedTime"`
	GID           int    `json:"gid"`
	Groups        []int  `json:"groups"`
	GUID          int    `json:"guid"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ParentSrvId   int    `json:"parentSrvId"`
	ParentSrvType string `json:"parentSrvType"`
	RGID          int    `json:"rgId"`
	RGName        string `json:"rgName"`
	SSHUser       string `json:"sshUser"`
	Status        string `json:"status"`
	TechStatus    string `json:"techStatus"`
	UpdatedBy     string `json:"updatedBy"`
	UpdatedTime   int    `json:"updatedTime"`
	UserManaged   bool   `json:"userManaged"`
}

type BasicServiceList []BasicService

type BasicServiceExtend struct {
	BasicService
	Computes   BasicServiceComputes  `json:"computes"`
	CPUTotal   int                   `json:"cpuTotal"`
	DiskTotal  int                   `json:"diskTotal"`
	GroupsName []string              `json:"groupsName"`
	Milestones int                   `json:"milestones"`
	RamTotal   int                   `json:"ramTotal"`
	Snapshots  BasicServiceSnapshots `json:"snapshots"`
	SSHKey     string                `json:"sshKey"`
}

type BasicServiceGroupOSUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BasicServiceGroupOSUsers []BasicServiceGroupOSUser

type BasicServicceGroupCompute struct {
	ID         int                      `json:"id"`
	IPAdresses []string                 `json:"ipAddresses"`
	Name       string                   `json:"name"`
	OSUsers    BasicServiceGroupOSUsers `json:"osUsers"`
}

type BasicServiceGroupComputes []BasicServicceGroupCompute

type BasicServiceGroup struct {
	AccountId    int                       `json:"accountId"`
	AccountName  string                    `json:"accountName"`
	Computes     BasicServiceGroupComputes `json:"computes"`
	Consistency  bool                      `json:"consistency"`
	CPU          int                       `json:"cpu"`
	CreatedBy    string                    `json:"createdBy"`
	CreatedTime  int                       `json:"createdTime"`
	DeletedBy    string                    `json:"deletedBy"`
	DeletedTime  int                       `json:"deletedTime"`
	Disk         int                       `json:"disk"`
	Driver       string                    `json:"driver"`
	Extnets      []int                     `json:"extnets"`
	GID          int                       `json:"gid"`
	GUID         int                       `json:"guid"`
	ID           int                       `json:"id"`
	ImageId      int                       `json:"imageId"`
	Milestones   int                       `json:"milestones"`
	Name         string                    `json:"name"`
	Parents      []int                     `json:"parents"`
	RAM          int                       `json:"ram"`
	RGID         int                       `json:"rgId"`
	RGName       string                    `json:"rgName"`
	Role         string                    `json:"role"`
	SepId        int                       `json:"sepId"`
	SeqNo        int                       `json:"seqNo"`
	ServiceId    int                       `json:"serviceId"`
	Status       string                    `json:"status"`
	TechStatus   string                    `json:"techStatus"`
	TimeoutStart int                       `json:"timeoutStart"`
	UpdatedBy    string                    `json:"updatedBy"`
	UpdatedTime  int                       `json:"updatedTime"`
	Vinses       []int                     `json:"vinses"`
}

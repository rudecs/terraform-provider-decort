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

package disks

type Disk struct {
	Acl                 map[string]interface{} `json:"acl"`
	AccountID           int                    `json:"accountId"`
	AccountName         string                 `json:"accountName"`
	BootPartition       int                    `json:"bootPartition"`
	Computes            map[string]string      `json:"computes"`
	CreatedTime         uint64                 `json:"creationTime"`
	DeletedTime         uint64                 `json:"deletionTime"`
	DeviceName          string                 `json:"devicename"`
	Desc                string                 `json:"desc"`
	DestructionTime     uint64                 `json:"destructionTime"`
	DiskPath            string                 `json:"diskPath"`
	GridID              int                    `json:"gid"`
	GUID                int                    `json:"guid"`
	ID                  uint                   `json:"id"`
	ImageID             int                    `json:"imageId"`
	Images              []int                  `json:"images"`
	IOTune              IOTune                 `json:"iotune"`
	IQN                 string                 `json:"iqn"`
	Login               string                 `json:"login"`
	Name                string                 `json:"name"`
	MachineId           int                    `json:"machineId"`
	MachineName         string                 `json:"machineName"`
	Milestones          uint64                 `json:"milestones"`
	Order               int                    `json:"order"`
	Params              string                 `json:"params"`
	Passwd              string                 `json:"passwd"`
	ParentId            int                    `json:"parentId"`
	PciSlot             int                    `json:"pciSlot"`
	Pool                string                 `json:"pool"`
	PresentTo           []int                  `json:"presentTo"`
	PurgeTime           uint64                 `json:"purgeTime"`
	PurgeAttempts       uint64                 `json:"purgeAttempts"`
	RealityDeviceNumber int                    `json:"realityDeviceNumber"`
	ReferenceId         string                 `json:"referenceId"`
	ResID               string                 `json:"resId"`
	ResName             string                 `json:"resName"`
	Role                string                 `json:"role"`
	SepType             string                 `json:"sepType"`
	Shareable           bool                   `json:"shareable"`
	SepID               int                    `json:"sepId"` // NOTE: absent from compute/get output
	SizeMax             int                    `json:"sizeMax"`
	SizeUsed            float64                `json:"sizeUsed"` // sum over all snapshots of this disk to report total consumed space
	Snapshots           []Snapshot             `json:"snapshots"`
	Status              string                 `json:"status"`
	TechStatus          string                 `json:"techStatus"`
	Type                string                 `json:"type"`
	UpdateBy            uint64                 `json:"updateBy"`
	VMID                int                    `json:"vmid"`
}

type Snapshot struct {
	Guid        string `json:"guid"`
	Label       string `json:"label"`
	ResId       string `json:"resId"`
	SnapSetGuid string `json:"snapSetGuid"`
	SnapSetTime uint64 `json:"snapSetTime"`
	TimeStamp   uint64 `json:"timestamp"`
}

type SnapshotList []Snapshot

type DisksList []Disk

type IOTune struct {
	ReadBytesSec     int `json:"read_bytes_sec"`
	ReadBytesSecMax  int `json:"read_bytes_sec_max"`
	ReadIopsSec      int `json:"read_iops_sec"`
	ReadIopsSecMax   int `json:"read_iops_sec_max"`
	SizeIopsSec      int `json:"size_iops_sec"`
	TotalBytesSec    int `json:"total_bytes_sec"`
	TotalBytesSecMax int `json:"total_bytes_sec_max"`
	TotalIopsSec     int `json:"total_iops_sec"`
	TotalIopsSecMax  int `json:"total_iops_sec_max"`
	WriteBytesSec    int `json:"write_bytes_sec"`
	WriteBytesSecMax int `json:"write_bytes_sec_max"`
	WriteIopsSec     int `json:"write_iops_sec"`
	WriteIopsSecMax  int `json:"write_iops_sec_max"`
}

type Pool struct {
	Name  string   `json:"name"`
	Types []string `json:"types"`
}

type PoolList []Pool

type TypeDetailed struct {
	Pools []Pool `json:"pools"`
	SepID int    `json:"sepId"`
}

type TypesDetailedList []TypeDetailed

type TypesList []string

type Unattached struct {
	Ckey                string                 `json:"_ckey"`
	Meta                []interface{}          `json:"_meta"`
	AccountID           int                    `json:"accountId"`
	AccountName         string                 `json:"accountName"`
	Acl                 map[string]interface{} `json:"acl"`
	BootPartition       int                    `json:"bootPartition"`
	CreatedTime         int                    `json:"createdTime"`
	DeletedTime         int                    `json:"deletedTime"`
	Desc                string                 `json:"desc"`
	DestructionTime     int                    `json:"destructionTime"`
	DiskPath            string                 `json:"diskPath"`
	GridID              int                    `json:"gid"`
	GUID                int                    `json:"guid"`
	ID                  int                    `json:"id"`
	ImageID             int                    `json:"imageId"`
	Images              []int                  `json:"images"`
	IOTune              IOTune                 `json:"iotune"`
	IQN                 string                 `json:"iqn"`
	Login               string                 `json:"login"`
	Milestones          int                    `json:"milestones"`
	Name                string                 `json:"name"`
	Order               int                    `json:"order"`
	Params              string                 `json:"params"`
	ParentID            int                    `json:"parentId"`
	Passwd              string                 `json:"passwd"`
	PciSlot             int                    `json:"pciSlot"`
	Pool                string                 `json:"pool"`
	PurgeAttempts       int                    `json:"purgeAttempts"`
	PurgeTime           int                    `json:"purgeTime"`
	RealityDeviceNumber int                    `json:"realityDeviceNumber"`
	ReferenceID         string                 `json:"referenceId"`
	ResID               string                 `json:"resId"`
	ResName             string                 `json:"resName"`
	Role                string                 `json:"role"`
	SepID               int                    `json:"sepId"`
	SizeMax             int                    `json:"sizeMax"`
	SizeUsed            float64                `json:"sizeUsed"`
	Snapshots           []Snapshot             `json:"snapshots"`
	Status              string                 `json:"status"`
	TechStatus          string                 `json:"techStatus"`
	Type                string                 `json:"type"`
	VMID                int                    `json:"vmid"`
}

type UnattachedList []Unattached

type Pair struct {
	intPort      int
	extPortStart int
}

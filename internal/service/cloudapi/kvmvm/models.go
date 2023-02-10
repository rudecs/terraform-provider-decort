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

package kvmvm

type DiskRecord struct {
	Acl                 map[string]interface{} `json:"acl"`
	AccountID           int                    `json:"accountId"`
	AccountName         string                 `json:"accountName"`
	BootPartition       int                    `json:"bootPartition"`
	CreatedTime         uint64                 `json:"creationTime"`
	ComputeID           int                    `json:"computeId"`
	ComputeName         string                 `json:"computeName"`
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
	IOTune              map[string]interface{} `json:"iotune"`
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
	PurgeTime           uint64                 `json:"purgeTime"`
	PurgeAttempts       uint64                 `json:"purgeAttempts"`
	RealityDeviceNumber int                    `json:"realityDeviceNumber"`
	ReferenceId         string                 `json:"referenceId"`
	ResID               string                 `json:"resId"`
	ResName             string                 `json:"resName"`
	Role                string                 `json:"role"`
	SepType             string                 `json:"sepType"`
	SepID               int                    `json:"sepId"` // NOTE: absent from compute/get output
	Shareable           bool                   `json:"shareable"`
	SizeMax             int                    `json:"sizeMax"`
	SizeUsed            float64                `json:"sizeUsed"` // sum over all snapshots of this disk to report total consumed space
	Snapshots           []SnapshotRecord       `json:"snapshots"`
	Status              string                 `json:"status"`
	TechStatus          string                 `json:"techStatus"`
	Type                string                 `json:"type"`
	UpdateBy            uint64                 `json:"updateBy"`
	VMID                int                    `json:"vmid"`
}

type InterfaceRecord struct {
	ConnID    int                `json:"connId"`   // This is VLAN ID or VxLAN ID, depending on ConnType
	ConnType  string             `json:"connType"` // Either "VLAN" or "VXLAN" tag
	DefaultGW string             `json:"defGw"`
	Guid      string             `json:"guid"`
	IPAddress string             `json:"ipAddress"` // without trailing network mask, i.e. "192.168.1.3"
	MAC       string             `json:"mac"`
	Name      string             `json:"name"`
	NetID     int                `json:"netId"` // This is either ExtNet ID or ViNS ID, depending on NetType
	NetMask   int                `json:"netMask"`
	NetType   string             `json:"netType"` // Either "EXTNET" or "VINS" tag
	PciSlot   int                `json:"pciSlot"`
	Target    string             `json:"target"`
	Type      string             `json:"type"`
	VNFs      []int              `json:"vnfs"`
	QOS       InterfaceQosRecord `json:"qos"`
}

type InterfaceQosRecord struct {
	ERate   int    `json:"eRate"`
	Guid    string `json:"guid"`
	InBurst int    `json:"inBurst"`
	InRate  int    `json:"inRate"`
}

type SnapshotRecord struct {
	Guid        string `json:"guid"`
	Label       string `json:"label"`
	ResId       string `json:"resId"`
	SnapSetGuid string `json:"snapSetGuid"`
	SnapSetTime uint64 `json:"snapSetTime"`
	TimeStamp   uint64 `json:"timestamp"`
}

type SnapshotRecordList []SnapshotRecord

type ComputeGetResp struct {
	// ACLs `json:"ACL"` - it is a dictionary, special parsing required
	AccountID          int               `json:"accountId"`
	AccountName        string            `json:"accountName"`
	Arch               string            `json:"arch"`
	BootDiskSize       int               `json:"bootdiskSize"`
	CloneReference     int               `json:"cloneReference"`
	Clones             []int             `json:"clones"`
	Cpu                int               `json:"cpus"`
	Desc               string            `json:"desc"`
	Disks              []DiskRecord      `json:"disks"`
	Driver             string            `json:"driver"`
	GridID             int               `json:"gid"`
	ID                 uint              `json:"id"`
	ImageID            int               `json:"imageId"`
	ImageName          string            `json:"imageName"`
	Interfaces         []InterfaceRecord `json:"interfaces"`
	LockStatus         string            `json:"lockStatus"`
	ManagerID          int               `json:"managerId"`
	ManagerType        string            `json:"manageType"`
	Name               string            `json:"name"`
	NatableVinsID      int               `json:"natableVinsId"`
	NatableVinsIP      string            `json:"natableVinsIp"`
	NatableVinsName    string            `json:"natableVinsName"`
	NatableVinsNet     string            `json:"natableVinsNetwork"`
	NatableVinsNetName string            `json:"natableVinsNetworkName"`
	OsUsers            []OsUserRecord    `json:"osUsers"`
	Ram                int               `json:"ram"`
	RgID               int               `json:"rgId"`
	RgName             string            `json:"rgName"`
	SnapSets           []SnapSetRecord   `json:"snapSets"`
	Status             string            `json:"status"`
	// Tags               []string          `json:"tags"` // Tags were reworked since DECORT 3.7.1
	TechStatus     string `json:"techStatus"`
	TotalDiskSize  int    `json:"totalDiskSize"`
	UpdatedBy      string `json:"updatedBy"`
	UpdateTime     uint64 `json:"updateTime"`
	UserManaged    bool   `json:"userManaged"`
	Vgpus          []int  `json:"vgpus"`
	VinsConnected  int    `json:"vinsConnected"`
	VirtualImageID int    `json:"virtualImageId"`
}

type OsUserRecord struct {
	Guid     string `json:"guid"`
	Login    string `json:"login"`
	Password string `json:"password"`
	PubKey   string `json:"pubkey"`
}

type SnapSetRecord struct {
	Disks     []int  `json:"disks"`
	Guid      string `json:"guid"`
	Label     string `json:"label"`
	TimeStamp uint64 `json:"timestamp"`
}

type ComputeBriefRecord struct { // this is a brief compute specifiaction as returned by API rg/listComputes
	// we do not even include here all fields as returned by this API, but only the most important that
	// are really necessary to identify and distinguish computes
	AccountID   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	Name        string `json:"name"`
	ID          uint   `json:"id"`
	RgID        int    `json:"rgId"`
	RgName      string `json:"rgName"`
	Status      string `json:"status"`
	TechStatus  string `json:"techStatus"`
}

type RgListComputesResp []ComputeBriefRecord

/*
Copyright (c) 2019-2020 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

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

	"time"
)

//
// timeouts for API calls from CRUD functions of Terraform plugin
var Timeout30s = time.Second * 30
var Timeout60s = time.Second * 60
var Timeout180s = time.Second * 180

//
// structures related to /cloudapi/rg/list API
//
type UserAclRecord struct {
	IsExplicit bool        `json:"explicit"`
	Rights string          `json:"right"`
	Status string          `json:"status"`
	Type string            `json:"type"`
	UgroupID string        `json:"userGroupId"`
	// CanBeDeleted bool      `json:"canBeDeleted"`
}

type AccountAclRecord struct {
	IsExplicit bool        `json:"explicit"`
	Guid string            `json:"guid"`
	Rights string          `json:"right"`
	Status string          `json:"status"`
	Type string            `json:"type"`
	UgroupID string        `json:"userGroupId"`
}

type ResgroupRecord struct {
	ACLs []UserAclRecord   `json:"ACLs"`
	Owner AccountAclRecord `json:"accountAcl"`
	AccountID int           `json:"accountId"`
	AccountName string      `json:"accountName"`
	CreatedBy string       `json:"createdBy"`
	CreatedTime uint64     `json:"createdTime"`
	DefaultNetID int       `json:"def_net_id"`
	DefaultNetType string  `json:"def_net_type"`
	Decsription string     `json:"desc"`
	GridID int             `json:"gid"`
	ID uint                `json:"id"`
	LockStatus string      `json:"lockStatus"`
	Name string            `json:"name"`
	Status string          `json:"status"`
	UpdatedBy string       `json:"updatedBy"`
	UpdatedTime uint64     `json:"updatedTime"`
	Vins []int             `json:"vins"`
	Computes []int         `json:"vms"`
}

const ResgroupListAPI = "/restmachine/cloudapi/rg/list"
type ResgroupListResp []ResgroupRecord

//
// structures related to /cloudapi/rg/create API call
//
const ResgroupCreateAPI= "/restmachine/cloudapi/rg/create"

//
// structures related to /cloudapi/rg/update API call
//
const ResgroupUpdateAPI= "/restmachine/cloudapi/rg/update"
type ResgroupUpdateParam struct {
	RgId int               `json:"rgId"`
	Name string            `json:"name"`
	Desc string            `json:"decs"`
	Ram int                `json:"maxMemoryCapacity"`
	Disk int               `json:"maxVDiskCapacity"`
	Cpu int                `json:"maxCPUCapacity"`
	NetTraffic int         `json:"maxNetworkPeerTransfer"`
	Reason string          `json:"reason"`
} 

//
// structures related to /cloudapi/rg/get API call
//
type ResourceRecord struct {
	Cpu int                `json:"cpu"`
	Disk int               `json:"disksize"`
	ExtIPs int             `json:"extips"`
	ExtTraffic int         `json:"exttraffic"`
	Gpu int                `json:"gpu"`
	Ram int                `json:"ram"`
}

type UsageRecord struct {
	Current ResourceRecord    `json:"Current"`
	Reserved ResourceRecord   `json:"Reserved"`
}

const ResgroupGetAPI= "/restmachine/cloudapi/rg/get"
type ResgroupGetResp struct {
	ACLs []UserAclRecord   `json:"ACLs"`
	Usage UsageRecord      `json:"Resources"`
	AccountID int          `json:"accountId"`
	AccountName string     `json:"accountName"`

	CreatedBy string       `json:"createdBy"`
	CreatedTime uint64     `json:"createdTime"`
	DefaultNetID int       `json:"def_net_id"`
	DefaultNetType string  `json:"def_net_type"`
	DeletedBy string       `json:"deletedBy"`
	DeletedTime uint64     `json:"deletedTime"`
	Decsription string     `json:"desc"`
	ID uint                `json:"id"`
	LockStatus string      `json:"lockStatus"`
	Name string            `json:"name"`
	Quotas QuotaRecord     `json:"resourceLimits"`
	Status string          `json:"status"`
	UpdatedBy string       `json:"updatedBy"`
	UpdatedTime uint64     `json:"updatedTime"`
	Vins []int             `json:"vins"`
	Computes []int         `json:"vms"`

	Ignored map[string]interface{} `json:"-"`
}

// 
// structures related to /cloudapi/rg/update API
//
const ResgroupUpdateAPI = "/restmachine/cloudapi/rg/update"
type ResgroupUpdateParam struct {
	ID uint                `json:"rgId"`
	Name string            `json:"name"`
	Decsription string     `json:"desc"`
	Cpu int                `json:"maxCPUCapacity"`
	Ram int                `json:"maxMemoryCapacity"`
	Disk int               `json:"maxVDiskCapacity"`
	NetTraffic int         `json:"maxNetworkPeerTransfer"`
	ExtIPs int             `json:"maxNumPublicIP"`
	Reason string          `json:"reason"`
}

// 
// structures related to /cloudapi/rg/delete API
//
const ResgroupDeleteAPI = "/restmachine/cloudapi/rg/delete"

//
// structures related to /cloudapi/kvmXXX/create APIs
//
const KvmX86CreateAPI = "/restmachine/cloudapi/kvmx86/create"
const KvmPPCCreateAPI = "/restmachine/cloudapi/kvmppc/create"
type KvmXXXCreateParam struct { // this is unified structure for both x86 and PPC based VMs creation
	RgID uint              `json:"rgId"`
	Name string            `json:"name"`
	Cpu int                `json:"cpu"`
	Ram int                `json:"ram"`
	ImageID int            `json:"imageId"`
	BootDisk int           `json:"bootDisk"`
	NetType string         `json:"netType"`
	NetId int              `json:"netId"`
	IPAddr string          `json:"ipAddr"`
	UserData string        `json:"userdata"`
	Description string     `json:"desc"`
	Start bool             `json:"start"`
}

// structures related to cloudapi/compute/delete API
const ComputeDeleteAPI = "/restmachine/cloudapi/compute/delete"

type ComputeDeleteParam struct {
	ComputeID int          `json:"computeId"`
	Permanently bool       `json:"permanently"`
}

// 
// structures related to /cloudapi/compute/list API
//

type InterfaceRecord struct {
	ConnID int             `json:"connId"`
	ConnType string        `json:"connType"`
	DefaultGW string       `json:"defGw"`
	Guid string            `json:"guid"`
	IPAddress string       `json:"ipAddress"` // without trailing network mask, i.e. "192.168.1.3"
	MAC string             `json:"mac"`
	Name string            `json:"name"`
	NetID int              `json:"netId"` 
	NetMaks int            `json:"netMask"`
	NetType string         `json:"netType"`
	PciSlot int            `json:"pciSlot"`
	Target string          `json:"target"`
	Type string            `json:"type"`
	VNFs []int             `json:"vnfs"`
}

type SnapSetRecord struct {
	Disks []int            `json:"disks"`
	Guid string            `json:"guid"`
	Label string           `json:"label"`
	TimeStamp uint64       `json:"timestamp"`
}

type ComputeRecord struct {
	AccountID int          `json:"accountId"`
	AccountName string     `json:"accountName"`
	ACLs []UserAclRecord   `json:"acl"`
	Arch string            `json:"arch"`
	BootDiskSize int       `json:"bootdiskSize"`
	CloneReference int     `json:"cloneReference"`
	Clones []int           `json:"clones"`
	Cpus int               `json:"cpus"`
	CreatedBy string       `json:"createdBy"`
	CreatedTime uint64     `json:"createdTime"`
	DeletedBy string       `json:"deletedBy"`
	DeletedTime uint64     `json:"deletedTime"`
	Desc string            `json:"desc"`
	Disks []int            `json:"disks"`
	GridID int             `json:"gid"`
	ID uint                `json:"id"`
	ImageID int            `json:"imageId"`
	Interfaces []InterfaceRecord `json:"interfaces`
	LockStatus string      `json:"lockStatus"`
	ManagerID int          `json:"managerId"`
	Name string            `json:"name"`
	Ram int                `json:"ram"`
	RgID int               `json:"rgId"`
	RgName string          `json:"rgName"`
	SnapSets []SnapSetRecord `json:"snapSets"`
	Status string          `json:"status"`
	Tags []string          `json:"tags"`
	TechStatus string      `json:"techStatus"`
	TotalDiskSize int      `json:"totalDiskSize"`
	UpdatedBy string       `json:"updatedBy"`
	UpdateTime uint64      `json:"updateTime"`
	UserManaged bool       `json:"userManaged"`
	Vgpus []int            `json:"vgpus"`
	VinsConnected int      `json:"vinsConnected"`
	VirtualImageID int     `json:"virtualImageId"`
}

const ComputeListAPI = "/restmachine/cloudapi/compute/list"
type ComputeListParam struct {
	IncludeDeleted bool    `json:"includedeleted"`
}
type ComputeListResp []ComputeRecord

//
// structures related to /cloudapi/compute/get
//
type SnapshotRecord struct {
	Guid string            `json:"guid"`
	Label string           `json:"label"`
	SnapSetGuid string     `json:"snapSetGuid"`
	SnapSetTime uint64     `json:"snapSetTime"`
	TimeStamp uint64       `json:"timestamp"`
}

type DiskRecord struct {
	// ACLs `json:"ACL"` - it is a dictionary, special parsing required
	// was - Acl map[string]string  `json:"acl"`
	AccountID int          `json:"accountId"`
	BootPartition int      `json:"bootPartition"`
	CreatedTime uint64     `json:"creationTime"`
	DeletedTime uint64     `json:"deletionTime"`
	Description string     `json:"descr"`
	DestructionTime uint64 `json:"destructionTime"`
	DiskPath string        `json:"diskPath"`
	GridID int             `json:"gid"`
	ID uint                `json:"id"`
	ImageID int            `json:"imageId"`
	Images []int           `json:"images"`
	// IOTune 'json:"iotune" - it is a dictionary
	Name string            `json:"name"`
	ParentId int           `json:"parentId"`
	PciSlot int            `json:"pciSlot"`
	// ResID string           `json:"resId"`
	// ResName string         `json:"resName"`
	// Params string          `json:"params"`
	Pool string            `json:"pool"`
	PurgeTime uint64       `json:"purgeTime"`
	// Role string            `json:"role"`
	SepType string         `json:"sepType"`
	SepID int              `json:"sepid"`
	SizeMax int            `json:"sizeMax"`
	SizeUsed int           `json:"sizeUsed"`
	Snapshots []SnapshotRecord `json:"snapshots"`
	Status string          `json:"status"`
	TechStatus string      `json:"techStatus"`
	Type string            `json:"type"`
	ComputeID int          `json:"vmId"`
}

type OsUserRecord struct {
	Guid string            `json:"guid"`
	Login string           `json:"login"`
	Password string        `json:"password"`
	PubKey string          `json:"pubkey"`
}

const ComputeGetAPI = "/restmachine/cloudapi/compute/get"
type ComputeGetParam struct {
	ComputeID int          `json:"computeId"`
}
type ComputeGetResp struct {
	// ACLs `json:"ACL"` - it is a dictionary, special parsing required
	AccountID int          `json:"accountId"`
	AccountName string     `json:"accountName"`
	Arch string            `json:"arch"`
	BootDiskSize int       `json:"bootdiskSize"`
	CloneReference int     `json:"cloneReference"`
	Clones []int           `json:"clones"`
	Cpus int               `json:"cpus"`
	Desc string            `json:"desc"`
	Disks []DiskRecord     `json:"disks"`
	GridID int             `json:"gid"`
	ID uint                `json:"id"`
	ImageID int            `json:"imageId"`
	ImageName string       `json:"imageName"`
	Interfaces []InterfaceRecord `json:"interfaces`
	LockStatus string      `json:"lockStatus"`
	ManagerID int          `json:"managerId"`
	ManagerType string     `json:"manageType"`
	Name string            `json:"name"`
	NatableVinsID int      `json:"natableVinsId"`
    NatableVinsIP string   `json:"natableVinsIp"`
    NatableVinsName string `json:"natableVinsName"`
    NatableVinsNet string  `json:"natableVinsNetwork"`
	NatableVinsNetName string `json:"natableVinsNetworkName"`
	OsUsers []OsUserRecord `json:"osUsers"`
	Ram int                `json:"ram"`
	RgID int               `json:"rgId"`
	RgName string          `json:"rgName"`
	SnapSets []SnapSetRecord `json:"snapSets"`
	Status string          `json:"status"`
	Tags []string          `json:"tags"`
	TechStatus string      `json:"techStatus"`
	TotalDiskSize int      `json:"totalDiskSize"`
	UpdatedBy string       `json:"updatedBy"`
	UpdateTime uint64      `json:"updateTime"`
	UserManaged bool       `json:"userManaged"`
	Vgpus []int            `json:"vgpus"`
	VinsConnected int      `json:"vinsConnected"`
	VirtualImageID int     `json:"virtualImageId"`
}

//
// structures related to /restmachine/cloudapi/images/list API
//
type ImageRecord struct {
	AccountID uint      `json:"accountId"`
	Arch string         `json:"architecture`
	BootType string     `json:"bootType"`
	IsBootable boo      `json:"bootable"`
	IsCdrom bool        `json:"cdrom"`
	Desc string         `json:"description"`
	IsHotResize bool    `json:"hotResize"`
	ID uint             `json:"id"`
	Name string         `json:"name"`
	Pool string         `json:"pool"`
	SepID int           `json:"sepid"`
	Size int            `json:"size"`
	Status string       `json:"status"`
	Type string         `json:"type"`
	Username string     `json:"username"`
	IsVirtual bool      `json:"virtual"`
}

const ImagesListAPI = "/restmachine/cloudapi/images/list"
type ImagesListParam struct {
	AccountID int        `json:"accountId"`
}
type ImagesListResp []ImageRecord

//
// structures related to /cloudapi/extnet/list API
//
type ExtNetRecord struct {
	Name string          `json:"name"`
	ID uint              `json:"id"`
	IPCIDR string        `json:"ipcidr"`
} 

const ExtNetListAPI = "/restmachine/cloudapi/extnet/list"
type ExtNetListParam struct {
	AccountID int        `json:"accountId"`
}
type ExtNetListResp []ExtNetRecord


//
// structures related to /cloudapi/accounts/list API
//
type AccountRecord struct {
	ACLs []UserAclRecord    `json:"acl"`
	CreatedTime uint64      `json:"creationTime"`
	DeletedTime uint64      `json:"deletionTime"`
	ID int                  `json:"id"`
	Name string             `json:"name"`
	Status string           `json:"status"`
	UpdatedTime uint64      `json:"updateTime"`
}

const AccountsListAPI = "/restmachine/cloudapi/accounts/list"
type AccountsListResp []AccountRecord

//
// structures related to /cloudapi/portforwarding/list API
//
type PfwRecord struct {
	ID int                 `json:"id"`
	LocalIP string         `json:"localIp`
	LocalPort int          `json:"localPort"`
	Protocol string        `json:"protocol"`
	PublicPortEnd int      `json:"publicPortEnd"`
	PublicPortStart int    `json:"publicPortStart"`
	ComputeID int          `json:"vmId"`
} 

const ComputePfwListAPI = "/restmachine/cloudapi/compute/pfwList"
type ComputePfwListResp []PfwRecord

type ComputePfwAddParam struct {
	ComputeID int          `json:"computeId"`
	PublicPortStart int    `json:"publicPortStart"`
	PublicPortEnd int      `json:"publicPortEnd"`
	LocalBasePort int      `json:"localBasePort"`
	Protocol string        `json:"proto"`
}
const ComputePfwAddAPI = "/restmachine/cloudapi/compute/pfwAdd"

type ComputePfwDelParam struct {
	ComputeID int          `json:"computeId"`
	RuleID int             `json:"ruleId"`
	PublicPortStart int    `json:"publicPortStart"`
	PublicPortEnd int      `json:"publicPortEnd"`
	LocalBasePort int      `json:"localBasePort"`
	Protocol string        `json:"proto"`
}
const ComputePfwDelAPI = "/restmachine/cloudapi/compute/pfwDel"

//
// structures related to /cloudapi/compute/net Attach/Detach API
//
type ComputeNetAttachParam struct {
	ComputeID int          `json:"computeId"`
	NetType string         `json:"netType"`
	NetID int              `json:"netId"`
	IPAddr string          `json:"apAddr"`
}
const ComputeNetAttachAPI = "/restmachine/cloudapi/compute/netAttach"

type ComputeNetDetachParam struct {
	ComputeID int          `json:"computeId"`
	IPAddr string          `json:"apAddr"`
	MAC string             `json:"mac"`
}
const ComputeNetDetachAPI = "/restmachine/cloudapi/compute/netDetach"


//
// structures related to /cloudapi/compute/disk Attach/Detach API
//
type ComputeDiskManipulationParam struct {
	ComputeID int          `json:"computeId"`
	DiskID int             `json:"diskId"`
}
const ComputeDiskAttachAPI = "/restmachine/cloudapi/compute/diskAttach"

const ComputeDiskDetachAPI = "/restmachine/cloudapi/compute/diskDetach"

//
// structures related to /cloudapi/disks/create
// 
type DiskCreateParam struct {
	AccountID int          `json:"accountId`
	GridID int             `json:"gid"`
	Name string            `json:"string"`
	Description string     `json:"description"`
	Size int               `json:"size"`
	Type string            `json:"type"`
	SepID int              `json:"sep_id"`
	Pool string            `json:"pool"`
}
const DiskCreateAPI = "/restmachine/cloudapi/disks/create"

//
// structures related to /cloudapi/disks/get
// 
type DisksGetParam struct {
	DiskID int           `json:"diskId`
}
const DisksCreateAPI = "/restmachine/cloudapi/disks/create"

const DisksGetAPI = "/restmachine/cloudapi/disks/get" // Returns single DiskRecord on success
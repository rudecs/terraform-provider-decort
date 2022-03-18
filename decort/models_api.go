/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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
	"bytes"
	"strconv"
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
	IsExplicit bool   `json:"explicit"`
	Rights     string `json:"right"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	UgroupID   string `json:"userGroupId"`
	// CanBeDeleted bool      `json:"canBeDeleted"`
}

type AccountAclRecord struct {
	IsExplicit bool   `json:"explicit"`
	Guid       string `json:"guid"`
	Rights     string `json:"right"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	UgroupID   string `json:"userGroupId"`
}

type ResgroupRecord struct {
	ACLs           []UserAclRecord  `json:"acl"`
	Owner          AccountAclRecord `json:"accountAcl"`
	AccountID      int              `json:"accountId"`
	AccountName    string           `json:"accountName"`
	CreatedBy      string           `json:"createdBy"`
	CreatedTime    uint64           `json:"createdTime"`
	DefaultNetID   int              `json:"def_net_id"`
	DefaultNetType string           `json:"def_net_type"`
	Decsription    string           `json:"desc"`
	GridID         int              `json:"gid"`
	ID             uint             `json:"id"`
	LockStatus     string           `json:"lockStatus"`
	Name           string           `json:"name"`
	Status         string           `json:"status"`
	UpdatedBy      string           `json:"updatedBy"`
	UpdatedTime    uint64           `json:"updatedTime"`
	Vins           []int            `json:"vins"`
	Computes       []int            `json:"vms"`
}

const ResgroupListAPI = "/restmachine/cloudapi/rg/list"

type ResgroupListResp []ResgroupRecord

//
// structures related to /cloudapi/rg/create API call
//
const ResgroupCreateAPI = "/restmachine/cloudapi/rg/create"

//
// structures related to /cloudapi/rg/update API call
//
const ResgroupUpdateAPI = "/restmachine/cloudapi/rg/update"

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

//
// structures related to /cloudapi/rg/get API call
//
type QuotaRecord struct { // this is how quota is reported by /api/.../rg/get
	Cpu        int     `json:"CU_C"`      // CPU count in pcs
	Ram        float64 `json:"CU_M"`      // RAM volume in MB, it is STILL reported as FLOAT
	Disk       int     `json:"CU_D"`      // Disk capacity in GB
	ExtIPs     int     `json:"CU_I"`      // Ext IPs count
	ExtTraffic int     `json:"CU_NP"`     // Ext network traffic
	GpuUnits   int     `json:"gpu_units"` // GPU count
}

type ResourceRecord struct { // this is how actual usage is reported by /api/.../rg/get
	Cpu        int `json:"cpu"`
	Disk       int `json:"disksize"`
	ExtIPs     int `json:"extips"`
	ExtTraffic int `json:"exttraffic"`
	Gpu        int `json:"gpu"`
	Ram        int `json:"ram"`
}

type UsageRecord struct {
	Current  ResourceRecord `json:"Current"`
	Reserved ResourceRecord `json:"Reserved"`
}

const ResgroupGetAPI = "/restmachine/cloudapi/rg/get"

type ResgroupGetResp struct {
	ACLs           []UserAclRecord `json:"ACLs"`
	Usage          UsageRecord     `json:"Resources"`
	AccountID      int             `json:"accountId"`
	AccountName    string          `json:"accountName"`
	GridID         int             `json:"gid"`
	CreatedBy      string          `json:"createdBy"`
	CreatedTime    uint64          `json:"createdTime"`
	DefaultNetID   int             `json:"def_net_id"`
	DefaultNetType string          `json:"def_net_type"`
	DeletedBy      string          `json:"deletedBy"`
	DeletedTime    uint64          `json:"deletedTime"`
	Desc           string          `json:"desc"`
	ID             uint            `json:"id"`
	LockStatus     string          `json:"lockStatus"`
	Name           string          `json:"name"`
	Quota          QuotaRecord     `json:"resourceLimits"`
	Status         string          `json:"status"`
	UpdatedBy      string          `json:"updatedBy"`
	UpdatedTime    uint64          `json:"updatedTime"`
	Vins           []int           `json:"vins"`
	Computes       []int           `json:"vms"`

	Ignored map[string]interface{} `json:"-"`
}

//
// structures related to /cloudapi/rg/delete API
//
const ResgroupDeleteAPI = "/restmachine/cloudapi/rg/delete"

//
// structures related to /cloudapi/rg/listComputes API
//
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

const RgListComputesAPI = "/restmachine/cloudapi/rg/listComputes"

type RgListComputesResp []ComputeBriefRecord

//
// structures related to /cloudapi/kvmXXX/create APIs
//
const KvmX86CreateAPI = "/restmachine/cloudapi/kvmx86/create"
const KvmPPCCreateAPI = "/restmachine/cloudapi/kvmppc/create"

type KvmVmCreateParam struct { // this is unified structure for both x86 and PPC based KVM VMs creation
	RgID     uint   `json:"rgId"`
	Name     string `json:"name"`
	Cpu      int    `json:"cpu"`
	Ram      int    `json:"ram"`
	ImageID  int    `json:"imageId"`
	BootDisk int    `json:"bootDisk"`
	NetType  string `json:"netType"`
	NetId    int    `json:"netId"`
	IPAddr   string `json:"ipAddr"`
	UserData string `json:"userdata"`
	Desc     string `json:"desc"`
	Start    bool   `json:"start"`
}

// structures related to cloudapi/compute/start API
const ComputeStartAPI = "/restmachine/cloudapi/compute/start"
const ComputeStopAPI = "/restmachine/cloudapi/compute/stop"

// structures related to cloudapi/compute/delete API
const ComputeDeleteAPI = "/restmachine/cloudapi/compute/delete"

//
// structures related to /cloudapi/compute/list API
//

type InterfaceQosRecord struct {
	ERate   int    `json:"eRate"`
	Guid    string `json:"guid"`
	InBurst int    `json:"inBurst"`
	InRate  int    `json:"inRate"`
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

type SnapSetRecord struct {
	Disks     []int  `json:"disks"`
	Guid      string `json:"guid"`
	Label     string `json:"label"`
	TimeStamp uint64 `json:"timestamp"`
}

type ComputeRecord struct {
	AccountID      int               `json:"accountId"`
	AccountName    string            `json:"accountName"`
	ACLs           []UserAclRecord   `json:"acl"`
	Arch           string            `json:"arch"`
	BootDiskSize   int               `json:"bootdiskSize"`
	CloneReference int               `json:"cloneReference"`
	Clones         []int             `json:"clones"`
	Cpus           int               `json:"cpus"`
	CreatedBy      string            `json:"createdBy"`
	CreatedTime    uint64            `json:"createdTime"`
	DeletedBy      string            `json:"deletedBy"`
	DeletedTime    uint64            `json:"deletedTime"`
	Desc           string            `json:"desc"`
	Disks          []int             `json:"disks"`
	GridID         int               `json:"gid"`
	ID             uint              `json:"id"`
	ImageID        int               `json:"imageId"`
	Interfaces     []InterfaceRecord `json:"interfaces"`
	LockStatus     string            `json:"lockStatus"`
	ManagerID      int               `json:"managerId"`
	Name           string            `json:"name"`
	Ram            int               `json:"ram"`
	RgID           int               `json:"rgId"`
	RgName         string            `json:"rgName"`
	SnapSets       []SnapSetRecord   `json:"snapSets"`
	Status         string            `json:"status"`
	// Tags           []string          `json:"tags"` // Tags were reworked since DECORT 3.7.1
	TechStatus     string `json:"techStatus"`
	TotalDiskSize  int    `json:"totalDiskSize"`
	UpdatedBy      string `json:"updatedBy"`
	UpdateTime     uint64 `json:"updateTime"`
	UserManaged    bool   `json:"userManaged"`
	Vgpus          []int  `json:"vgpus"`
	VinsConnected  int    `json:"vinsConnected"`
	VirtualImageID int    `json:"virtualImageId"`
}

const ComputeListAPI = "/restmachine/cloudapi/compute/list"

type ComputeListResp []ComputeRecord

const ComputeResizeAPI = "/restmachine/cloudapi/compute/resize"

//
// structures related to /cloudapi/compute/get
//
type SnapshotRecord struct {
	Guid        string `json:"guid"`
	Label       string `json:"label"`
	SnapSetGuid string `json:"snapSetGuid"`
	SnapSetTime uint64 `json:"snapSetTime"`
	TimeStamp   uint64 `json:"timestamp"`
}

type DiskRecord struct {
	// ACLs `json:"ACL"` - it is a dictionary, special parsing required
	// was - Acl map[string]string  `json:"acl"`
	AccountID       int    `json:"accountId"`
	AccountName     string `json:"accountName"` // NOTE: absent from compute/get output
	BootPartition   int    `json:"bootPartition"`
	CreatedTime     uint64 `json:"creationTime"`
	DeletedTime     uint64 `json:"deletionTime"`
	Desc            string `json:"desc"`
	DestructionTime uint64 `json:"destructionTime"`
	DiskPath        string `json:"diskPath"`
	GridID          int    `json:"gid"`
	ID              uint   `json:"id"`
	ImageID         int    `json:"imageId"`
	Images          []int  `json:"images"`
	// IOTune 'json:"iotune" - it is a dictionary
	Name string `json:"name"`
	// Order                   `json:"order"`
	ParentId int `json:"parentId"`
	PciSlot  int `json:"pciSlot"`
	// ResID string           `json:"resId"`
	// ResName string         `json:"resName"`
	// Params string          `json:"params"`
	Pool      string `json:"pool"`
	PurgeTime uint64 `json:"purgeTime"`
	// Role string            `json:"role"`
	SepType    string           `json:"sepType"`
	SepID      int              `json:"sepId"` // NOTE: absent from compute/get output
	SizeMax    int              `json:"sizeMax"`
	SizeUsed   int              `json:"sizeUsed"` // sum over all snapshots of this disk to report total consumed space
	Snapshots  []SnapshotRecord `json:"snapshots"`
	Status     string           `json:"status"`
	TechStatus string           `json:"techStatus"`
	Type       string           `json:"type"`
	ComputeID  int              `json:"vmid"`
}

type OsUserRecord struct {
	Guid     string `json:"guid"`
	Login    string `json:"login"`
	Password string `json:"password"`
	PubKey   string `json:"pubkey"`
}

const ComputeGetAPI = "/restmachine/cloudapi/compute/get"

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

//
// structures related to /restmachine/cloudapi/image/list API
//
type ImageRecord struct {
	AccountID   uint   `json:"accountId"`
	Arch        string `json:"architecture"`
	BootType    string `json:"bootType"`
	IsBootable  bool   `json:"bootable"`
	IsCdrom     bool   `json:"cdrom"`
	Desc        string `json:"desc"`
	IsHotResize bool   `json:"hotResize"`
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Pool        string `json:"pool"`
	SepID       int    `json:"sepId"`
	Size        int    `json:"size"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Username    string `json:"username"`
	IsVirtual   bool   `json:"virtual"`
}

const ImagesListAPI = "/restmachine/cloudapi/image/list"

type ImagesListResp []ImageRecord

//
// structures related to /cloudapi/extnet/list API
//
type ExtNetRecord struct {
	Name   string `json:"name"`
	ID     uint   `json:"id"`
	IPCIDR string `json:"ipcidr"`
}

const ExtNetListAPI = "/restmachine/cloudapi/extnet/list"

type ExtNetListResp []ExtNetRecord

//
// structures related to /cloudapi/account/list API
//
type AccountRecord struct {
	// ACLs        []UserAclRecord `json:"acl"`
	// CreatedTime uint64          `json:"creationTime"`
	// DeletedTime uint64          `json:"deletionTime"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	// UpdatedTime uint64          `json:"updateTime"`
}

const AccountsGetAPI = "/restmachine/cloudapi/account/get" // returns AccountRecord superset

const AccountsListAPI = "/restmachine/cloudapi/account/list" // returns list of abdridged info about accounts
type AccountsListResp []AccountRecord

//
// structures related to /cloudapi/portforwarding/list API
//
type PfwRecord struct {
	ID              int    `json:"id"`
	LocalIP         string `json:"localIp"`
	LocalPort       int    `json:"localPort"`
	Protocol        string `json:"protocol"`
	PublicPortEnd   int    `json:"publicPortEnd"`
	PublicPortStart int    `json:"publicPortStart"`
	ComputeID       int    `json:"vmId"`
}

const ComputePfwListAPI = "/restmachine/cloudapi/compute/pfwList"

type ComputePfwListResp []PfwRecord

const ComputePfwAddAPI = "/restmachine/cloudapi/compute/pfwAdd"

const ComputePfwDelAPI = "/restmachine/cloudapi/compute/pfwDel"

//
// structures related to /cloudapi/compute/net Attach/Detach API
//
type ComputeNetMgmtRecord struct { // used to "cache" network specs when preparing to manage compute networks
	ID        int
	Type      string
	IPAddress string
	MAC       string
}

const ComputeNetAttachAPI = "/restmachine/cloudapi/compute/netAttach"

const ComputeNetDetachAPI = "/restmachine/cloudapi/compute/netDetach"

//
// structures related to /cloudapi/compute/disk Attach/Detach API
//
const ComputeDiskAttachAPI = "/restmachine/cloudapi/compute/diskAttach"

const ComputeDiskDetachAPI = "/restmachine/cloudapi/compute/diskDetach"

//
// structures related to /cloudapi/disks/create
//
const DisksCreateAPI = "/restmachine/cloudapi/disks/create"

//
// structures related to /cloudapi/disks/get
//
const DisksGetAPI = "/restmachine/cloudapi/disks/get" // Returns single DiskRecord on success

const DisksListAPI = "/restmachine/cloudapi/disks/list" // Returns list of DiskRecord on success
type DisksListResp []DiskRecord

//
// structures related to /cloudapi/disks/resize
//
const DisksResizeAPI = "/restmachine/cloudapi/disks/resize2"

//
// structures related to /cloudapi/disks/resize
//
const DisksRenameAPI = "/restmachine/cloudapi/disks/rename"

//
// structures related to /cloudapi/disks/delete
//
const DisksDeleteAPI = "/restmachine/cloudapi/disks/delete"

//
// ViNS structures
//

// this is the structure of the element in the list returned by vins/search API
type VinsSearchRecord struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IPCidr      string `json:"network"`
	VxLanID     int    `json:"vxlanId"`
	ExternalIP  string `json:"externalIP"`
	AccountID   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	RgID        int    `json:"rgId"`
	RgName      string `json:"rgName"`
}

const VinsSearchAPI = "/restmachine/cloudapi/vins/search"

type VinsSearchResp []VinsSearchRecord

type VnfRecord struct {
	ID        int                    `json:"id"`
	AccountID int                    `json:"accountId"`
	Type      string                 `json:"type"`   // "DHCP", "NAT", "GW" etc
	Config    map[string]interface{} `json:"config"` // NOTE: VNF specs vary by VNF type
}

type VnfGwConfigRecord struct { // describes GW VNF config structure inside ViNS, as returned by API vins/get
	ExtNetID   int    `json:"ext_net_id"`
	ExtNetIP   string `json:"ext_net_ip"`
	ExtNetMask int    `json:"ext_net_mask"`
	DefaultGW  string `json:"default_gw"`
}
type VinsRecord struct { // represents part of the response from API vins/get
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	IPCidr      string               `json:"network"`
	VxLanID     int                  `json:"vxlanId"`
	ExternalIP  string               `json:"externalIP"`
	AccountID   int                  `json:"accountId"`
	AccountName string               `json:"accountName"`
	RgID        int                  `json:"rgid"`
	RgName      string               `json:"rgName"`
	VNFs        map[string]VnfRecord `json:"vnfs"`
	Desc        string               `json:"desc"`
}

const VinsGetAPI = "/restmachine/cloudapi/vins/get"

const VinsCreateInAccountAPI = "/restmachine/cloudapi/vins/createInAccount"
const VinsCreateInRgAPI = "/restmachine/cloudapi/vins/createInRG"

const VinsExtNetConnectAPI = "/restmachine/cloudapi/vins/extNetConnect"
const VinsExtNetDisconnectAPI = "/restmachine/cloudapi/vins/extNetDisconnect"

const VinsDeleteAPI = "/restmachine/cloudapi/vins/delete"

//
// K8s structures
//

//K8sNodeRecord represents a worker/master group
type K8sNodeRecord struct {
	ID   int `json:"id"`
	Disk int `json:"disk"`
	Cpu  int `json:"cpu"`
	Num  int `json:"num"`
	Ram  int `json:"ram"`
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
	Name   string `json:"name"`
	RgID   int    `json:"rgId"`
	RgName string `json:"rgName"`
}

const K8sCreateAPI = "/restmachine/cloudapi/k8s/create"
const K8sGetAPI = "/restmachine/cloudapi/k8s/get"
const K8sUpdateAPI = "/restmachine/cloudapi/k8s/update"
const K8sDeleteAPI = "/restmachine/cloudapi/k8s/delete"

const K8sWgCreateAPI = "/restmachine/cloudapi/k8s/workersGroupAdd"
const K8sWgDeleteAPI = "/restmachine/cloudapi/k8s/workersGroupDelete"

//Blasphemous workaround for parsing Result value
type TaskResult int

func (r *TaskResult) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 {
		*r = 0
		return nil
	}

	n, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*r = TaskResult(n)
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

const AsyncTaskGetAPI = "/restmachine/cloudapi/tasks/get"

//
// Grid ID structures
//
type LocationRecord struct {
	GridID       int    `json:"gid"`
	Id           int    `json:"id"`
	LocationCode string `json:"locationCode"`
	Name         string `json:"name"`
	Flag         string `json:"flag"`
}

const LocationsListAPI = "/restmachine/cloudapi/locations/list" // Returns list of GridRecord on success
type LocationsListResp []LocationRecord

//
// Auxiliary structures
//
type SshKeyConfig struct {
	User      string
	SshKey    string
	UserShell string
}

////////////////////
// 	 IMAGE API	  //
////////////////////
const imageCreateAPI = "/restmachine/cloudbroker/image/createImage"
const imageSyncCreateAPI = "/restmachine/cloudbroker/image/syncCreateImage"
const imageCreateVirtualAPI = "/restmachine/cloudbroker/image/createVirtual"
const imageCreateCDROMAPI = "/restmachine/cloudbroker/image/createCDROMImage"
const imageListStacksApi = "/restmachine/cloudbroker/image/listStacks"
const imageGetAPI = "/restmachine/cloudbroker/image/get"
const imageListGetAPI = "/restmachine/cloudbroker/image/list"
const imageEditAPI = "/restmachine/cloudbroker/image/edit"
const imageDeleteAPI = "/restmachine/cloudbroker/image/delete"
const imageDeleteCDROMAPI = "/restmachine/cloudbroker/image/deleteCDROMImage"
const imageEnableAPI = "/restmachine/cloudbroker/image/enable"
const imageDisableAPI = "/restmachine/cloudbroker/image/disable"
const imageEditNameAPI = "/restmachine/cloudbroker/image/rename"
const imageLinkAPI = "/restmachine/cloudbroker/image/link"
const imageShareAPI = "/restmachine/cloudbroker/image/share"
const imageComputeciSetAPI = "/restmachine/cloudbroker/image/computeciSet"
const imageComputeciUnsetAPI = "/restmachine/cloudbroker/image/computeciUnset"
const imageUpdateNodesAPI = "/restmachine/cloudbroker/image/updateNodes"
const imageDeleteImagesAPI = "/restmachine/cloudbroker/image/deleteImages"

type History struct {
	Guid      string `json:"guid"`
	Id        int    `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type Image struct {
	ImageId       int           `json:"id"`
	Name          string        `json:"name"`
	Url           string        `json:"url"`
	Gid           int           `json:"gid"`
	Guid          int           `json:"guid"`
	Boottype      string        `json:"bootType"`
	Imagetype     string        `json:"type"`
	Drivers       []string      `json:"drivers"`
	Hotresize     bool          `json:"hotResize"`
	Bootable      bool          `json:"bootable"`
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	AccountId     int           `json:"accountId"`
	UsernameDL    string        `json:"usernameDL"`
	PasswordDL    string        `json:"passwordDL"`
	SepId         int           `json:"sepId"`
	PoolName      string        `json:"pool"`
	Architecture  string        `json:"architecture"`
	UNCPath       string        `json:"UNCPath"`
	LinkTo        int           `json:"linkTo"`
	Status        string        `json:"status"`
	TechStatus    string        `json:"techStatus"`
	Size          int           `json:"size"`
	Version       string        `json:"version"`
	Enabled       bool          `json:"enabled"`
	ComputeciId   int           `json:"computeciId"`
	Milestones    int           `json:"milestones"`
	ProviderName  string        `json:"provider_name"`
	PurgeAttempts int           `json:"purgeAttempts"`
	ReferenceId   string        `json:"referenceId"`
	ResId         string        `json:"resId"`
	ResName       string        `json:"resName"`
	Rescuecd      bool          `json:"rescuecd"`
	Meta          []interface{} `json:"_meta"`
	History       []History     `json:"history"`
	LastModified  int64         `json:"lastModified"`
	Desc          string        `json:"desc"`
	SharedWith    []int         `json:"sharedWith"`
}

type ImageList []Image

type ImageStack struct {
	ApiURL      string   `json:"apiUrl"`
	ApiKey      string   `json:"apikey"`
	AppId       string   `json:"appId"`
	Desc        string   `json:"desc"`
	Drivers     []string `json:"drivers"`
	Error       int      `json:"error"`
	Guid        int      `json:"guid"`
	Id          int      `json:"id"`
	Images      []int    `json:"images"`
	Login       string   `json:"login"`
	Name        string   `json:"name"`
	Passwd      string   `json:"passwd"`
	ReferenceId string   `json:"referenceId"`
	Status      string   `json:"status"`
	Type        string   `json:"type"`
}

type ImageListStacks []ImageStack

/////////////////
//  GRID API   //
/////////////////
const GridListGetAPI = "/restmachine/cloudbroker/grid/list"
const GridGetAPI = "/restmachine/cloudbroker/grid/get"

type Grid struct {
	Flag         string `json:"flag"`
	Gid          int    `json:"gid"`
	Guid         int    `json:"guid"`
	Id           int    `json:"id"`
	LocationCode string `json:"locationCode"`
	Name         string `json:"name"`
}

type GridList []Grid

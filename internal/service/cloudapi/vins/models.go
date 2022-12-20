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

package vins

type VINSRecord struct {
	AccountID   uint64 `json:"accountId"`
	AccountName string `json:"accountName"`
	CreatedBy   string `json:"createdBy"`
	CreatedTime uint64 `json:"createdTime"`
	DeletedBy   string `json:"deletedBy"`
	DeletedTime uint64 `json:"deletedTime"`
	ExternalIP  string `json:"externalIP"`
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Network     string `json:"network"`
	RGID        uint64 `json:"rgId"`
	RGName      string `json:"rgName"`
	Status      string `json:"status"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedTime uint64 `json:"updatedTime"`
	VXLANID     uint64 `json:"vxlanId"`
}

type VINSList []VINSRecord

type VINSAudits struct {
	Call         string  `json:"call"`
	ResponseTime float64 `json:"responsetime"`
	StatusCode   uint64  `json:"statuscode"`
	Timestamp    float64 `json:"timestamp"`
	User         string  `json:"user"`
}

type VINSAuditsList []VINSAudits

type VINSExtNet struct {
	DefaultGW  string `json:"default_gw"`
	ExtNetID   uint64 `json:"ext_net_id"`
	IP         string `json:"ip"`
	PrefixLen  uint64 `json:"prefixlen"`
	Status     string `json:"status"`
	TechStatus string `json:"techStatus"`
}

type ExtNetList []VINSExtNet

type IP struct {
	ClientType string `json:"clientType"`
	DomainName string `json:"domainname"`
	HostName   string `json:"hostname"`
	IP         string `json:"ip"`
	MAC        string `json:"mac"`
	Type       string `json:"type"`
	VMID       uint64 `json:"vmId"`
}

type IPList []IP

type VNFDev struct {
	CKey            string           `json:"_ckey"`
	AccountID       uint64           `json:"accountId"`
	Capabilities    []string         `json:"capabilities"`
	Config          VNFConfig        `json:"config"`
	ConfigSaved     bool             `json:"configSaved"`
	CustomPreConfig bool             `json:"customPrecfg"`
	Description     string           `json:"desc"`
	GID             uint64           `json:"gid"`
	GUID            uint64           `json:"guid"`
	ID              uint64           `json:"id"`
	Interfaces      VNFInterfaceList `json:"interfaces"`
	LockStatus      string           `json:"lockStatus"`
	Milestones      uint64           `json:"milestones"`
	Name            string           `json:"name"`
	Status          string           `json:"status"`
	TechStatus      string           `json:"techStatus"`
	Type            string           `json:"type"`
	VINS            []uint64         `json:"vins"`
}

type VNFConfig struct {
	MGMT      VNFConfigMGMT      `json:"mgmt"`
	Resources VNFConfigResources `json:"resources"`
}

type VNFConfigMGMT struct {
	IPAddr   string `json:"ipaddr"`
	Password string `json:"password"`
	SSHKey   string `json:"sshkey"`
	User     string `json:"user"`
}

type VNFConfigResources struct {
	CPU     uint64 `json:"cpu"`
	RAM     uint64 `json:"ram"`
	StackID uint64 `json:"stackId"`
	UUID    string `json:"uuid"`
}

type VNFInterface struct {
	ConnID      uint64   `json:"connId"`
	ConnType    string   `json:"connType"`
	DefGW       string   `json:"defGw"`
	FlipGroupID uint64   `json:"flipgroupId"`
	GUID        string   `json:"guid"`
	IPAddress   string   `json:"ipAddress"`
	ListenSSH   bool     `json:"listenSsh"`
	MAC         string   `json:"mac"`
	Name        string   `json:"name"`
	NetID       uint64   `json:"netId"`
	NetMask     uint64   `json:"netMask"`
	NetType     string   `json:"netType"`
	PCISlot     uint64   `json:"pciSlot"`
	QOS         QOS      `json:"qos"`
	Target      string   `json:"target"`
	Type        string   `json:"type"`
	VNFS        []uint64 `json:"vnfs"`
}

type QOS struct {
	ERate   uint64 `json:"eRate"`
	GUID    string `json:"guid"`
	InBurst uint64 `json:"inBurst"`
	InRate  uint64 `json:"inRate"`
}

type VNFInterfaceList []VNFInterface

type VINSCompute struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type VINSComputeList []VINSCompute

type VNFS struct {
	DHCP DHCP `json:"DHCP"`
	GW   GW   `json:"GW"`
	NAT  NAT  `json:"NAT"`
}

type NAT struct {
	CKey        string    `json:"_ckey"`
	AccountID   uint64    `json:"accountId"`
	CreatedTime uint64    `json:"createdTime"`
	Config      NATConfig `json:"config"`
	Devices     Devices   `json:"devices"`
	GID         uint64    `json:"gid"`
	GUID        uint64    `json:"guid"`
	ID          uint64    `json:"id"`
	LockStatus  string    `json:"lockStatus"`
	Milestones  uint64    `json:"milestones"`
	OwnerID     uint64    `json:"ownerId"`
	OwnerType   string    `json:"ownerType"`
	PureVirtual bool      `json:"pureVirtual"`
	Status      string    `json:"status"`
	TechStatus  string    `json:"techStatus"`
	Type        string    `json:"type"`
}

type NATConfig struct {
	NetMask uint64       `json:"netmask"`
	Network string       `json:"network"`
	Rules   ListNATRules `json:"rules"`
}

type ItemNATRule struct {
	ID              uint64 `json:"id"`
	LocalIP         string `json:"localIp"`
	LocalPort       uint64 `json:"localPort"`
	Protocol        string `json:"protocol"`
	PublicPortEnd   uint64 `json:"publicPortEnd"`
	PublicPortStart uint64 `json:"publicPortStart"`
	VMID            uint64 `json:"vmId"`
	VMName          string `json:"vmName"`
}

type ListNATRules []ItemNATRule

type GW struct {
	CKey        string   `json:"_ckey"`
	AccountID   uint64   `json:"accountId"`
	Config      GWConfig `json:"config"`
	CreatedTime uint64   `json:"createdTime"`
	Devices     Devices  `json:"devices"`
	GID         uint64   `json:"gid"`
	GUID        uint64   `json:"guid"`
	ID          uint64   `json:"id"`
	LockStatus  string   `json:"lockStatus"`
	Milestones  uint64   `json:"milestones"`
	OwnerID     uint64   `json:"ownerId"`
	OwnerType   string   `json:"ownerType"`
	PureVirtual bool     `json:"pureVirtual"`
	Status      string   `json:"status"`
	TechStatus  string   `json:"techStatus"`
	Type        string   `json:"type"`
}

type GWConfig struct {
	DefaultGW  string `json:"default_gw"`
	ExtNetID   uint64 `json:"ext_net_id"`
	ExtNetIP   string `json:"ext_net_ip"`
	ExtNetMask uint64 `json:"ext_netmask"`
	QOS        QOS    `json:"qos"`
}

type Devices struct {
	Primary DevicePrimary `json:"primary"`
}

type DevicePrimary struct {
	DevID   uint64 `json:"devId"`
	IFace01 string `json:"iface01"`
	IFace02 string `json:"iface02"`
}

type DHCP struct {
	CKey        string     `json:"_ckey"`
	AccountID   uint64     `json:"accountId"`
	Config      DHCPConfig `json:"config"`
	CreatedTime uint64     `json:"createdTime"`
	Devices     Devices    `json:"devices"`
	GID         uint64     `json:"gid"`
	GUID        uint64     `json:"guid"`
	ID          uint64     `json:"id"`
	LockStatus  string     `json:"lockStatus"`
	Milestones  uint64     `json:"milestones"`
	OwnerID     uint64     `json:"ownerId"`
	OwnerType   string     `json:"ownerType"`
	PureVirtual bool       `json:"pureVirtual"`
	Status      string     `json:"status"`
	TechStatus  string     `json:"techStatus"`
	Type        string     `json:"type"`
}

type DHCPConfig struct {
	DefaultGW    string          `json:"default_gw"`
	DNS          []string        `json:"dns"`
	IPEnd        string          `json:"ip_end"`
	IPStart      string          `json:"ip_start"`
	Lease        uint64          `json:"lease"`
	Netmask      uint64          `json:"netmask"`
	Network      string          `json:"network"`
	Reservations ReservationList `json:"reservations"`
}

type VINSDetailed struct {
	VNFDev            VNFDev          `json:"VNFDev"`
	CKey              string          `json:"_ckey"`
	AccountID         uint64          `json:"accountId"`
	AccountName       string          `json:"accountName"`
	Computes          VINSComputeList `json:"computes"`
	DefaultGW         string          `json:"defaultGW"`
	DefaultQOS        QOS             `json:"defaultQos"`
	Description       string          `json:"desc"`
	GID               uint64          `json:"gid"`
	GUID              uint64          `json:"guid"`
	ID                uint64          `json:"id"`
	LockStatus        string          `json:"lockStatus"`
	ManagerID         uint64          `json:"managerId"`
	ManagerType       string          `json:"managerType"`
	Milestones        uint64          `json:"milestones"`
	Name              string          `json:"name"`
	NetMask           uint64          `json:"netMask"`
	Network           string          `json:"network"`
	PreReservaionsNum uint64          `json:"preReservationsNum"`
	Redundant         bool            `json:"redundant"`
	RGID              uint64          `json:"rgId"`
	RGName            string          `json:"rgName"`
	SecVNFDevID       uint64          `json:"secVnfDevId"`
	Status            string          `json:"status"`
	UserManaged       bool            `json:"userManaged"`
	VNFS              VNFS            `json:"vnfs"`
	VXLanID           uint64          `json:"vxlanId"`
}

type Reservation struct {
	ClientType  string `json:"clientType"`
	Description string `json:"desc"`
	DomainName  string `json:"domainname"`
	HostName    string `json:"hostname"`
	IP          string `json:"ip"`
	MAC         string `json:"mac"`
	Type        string `json:"type"`
	VMID        int    `json:"vmId"`
}

type ReservationList []Reservation

type NATRule struct {
	ID              uint64 `json:"id"`
	LocalIP         string `json:"localIp"`
	LocalPort       uint64 `json:"localPort"`
	Protocol        string `json:"protocol"`
	PublicPortEnd   uint64 `json:"publicPortEnd"`
	PublicPortStart uint64 `json:"publicPortStart"`
	VMID            uint64 `json:"vmId"`
	VMName          string `json:"vmName"`
}

type NATRuleList []NATRule

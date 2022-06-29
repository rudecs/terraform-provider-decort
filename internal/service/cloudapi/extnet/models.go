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

package extnet

type Extnet struct {
	ID     int    `json:"id"`
	IPCidr string `json:"ipcidr"`
	Name   string `json:"name"`
}
type ExtnetExtend struct {
	Extnet
	IPAddr string `json:"ipaddr"`
}

type ExtnetList []Extnet
type ExtnetExtendList []ExtnetExtend

type ExtnetComputes struct {
	AccountId   int              `json:"accountId"`
	AccountName string           `json:"accountName"`
	Extnets     ExtnetExtendList `json:"extnets"`
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	RGID        int              `json:"rgId"`
	RGName      string           `json:"rgName"`
}

type ExtnetComputesList []ExtnetComputes

type ExtnetQos struct {
	ERate   int    `json:"eRate"`
	GUID    string `json:"guid"`
	InBurst int    `json:"inBurst"`
	InRate  int    `json:"inRate"`
}

type ExtnetReservation struct {
	ClientType string `json:"clientType"`
	Desc       string `json:"desc"`
	DomainName string `json:"domainname"`
	HostName   string `json:"hostname"`
	IP         string `json:"ip"`
	MAC        string `json:"mac"`
	Type       string `json:"type"`
	VMID       int    `json:"vmId"`
}

type ExtnetReservations []ExtnetReservation

type ExtnetVNFS struct {
	DHCP int `json:"dhcp"`
}

type ExtnetDetailed struct {
	CKey               string             `json:"_ckey"`
	Meta               []interface{}      `json:"_meta"`
	CheckIPs           []string           `json:"checkIPs"`
	CheckIps           []string           `json:"checkIps"`
	Default            bool               `json:"default"`
	DefaultQos         ExtnetQos          `json:"defaultQos"`
	Desc               string             `json:"desc"`
	Dns                []string           `json:"dns"`
	Excluded           []string           `json:"excluded"`
	FreeIps            int                `json:"free_ips"`
	Gateway            string             `json:"gateway"`
	GID                int                `json:"gid"`
	GUID               int                `json:"guid"`
	ID                 int                `json:"id"`
	IPCidr             string             `json:"ipcidr"`
	Milestones         int                `json:"milestones"`
	Name               string             `json:"name"`
	Network            string             `json:"network"`
	NetworkId          int                `json:"networkId"`
	PreReservationsNum int                `json:"preReservationsNum"`
	Prefix             int                `json:"prefix"`
	PriVnfDevId        int                `json:"priVnfDevId"`
	Reservations       ExtnetReservations `json:"reservations"`
	SharedWith         []int              `json:"sharedWith"`
	Status             string             `json:"status"`
	VlanID             int                `json:"vlanId"`
	VNFS               ExtnetVNFS         `json:"vnfs"`
}

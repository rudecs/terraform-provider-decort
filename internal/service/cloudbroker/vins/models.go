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

package vins

type Vins struct {
	AccountId   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	CreatedBy   string `json:"createdBy"`
	CreatedTime int    `json:"createdTime"`
	DeletedBy   string `json:"deletedBy"`
	DeletedTime int    `json:"deletedTime"`
	ExternalIP  string `json:"externalIP"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Network     string `json:"network"`
	RGID        int    `json:"rgId"`
	RGName      string `json:"rgName"`
	Status      string `json:"status"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedTime int    `json:"updatedTime"`
	VXLanID     int    `json:"vxlanId"`
}

type VinsList []Vins

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

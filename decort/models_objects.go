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

type DiskConfig struct {
	Label string
	Size int
	Pool string
	Provider string
	ID int
}

type NetworkConfig struct {
	Label string
	NetworkID int
}

type PortforwardConfig struct {
	Label string
	ExtPort int
	IntPort int
	Proto string
}

type SshKeyConfig struct {
	User string
	SshKey string
	UserShell string
}

type ComputeConfig struct {
	ResGroupID int
	Name string
	ID int
	Cpu int
	Ram int
	ImageID int
	BootDisk DiskConfig
	DataDisks []DiskConfig
	Networks []NetworkConfig
	PortForwards []PortforwardConfig
	SshKeys []SshKeyConfig
	Description string
	// The following two parameters are required to create data disks by 
	// a separate disks/create API call
	TenantID int
	GridID int
	// The following one paratmeter is required to create port forwards
	// it will be obsoleted when we implement true Resource Groups
	ExtIP string
}

type ResgroupQuotaConfig struct {
	Cpu int
	Ram float32 // NOTE: it is float32! However, int would be enough here
	Disk int
	NetTraffic int
	ExtIPs int
}

type ResgroupConfig struct {
	TenantID int
	TenantName string
	Location string
	Name string
	ID int
	GridID int
	ExtIP string   // legacy field for VDC - this will eventually become obsoleted by true Resource Groups
	Quota ResgroupQuotaConfig
	Network NetworkConfig
}
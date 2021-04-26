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

// LimitMaxVinsPerResgroup set maximum number of VINs instances per Resource Group
const LimitMaxVinsPerResgroup=4

// MaxSshKeysPerCompute sets maximum number of user:ssh_key pairs to authorize when creating new compute
const MaxSshKeysPerCompute=12

// MaxExtraDisksPerCompute sets maximum number of extra disks that can be added when creating new compute
const MaxExtraDisksPerCompute=12

// MaxNetworksPerCompute sets maximum number of vNICs per compute
const MaxNetworksPerCompute=8

// MaxCpusPerCompute sets maximum number of vCPUs per compute
const MaxCpusPerCompute=128

// MinRamPerCompute sets minimum amount of RAM per compute in MB
const MinRamPerCompute=128
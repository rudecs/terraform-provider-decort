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

package kvmvm

const KvmX86CreateAPI = "/restmachine/cloudbroker/kvmx86/create"
const KvmPPCCreateAPI = "/restmachine/cloudbroker/kvmppc/create"
const ComputeGetAPI = "/restmachine/cloudbroker/compute/get"
const RgListComputesAPI = "/restmachine/cloudbroker/rg/listComputes"
const ComputeNetAttachAPI = "/restmachine/cloudbroker/compute/netAttach"
const ComputeNetDetachAPI = "/restmachine/cloudbroker/compute/netDetach"
const ComputeDiskAttachAPI = "/restmachine/cloudbroker/compute/diskAttach"
const ComputeDiskDetachAPI = "/restmachine/cloudbroker/compute/diskDetach"
const ComputeStartAPI = "/restmachine/cloudbroker/compute/start"
const ComputeStopAPI = "/restmachine/cloudbroker/compute/stop"
const ComputeResizeAPI = "/restmachine/cloudbroker/compute/resize"
const DisksResizeAPI = "/restmachine/cloudbroker/disks/resize2"
const ComputeDeleteAPI = "/restmachine/cloudbroker/compute/delete"

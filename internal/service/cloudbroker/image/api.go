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

package image

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

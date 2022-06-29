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

package account

const accountAddUserAPI = "/restmachine/cloudapi/account/addUser"
const accountAuditsAPI = "/restmachine/cloudapi/account/audits"
const accountCreateAPI = "/restmachine/cloudapi/account/create"
const accountDeleteAPI = "/restmachine/cloudapi/account/delete"
const accountDeleteUserAPI = "/restmachine/cloudapi/account/deleteUser"
const accountDisableAPI = "/restmachine/cloudapi/account/disable"
const accountEnableAPI = "/restmachine/cloudapi/account/enable"
const accountGetAPI = "/restmachine/cloudapi/account/get"
const accountGetConsumedUnitsAPI = "/restmachine/cloudapi/account/getConsumedAccountUnits"
const accountGetConsumedUnitsByTypeAPI = "/restmachine/cloudapi/account/getConsumedCloudUnitsByType"
const accountGetReservedUnitsAPI = "/restmachine/cloudapi/account/getReservedAccountUnits"
const accountListAPI = "/restmachine/cloudapi/account/list"
const accountListComputesAPI = "/restmachine/cloudapi/account/listComputes"
const accountListDeletedAPI = "/restmachine/cloudapi/account/listDeleted"
const accountListDisksAPI = "/restmachine/cloudapi/account/listDisks"
const accountListFlipGroupsAPI = "/restmachine/cloudapi/account/listFlipGroups"
const accountListRGAPI = "/restmachine/cloudapi/account/listRG"
const accountListTemplatesAPI = "/restmachine/cloudapi/account/listTemplates"
const accountListVinsAPI = "/restmachine/cloudapi/account/listVins"
const accountRestoreAPI = "/restmachine/cloudapi/account/restore"
const accountUpdateAPI = "/restmachine/cloudapi/account/update"
const accountUpdateUserAPI = "/restmachine/cloudapi/account/updateUser"

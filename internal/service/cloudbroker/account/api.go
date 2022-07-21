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

const accountAddUserAPI = "/restmachine/cloudbroker/account/addUser"
const accountAuditsAPI = "/restmachine/cloudbroker/account/audits"
const accountCreateAPI = "/restmachine/cloudbroker/account/create"
const accountDeleteAPI = "/restmachine/cloudbroker/account/delete"
const accountDeleteUserAPI = "/restmachine/cloudbroker/account/deleteUser"
const accountDisableAPI = "/restmachine/cloudbroker/account/disable"
const accountEnableAPI = "/restmachine/cloudbroker/account/enable"
const accountGetAPI = "/restmachine/cloudbroker/account/get"
const accountListAPI = "/restmachine/cloudbroker/account/list"
const accountListComputesAPI = "/restmachine/cloudbroker/account/listComputes"
const accountListDeletedAPI = "/restmachine/cloudbroker/account/listDeleted"
const accountListDisksAPI = "/restmachine/cloudbroker/account/listDisks"
const accountListFlipGroupsAPI = "/restmachine/cloudbroker/account/listFlipGroups"
const accountListRGAPI = "/restmachine/cloudbroker/account/listRG"
const accountListVinsAPI = "/restmachine/cloudbroker/account/listVins"
const accountRestoreAPI = "/restmachine/cloudbroker/account/restore"
const accountUpdateAPI = "/restmachine/cloudbroker/account/update"
const accountUpdateUserAPI = "/restmachine/cloudbroker/account/updateUser"

//currently unused
//const accountsEnableAPI = "/restmachine/cloudbroker/account/enableAccounts"
//const accountsDisableAPI = "/restmachine/cloudbroker/account/disableAccounts"
//const accountsDeleteAPI = "/restmachine/cloudbroker/account/deleteAccounts"

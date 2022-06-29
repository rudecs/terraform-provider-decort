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

package bservice

const bserviceCreateAPI = "/restmachine/cloudapi/bservice/create"
const bserviceDeleteAPI = "/restmachine/cloudapi/bservice/delete"
const bserviceDisableAPI = "/restmachine/cloudapi/bservice/disable"
const bserviceEnableAPI = "/restmachine/cloudapi/bservice/enable"
const bserviceGetAPI = "/restmachine/cloudapi/bservice/get"
const bserviceGroupAddAPI = "/restmachine/cloudapi/bservice/groupAdd"
const bserviceGroupComputeRemoveAPI = "/restmachine/cloudapi/bservice/groupComputeRemove"
const bserviceGroupGetAPI = "/restmachine/cloudapi/bservice/groupGet"
const bserviceGroupParentAddAPI = "/restmachine/cloudapi/bservice/groupParentAdd"
const bserviceGroupParentRemoveAPI = "/restmachine/cloudapi/bservice/groupParentRemove"
const bserviceGroupRemoveAPI = "/restmachine/cloudapi/bservice/groupRemove"
const bserviceGroupResizeAPI = "/restmachine/cloudapi/bservice/groupResize"
const bserviceGroupStartAPI = "/restmachine/cloudapi/bservice/groupStart"
const bserviceGroupStopAPI = "/restmachine/cloudapi/bservice/groupStop"
const bserviceGroupUpdateAPI = "/restmachine/cloudapi/bservice/groupUpdate"
const bserviceGroupUpdateExtnetAPI = "/restmachine/cloudapi/bservice/groupUpdateExtnet"
const bserviceGroupUpdateVinsAPI = "/restmachine/cloudapi/bservice/groupUpdateVins"
const bserviceListAPI = "/restmachine/cloudapi/bservice/list"
const bserviceListDeletedAPI = "/restmachine/cloudapi/bservice/listDeleted"
const bserviceRestoreAPI = "/restmachine/cloudapi/bservice/restore"
const bserviceSnapshotCreateAPI = "/restmachine/cloudapi/bservice/snapshotCreate"
const bserviceSnapshotDeleteAPI = "/restmachine/cloudapi/bservice/snapshotDelete"
const bserviceSnapshotListAPI = "/restmachine/cloudapi/bservice/snapshotList"
const bserviceSnapshotRollbackAPI = "/restmachine/cloudapi/bservice/snapshotRollback"
const bserviceStartAPI = "/restmachine/cloudapi/bservice/start"
const bserviceStopAPI = "/restmachine/cloudapi/bservice/stop"

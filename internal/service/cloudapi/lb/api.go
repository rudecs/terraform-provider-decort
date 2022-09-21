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

package lb

const lbListAPI = "/restmachine/cloudapi/lb/list"
const lbListDeletedAPI = "/restmachine/cloudapi/lb/listDeleted"
const lbGetAPI = "/restmachine/cloudapi/lb/get"
const lbCreateAPI = "/restmachine/cloudapi/lb/create"
const lbDeleteAPI = "/restmachine/cloudapi/lb/delete"
const lbDisableAPI = "/restmachine/cloudapi/lb/disable"
const lbEnableAPI = "/restmachine/cloudapi/lb/enable"
const lbUpdateAPI = "/restmachine/cloudapi/lb/update"
const lbStartAPI = "/restmachine/cloudapi/lb/start"
const lbStopAPI = "/restmachine/cloudapi/lb/stop"
const lbRestartAPI = "/restmachine/cloudapi/lb/restart"
const lbRestoreAPI = "/restmachine/cloudapi/lb/restore"
const lbConfigResetAPI = "/restmachine/cloudapi/lb/configReset"
const lbBackendCreateAPI = "/restmachine/cloudapi/lb/backendCreate"
const lbBackendDeleteAPI = "/restmachine/cloudapi/lb/backendDelete"
const lbBackendUpdateAPI = "/restmachine/cloudapi/lb/backendUpdate"
const lbBackendServerAddAPI = "/restmachine/cloudapi/lb/backendServerAdd"
const lbBackendServerDeleteAPI = "/restmachine/cloudapi/lb/backendServerDelete"
const lbBackendServerUpdateAPI = "/restmachine/cloudapi/lb/backendServerUpdate"
const lbFrontendCreateAPI = "/restmachine/cloudapi/lb/frontendCreate"
const lbFrontendDeleteAPI = "/restmachine/cloudapi/lb/frontendDelete"
const lbFrontendBindAPI = "/restmachine/cloudapi/lb/frontendBind"
const lbFrontendBindDeleteAPI = "/restmachine/cloudapi/lb/frontendBindDelete"
const lbFrontendBindUpdateAPI = "/restmachine/cloudapi/lb/frontendBindingUpdate"

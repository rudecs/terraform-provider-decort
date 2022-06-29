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

package sep

const sepAddConsumerNodesAPI = "/restmachine/cloudbroker/sep/addConsumerNodes"
const sepDelConsumerNodesAPI = "/restmachine/cloudbroker/sep/delConsumerNodes"
const sepAddProviderNodesAPI = "/restmachine/cloudbroker/sep/addProviderNodes"

const sepConfigFieldEditAPI = "/restmachine/cloudbroker/sep/configFieldEdit"
const sepConfigInsertAPI = "/restmachine/cloudbroker/sep/configInsert"
const sepConfigValidateAPI = "/restmachine/cloudbroker/sep/configValidate"

const sepConsumptionAPI = "/restmachine/cloudbroker/sep/consumption"

const sepDecommissionAPI = "/restmachine/cloudbroker/sep/decommission"

const sepEnableAPI = "/restmachine/cloudbroker/sep/enable"
const sepDisableAPI = "/restmachine/cloudbroker/sep/disable"

const sepDiskListAPI = "/restmachine/cloudbroker/sep/diskList"

const sepGetAPI = "/restmachine/cloudbroker/sep/get"
const sepGetConfigAPI = "/restmachine/cloudbroker/sep/getConfig"
const sepGetPoolAPI = "/restmachine/cloudbroker/sep/getPool"

const sepCreateAPI = "/restmachine/cloudbroker/sep/create"
const sepDeleteAPI = "/restmachine/cloudbroker/sep/delete"
const sepListAPI = "/restmachine/cloudbroker/sep/list"

const sepUpdateCapacityLimitAPI = "/restmachine/cloudbroker/sep/updateCapacityLimit"

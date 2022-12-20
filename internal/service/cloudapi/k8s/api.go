/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

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

package k8s

const (
	K8sCreateAPI      = "/restmachine/cloudapi/k8s/create"
	K8sGetAPI         = "/restmachine/cloudapi/k8s/get"
	K8sUpdateAPI      = "/restmachine/cloudapi/k8s/update"
	K8sDeleteAPI      = "/restmachine/cloudapi/k8s/delete"
	K8sListAPI        = "/restmachine/cloudapi/k8s/list"
	K8sListDeletedAPI = "/restmachine/cloudapi/k8s/listDeleted"

	K8sWgCreateAPI = "/restmachine/cloudapi/k8s/workersGroupAdd"
	K8sWgDeleteAPI = "/restmachine/cloudapi/k8s/workersGroupDelete"

	K8sWorkerAddAPI    = "/restmachine/cloudapi/k8s/workerAdd"
	K8sWorkerDeleteAPI = "/restmachine/cloudapi/k8s/deleteWorkerFromGroup"

	K8sGetConfigAPI = "/restmachine/cloudapi/k8s/getConfig"

	LbGetAPI = "/restmachine/cloudapi/lb/get"

	AsyncTaskGetAPI = "/restmachine/cloudapi/tasks/get"
)

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

package status

type Status = string

var (
	// The disk is linked to any Compute
	// Status available for:
	//  - Disk
	Assigned Status = "ASSIGNED"

	// An object enabled for operations
	// Status available for:
	//  - Compute
	//  - Disk
	//  - Vins
	Enabled Status = "ENABLED"

	// Enabling in process
	// Status available for:
	//  - Disk
	//  - Vins
	Enabling Status = "ENABLING"

	// An object disabled for operations
	// Status available for:
	//  - Compute
	//  - Disk
	//  - Vins
	Disabled Status = "DISABLED"

	// Disabling in process
	// Status available for:
	//  - Disk
	//  - Vins
	Disabling Status = "DISABLING"

	// An object model has been created in the database
	// Status available for:
	//  - Image
	//  - Disk
	//  - Compute
	//  - Vins
	Modeled Status = "MODELED"

	// In the process of creation
	// Status available for:
	//  - Image
	Creating Status = "CREATING"

	// An object was created successfully
	// Status available for:
	//  - Image
	//  - Disk
	//  - Compute
	//  - Vins
	Created Status = "CREATED"

	// Physical resources are allocated for the object
	// Status available for:
	//  - Compute
	Allocated Status = "ALLOCATED"

	// The object has released (returned to the platform) the physical resources that it occupied
	// Status available for:
	//  - Compute
	Unallocated Status = "UNALLOCATED"

	// Destroying in progress
	// Status available for:
	//  - Disk
	//  - Compute
	//  - Vins
	Destroying Status = "DESTROYING"

	// Permanently deleted
	// Status available for:
	//  - Image
	//  - Disk
	//  - Compute
	//  - Vins
	Destroyed Status = "DESTROYED"

	// Deleting in progress to Trash
	// Status available for:
	//  - Compute
	//  - Vins
	Deleting Status = "DELETING"

	// Deleted to Trash
	// Status available for:
	//  - Compute
	//  - Vins
	Deleted Status = "DELETED"

	// Deleted from storage
	// Status available for:
	//  - Image
	Purged Status = "PURGED"

	// Repeating deploy of the object in progress
	// Status available for:
	//  - Compute
	Redeploying Status = "REDEPLOYING"

	// The resource is not bound to vnf device
	// Status available for:
	//  - vins vnf
	Stashed Status = "STASHED"
)

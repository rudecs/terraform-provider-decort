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

package techstatus

type TechStatus = string

var (
	// Start in progress - send an execution command
	// Status available for:
	//  - Compute
	Starting TechStatus = "STARTING"

	// An object started
	// Can be stopped
	// Correctly running
	// Status available for:
	//  - Compute
	Started TechStatus = "STARTED"

	// Stop in progress - send an execution command
	// Status available for:
	//  - Compute
	Stopping TechStatus = "STOPPING"

	// An object stopped
	// Can be started
	// Limited functionality
	// Status available for:
	//  - Compute
	Stopped TechStatus = "STOPPED"

	// Pause in progress - send an execution command
	// Status available for:
	//  - Compute
	Pausing TechStatus = "PAUSING"

	// An object paused
	// Can be restarted
	// Currently running
	// Status available for:
	//  - Compute
	Paused TechStatus = "PAUSED"

	// Migrate in progress
	// Status available for:
	//  - Compute
	Migrating TechStatus = "MIGRATING"

	// An object failure status
	// Can be reastarted
	// Limited functionality
	// Status available for:
	//  - Compute
	Down TechStatus = "DOWN"

	// An object configuration process
	// Status available for:
	//  - Compute
	Scheduled TechStatus = "SCHEDULED"

	// Physical resources are allocated for the object
	// Status available for:
	//  - Image
	Allocated TechStatus = "ALLOCATED"

	// The object has released (returned to the platform) the physical resources that it occupied
	// Status available for:
	//  - Image
	Unallocated TechStatus = "UNALLOCATED"
)

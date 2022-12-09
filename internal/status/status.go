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
	Enabled Status = "ENABLED"

	// Enabling in process
	// Status available for:
	//  - Disk
	Enabling Status = "ENABLING"

	// An object disabled for operations
	// Status available for:
	//  - Compute
	//  - Disk
	Disabled Status = "DISABLED"

	// Disabling in process
	// Status available for:
	//  - Disk
	Disabling Status = "DISABLING"

	// An object model has been created in the database
	// Status available for:
	//  - Image
	//  - Disk
	//  - Compute
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
	Destroying Status = "DESTROYING"

	// Permanently deleted
	// Status available for:
	//  - Image
	//  - Disk
	//  - Compute
	Destroyed Status = "DESTROYED"

	// Deleting in progress to Trash
	// Status available for:
	//  - Compute
	Deleting Status = "DELETING"

	// Deleted to Trash
	// Status available for:
	//  - Compute
	Deleted Status = "DELETED"

	// Deleted from storage
	// Status available for:
	//  - Image
	Purged Status = "PURGED"

	// Repeating deploy of the object in progress
	// Status available for:
	//  - Compute
	Redeploying Status = "REDEPLOYING"
)

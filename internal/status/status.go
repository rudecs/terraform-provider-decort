package status

type Status = string

var (
	//The disk is linked to any Compute
	Assigned Status = "ASSIGNED"

	//An object model has been created in the database
	Modeled Status = "MODELED"

	//In the process of creation
	Creating Status = "CREATING"

	//Creating
	Created Status = "CREATED"

	//Physical resources are allocated for the object
	Allocated Status = "ALLOCATED"

	//The object has released (returned to the platform) the physical resources that it occupied
	Unallocated Status = "UNALLOCATED"

	//Permanently deleted
	Destroyed Status = "DESTROYED"

	//Deleted to Trash
	Deleted Status = "DELETED"

	//Deleted from storage
	Purged Status = "PURGED"
)

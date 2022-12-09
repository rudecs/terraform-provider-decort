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

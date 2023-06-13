package specs

type Sync struct {
	Concurrency uint64 `json:"concurrency,omitempty"`
	// Tables to sync from the source plugin
	Tables []string `json:"tables,omitempty"`
	// SkipTables defines tables to skip when syncing data. Useful if a glob pattern is used in Tables
	SkipTables []string `json:"skip_tables,omitempty"`
	// SkipDependentTables changes the matching behavior with regard to dependent tables. If set to true, dependent tables will not be synced unless they are explicitly matched by Tables.
	SkipDependentTables bool `json:"skip_dependent_tables,omitempty"`
	// Destinations are the names of destination plugins to send sync data to
	Destinations []string `json:"destinations,omitempty"`
	// Scheduler defines the scheduling algorithm that should be used to sync data
	Scheduler Scheduler `json:"scheduler,omitempty"`
	// DeterministicCQID is a flag that indicates whether the source plugin should generate a random UUID as the value of _cq_id
	// or whether it should calculate a UUID that is a hash of the primary keys (if they exist) or the entire resource.
	DeterministicCQID bool `json:"deterministic_cq_id,omitempty"`
}

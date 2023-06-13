package specs

type Write struct {
	WriteMode      WriteMode   `json:"write_mode,omitempty"`
	MigrateMode    MigrateMode `json:"migrate_mode,omitempty"`
	BatchSize      int         `json:"batch_size,omitempty"`
	BatchSizeBytes int         `json:"batch_size_bytes,omitempty"`
	PKMode         PKMode      `json:"pk_mode,omitempty"`
}

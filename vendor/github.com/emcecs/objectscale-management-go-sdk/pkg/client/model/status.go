package model

// RebuildInfo is the information about the rebuild status
type RebuildInfo struct {
	// Status of the storage server
	Status string `json:"status"`

	// TotalBytes of the storage server
	TotalBytes int `json:"total_bytes,string"`

	// RemainingBytes of the storage server
	RemainingBytes int `json:"remaining_bytes,string"`

	// Level of the storage server
	Level int `json:"level,string"`

	// Disk of the storage node
	Disk string `json:"disk"`

	// Message from the storage server
	Message string `json:"message"`

	// Host of the storage server
	Host string `json:"host"`

	// Progress of the recovery
	Progress string `json:"progress"`
}

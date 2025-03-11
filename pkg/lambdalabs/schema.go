package lambdalabs

type Instance struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	IP              string   `json:"ip"`
	Status          string   `json:"status"`
	FileSystemNames []string `json:"file_system_names,omitempty"`
}

type Region struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type InstanceTypeSpecs struct {
	VCPUs      int `json:"vcpus"`
	MemoryGiB  int `json:"memory_gib"`
	StorageGiB int `json:"storage_gib"`
	GPUs       int `json:"gpus"`
}

type InstanceType struct {
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	GPUDescription    string            `json:"gpu_description"`
	PriceCentsPerHour int               `json:"price_cents_per_hour"`
	Specs             InstanceTypeSpecs `json:"specs"`
}

type InstanceTypeInfo struct {
	InstanceType                 InstanceType `json:"instance_type"`
	RegionsWithCapacityAvailable []Region     `json:"regions_with_capacity_available"`
}

// User represents a Lambda Labs user
type User struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// FileSystem represents a Lambda Labs file system
type FileSystem struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	MountPoint string  `json:"mount_point"`
	Created    string  `json:"created"`
	CreatedBy  User    `json:"created_by"`
	IsInUse    bool    `json:"is_in_use"`
	Region     Region  `json:"region"`
	BytesUsed  int64   `json:"bytes_used"`
}

// Image represents a Lambda Labs image
type Image struct {
	ID           string  `json:"id"`
	CreatedTime  string  `json:"created_time"`
	UpdatedTime  string  `json:"updated_time"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Family       string  `json:"family"`
	Version      string  `json:"version"`
	Architecture string  `json:"architecture"`
	Region       Region  `json:"region"`
}

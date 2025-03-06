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

package lambdalabs

type Instance struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	IP              string   `json:"ip"`
	Status          string   `json:"status"`
	FileSystemNames []string `json:"file_system_names,omitempty"`
}

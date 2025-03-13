package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// filesystemRegionModel represents a region for a filesystem
type filesystemRegionModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// filesystemUserModel represents a user who created a filesystem
type filesystemUserModel struct {
	ID     types.String `tfsdk:"id"`
	Email  types.String `tfsdk:"email"`
	Status types.String `tfsdk:"status"`
}

// filesystemDataModel represents a filesystem with all its attributes
type filesystemDataModel struct {
	ID         types.String           `tfsdk:"id"`
	Name       types.String           `tfsdk:"name"`
	MountPoint types.String           `tfsdk:"mount_point"`
	Created    types.String           `tfsdk:"created"`
	CreatedBy  *filesystemUserModel   `tfsdk:"created_by"`
	IsInUse    types.Bool             `tfsdk:"is_in_use"`
	Region     *filesystemRegionModel `tfsdk:"region"`
	BytesUsed  types.Int64            `tfsdk:"bytes_used"`
}

// filesystemResourceModel represents a filesystem model for resource operations
// containing fields for creation, management and display
type filesystemResourceModel struct {
	ID         types.String          `tfsdk:"id"`
	Name       types.String          `tfsdk:"name"`
	Region     types.String          `tfsdk:"region"`
	MountPoint types.String          `tfsdk:"mount_point"`
	Created    types.String          `tfsdk:"created"`
	IsInUse    types.Bool            `tfsdk:"is_in_use"`
	BytesUsed  types.Int64           `tfsdk:"bytes_used"`
	CreatedBy  filesystemUserModel   `tfsdk:"created_by"`
	RegionInfo filesystemRegionModel `tfsdk:"region_info"`
}

// filesystemsFilterModel represents filtering options for filesystems
type filesystemsFilterModel struct {
	Region types.String `tfsdk:"region"`
}

// filesystemsDataModel represents the data source model for filesystems
type filesystemsDataModel struct {
	ID          types.String            `tfsdk:"id"`
	Filter      *filesystemsFilterModel `tfsdk:"filter"`
	FileSystems []filesystemDataModel   `tfsdk:"filesystems"`
}

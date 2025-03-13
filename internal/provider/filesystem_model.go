package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// 共用的檔案系統區域模型
type filesystemRegionModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// 共用的檔案系統使用者模型
type filesystemUserModel struct {
	ID     types.String `tfsdk:"id"`
	Email  types.String `tfsdk:"email"`
	Status types.String `tfsdk:"status"`
}

// 共用的檔案系統資料模型
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

// 檔案系統資源模型 (簡化版本，只包含必要欄位)
type filesystemResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
}

// 檔案系統過濾器模型
type filesystemsFilterModel struct {
	Region types.String `tfsdk:"region"`
}

// 檔案系統資料來源模型
type filesystemsDataModel struct {
	ID          types.String            `tfsdk:"id"`
	Filter      *filesystemsFilterModel `tfsdk:"filter"`
	FileSystems []filesystemDataModel   `tfsdk:"filesystems"`
}

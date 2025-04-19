package models

import (
	"nostalgia/internal/common/rbac"
	"nostalgia/internal/common/util"
)

type UserPermissions struct {
	CanViewMedia         bool `json:"can_view_media"`
	CanUploadMedia       bool `json:"can_upload_media"`
	CanDeleteMedia       bool `json:"can_delete_media"`
	CanUpdateMedia       bool `json:"can_update_media"`
	Whitelisted          bool `json:"whitelisted"`
	CanManageWhitelist   bool `json:"can_manage_whitelist"`
	CanManagePermissions bool `json:"can_manage_permissions"`
}

func NewUserPermissionsFromRoles(roles []string) UserPermissions {
	return UserPermissions{
		CanViewMedia:         util.Contains[string](roles, rbac.RoleCanViewMedia),
		CanUploadMedia:       util.Contains[string](roles, rbac.RoleCanUploadMedia),
		CanDeleteMedia:       util.Contains[string](roles, rbac.RoleCanDeleteMedia),
		CanUpdateMedia:       util.Contains[string](roles, rbac.RoleCanUpdateMedia),
		Whitelisted:          util.Contains[string](roles, rbac.RoleWhitelisted),
		CanManagePermissions: util.Contains[string](roles, rbac.RoleCanManagePermissions),
	}
}

func (p UserPermissions) RolesFromUserPermissions() ([]string, []string) {
	var match = map[string]bool{
		rbac.RoleCanViewMedia:         p.CanViewMedia,
		rbac.RoleCanUploadMedia:       p.CanUploadMedia,
		rbac.RoleCanDeleteMedia:       p.CanDeleteMedia,
		rbac.RoleCanUpdateMedia:       p.CanUpdateMedia,
		rbac.RoleWhitelisted:          p.Whitelisted,
		rbac.RoleCanManagePermissions: p.CanManagePermissions,
	}
	var add []string
	var remove []string
	for role, shouldAdd := range match {
		if shouldAdd {
			add = append(add, role)
		} else {
			remove = append(remove, role)
		}
	}
	return add, remove
}

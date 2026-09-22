package authz

import (
	"latihan-fiber/app/model"
	"latihan-fiber/helper"
	"strings"
)

func CanAccessUser(
	current model.AuthUser,
	targetID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.UserID == targetID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}

func ValidateAssignRole(
	current model.AuthUser,
	targetID int,
	req model.AssignRoleRequest,
	perms *helper.PermissionSet,
) map[string]string {
	errs := map[string]string{}
	role := strings.TrimSpace(req.Role)
	if role == "" {
		errs["role"] = "Wajib diisi"
		return errs
	}
	if !perms.IsKnownRole(role) {
		errs["role"] = "Unknown role, choose one: " +
			strings.Join(perms.KnownRoles(), ", ")
	}
	if current.UserID == targetID {
		errs["role"] = "Can't change your roles"
	}
	return errs
}

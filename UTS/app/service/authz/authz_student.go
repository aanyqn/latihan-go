package authz

import (
	"siakad/app/model"
	"siakad/helper"
)

func CanAccessStudent(
	current model.AuthUser,
	userID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.UserID == userID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}
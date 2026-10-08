package user

import (
	"siakad/app/model"
	"strings"
)

func ApplyPatch(current model.User, req model.PatchUserRequest) model.User {
	if req.Email != nil {
		current.Email = strings.TrimSpace(*req.Email)
	}
	return current
}

func IsEmptyPatch(req model.PatchUserRequest) bool {
	return req.Email == nil
}

func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1
}

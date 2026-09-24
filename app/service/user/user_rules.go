package user

import (
	"latihan-fiber/app/model"
	"strings"
)

// func ValidateCreate(req model.CreateUserRequest) map[string]string {
// 	errs := map[string]string{}

// 	if strings.TrimSpace(req.Username) == "" {
// 		errs["username"] = "Must be filled"
// 	}
// 	if !isValidEmail(req.Email) {
// 		errs["email"] = "Must be valid"
// 	}
// 	if len(req.Password) < 8 {
// 		errs["password"] = "8 characters required"
// 	}
// 	return errs
// }

// func ValidateReplace(req model.ReplaceUserRequest) map[string]string {
// 	errs := map[string]string{}
// 	if strings.TrimSpace(req.Username) == "" {
// 		errs["username"] = "must be filled"
// 	}
// 	if !isValidEmail(req.Email) {
// 		errs["email"] = "must be valid"
// 	}
// 	return errs
// }

func ApplyPatch(current model.User, req model.PatchUserRequest) model.User {
	if req.Username != nil {
		current.Username = strings.TrimSpace(*req.Username)
	}
	if req.Email != nil {
		current.Email = strings.TrimSpace(*req.Email)
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}
	return current
}

func IsEmptyPatch(req model.PatchUserRequest) bool {
	return req.Username == nil && req.Email == nil && req.IsActive == nil
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

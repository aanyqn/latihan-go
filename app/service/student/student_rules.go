package student

import (
	"latihan-fiber/app/model"
	"strings"
)

func ValidateCreate(req model.CreateStudentRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "Must be filled"
	}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "Must be filled"
	}
	return errs
}

func ValidateReplace(req model.ReplaceStudentRequest) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(*req.Username) == "" {
		errs["username"] = "must be filled"
	}
	if strings.TrimSpace(*req.NIM) == "" {
		errs["nim"] = "must be filled"
	}
	if req.Grade == nil {
		errs["grade"] = "must be filled"
	}
	return errs
}

func ApplyPatch(
	current model.Student, req model.PatchStudentRequest,
) (model.Student, map[string]string) {
	errs := map[string]string{}

	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			errs["username"] = "must be filled"
		} else {
			current.Username = *req.Username
		}
	}
	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			errs["nim"] = "must be filled"
		} else {
			current.NIM = *req.NIM
		}
	}
	
	if req.Grade != nil {
		current.Grade = *req.Grade
	}
	return current, errs
}

func IsEmptyPatch(req model.PatchStudentRequest) bool {
	return req.Username == nil && req.NIM == nil && req.IsActive == nil
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

package student_test

import (
	"latihan-fiber/app/model"
	"latihan-fiber/app/service/student"
	"testing"
)

func TestValidateCreate(t *testing.T) {
	t.Run("Valid request", func(t *testing.T) {
		req := model.CreateStudentRequest{
			Username: "johndoe",
			NIM:      "123456",
		}
		errs := student.ValidateCreate(req)
		if len(errs) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(errs))
		}
	})

	t.Run("Empty fields", func(t *testing.T) {
		req := model.CreateStudentRequest{
			Username: "",
			NIM:      "   ",
		}
		errs := student.ValidateCreate(req)
		if len(errs) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(errs))
		}
		if errs["username"] != "Must be filled" {
			t.Errorf("Expected username error 'Must be filled', got '%s'", errs["username"])
		}
		if errs["nim"] != "Must be filled" {
			t.Errorf("Expected nim error 'Must be filled', got '%s'", errs["nim"])
		}
	})
}

func TestValidateReplace(t *testing.T) {
	t.Run("Valid request", func(t *testing.T) {
		username := "jhondoe"
		nim := "232324"
		grade := 3.5
		req := model.ReplaceStudentRequest{
			Username: &username,
			NIM:      &nim,
			Grade:    &grade,
		}
		errs := student.ValidateReplace(req)
		if len(errs) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(errs))
		}
	})

	t.Run("Missing grade and empty username", func(t *testing.T) {
		username := "  "
		nim := "123"
		req := model.ReplaceStudentRequest{
			Username: &username,
			NIM:      &nim,
			Grade:    nil,
		}
		errs := student.ValidateReplace(req)
		if len(errs) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(errs))
		}
		if errs["grade"] != "must be filled" {
			t.Errorf("Expected grade error 'must be filled', got '%s'", errs["grade"])
		}
		if errs["username"] != "must be filled" {
			t.Errorf("Expected username error 'must be filled', got '%s'", errs["username"])
		}
	})
}

func TestApplyPatch(t *testing.T) {
	t.Run("Valid patch", func(t *testing.T) {
		current := model.Student{
			Username: "old_user",
			NIM:      "111",
			Grade:    2.0,
		}
		newUsername := "new_user"
		newGrade := 4.0
		req := model.PatchStudentRequest{
			Username: &newUsername,
			Grade:    &newGrade,
		}

		updated, errs := student.ApplyPatch(current, req)
		if len(errs) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(errs))
		}
		if updated.Username != "new_user" {
			t.Errorf("Expected username 'new_user', got '%s'", updated.Username)
		}
		if updated.Grade != 4.0 {
			t.Errorf("Expected grade 4.0, got %f", updated.Grade)
		}
		// NIM should remain unchanged
		if updated.NIM != "111" {
			t.Errorf("Expected nim '111', got '%s'", updated.NIM)
		}
	})

	t.Run("Invalid patch with empty fields", func(t *testing.T) {
		current := model.Student{
			Username: "old_user",
			NIM:      "111",
		}
		emptyStr := "   "
		req := model.PatchStudentRequest{
			Username: &emptyStr,
			NIM:      &emptyStr,
		}

		_, errs := student.ApplyPatch(current, req)
		if len(errs) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(errs))
		}
		if errs["username"] != "must be filled" {
			t.Errorf("Expected username error 'must be filled', got '%s'", errs["username"])
		}
		if errs["nim"] != "must be filled" {
			t.Errorf("Expected nim error 'must be filled', got '%s'", errs["nim"])
		}
	})
}

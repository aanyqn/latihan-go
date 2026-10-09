package enrollment

import (
	"context"
	"errors"
	"siakad/app/model"
	"siakad/app/repository"
	"siakad/helper"
	"strings"
)

type EnrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	studentRepo repository.StudentRepository
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository, studentRepo repository.StudentRepository) *EnrollmentService {
	return &EnrollmentService{enrollmentRepo: enrollmentRepo, studentRepo: studentRepo}
}

func (e *EnrollmentService) Create(ctx context.Context, req model.CreateEnrollmentRequest, user_id int) (model.Enrollment, error) {
	req.TahunAkademik = strings.TrimSpace(req.TahunAkademik)

	student, err := e.studentRepo.FindByUserID(ctx, user_id)
	if err != nil {
		return model.Enrollment{}, helper.Forbidden("Bukan mahasiswa")
	}

	ipk := 0.0
	if student.IPKTerakhir != nil {
		ipk = *student.IPKTerakhir
	}

	enrollment, err := e.enrollmentRepo.Create(ctx, model.Enrollment{
		StudentID:     student.ID,
		CourseID:      req.CourseID,
		TahunAkademik: req.TahunAkademik,
	}, ipk)

	if err != nil {
		if errors.Is(err, helper.ErrDuplicate) {
			return model.Enrollment{}, helper.ErrDuplicate
		}
		return model.Enrollment{}, err
	}

	return enrollment, nil
}

func (e *EnrollmentService) Delete(ctx context.Context, id int, userID int) error {
	student, err := e.studentRepo.FindByUserID(ctx, userID)
	if err != nil {
		return helper.Forbidden("Bukan mahasiswa")
	}

	enrollment, err := e.enrollmentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return helper.ErrNotFound
		}
		return err
	}

	if enrollment.StudentID != student.ID {
		return helper.Forbidden("Tidak dapat menghapus KRS milik mahasiswa lain")
	}

	err = e.enrollmentRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	return nil
}

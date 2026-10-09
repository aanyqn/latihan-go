package student

import (
	"context"
	"errors"
	"strings"

	"siakad/app/model"
	"siakad/app/repository"
	"siakad/helper"
)

type StudentService struct {
	studentRepo repository.StudentRepository
	userRepo    repository.UserRepository
}

func NewStudentService(studentRepo repository.StudentRepository, userRepo repository.UserRepository) *StudentService {
	return &StudentService{studentRepo: studentRepo, userRepo: userRepo}
}

func (s *StudentService) List(ctx context.Context, q helper.ListQuery) ([]model.Student, *helper.Meta, error) {
	students, total, err := s.studentRepo.FindAll(ctx, q)
	if err != nil {
		return nil, nil, err
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	meta := &helper.Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return students, meta, nil
}

func (s *StudentService) Get(ctx context.Context, id int) (model.Student, error) {
	student, err := s.studentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return model.Student{}, helper.ErrNotFound
		}
		return model.Student{}, err
	}
	return student, nil
}

func (s *StudentService) Create(ctx context.Context, req model.CreateStudentRequest) (model.Student, error) {
	req.Nama = strings.TrimSpace(req.Nama)

	req.Email = strings.TrimSpace(req.Email)

	user, err := s.userRepo.Create(ctx, model.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     "mahasiswa",
	})

	if err != nil {
		if errors.Is(err, helper.ErrDuplicate) {
			return model.Student{}, helper.ErrDuplicate
		}
		return model.Student{}, err
	}

	if req.Nama == "" || req.NIM == "" {
		return model.Student{}, helper.ErrInvalidInput
	}

	newStudent, err := s.studentRepo.Create(ctx, model.Student{
		Nama:     req.Nama,
		NIM:      req.NIM,
		Prodi:    req.Prodi,
		Angkatan: req.Angkatan,
		User:     user,
	})

	if err != nil {
		if errors.Is(err, helper.ErrDuplicate) {
			return model.Student{}, helper.ErrDuplicate
		}
		return model.Student{}, err
	}

	return newStudent, nil
}

func (s *StudentService) Replace(ctx context.Context, id int, req model.ReplaceStudentRequest) (model.Student, error) {
	if req.Nama == nil || req.NIM == nil || req.Prodi == nil || req.Angkatan == nil || req.IPKTerakhir == nil {
		return model.Student{}, helper.ErrInvalidInput
	}

	nama := strings.TrimSpace(*req.Nama)
	nim := strings.TrimSpace(*req.NIM)

	updated, err := s.studentRepo.Update(ctx, model.Student{
		ID:          id,
		Nama:        nama,
		NIM:         nim,
		Prodi:       *req.Prodi,
		Angkatan:    *req.Angkatan,
		IPKTerakhir: req.IPKTerakhir,
	})
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return model.Student{}, helper.ErrNotFound
		}
		if errors.Is(err, helper.ErrDuplicate) {
			return model.Student{}, helper.ErrDuplicate
		}
		return model.Student{}, err
	}

	return updated, nil
}

func (s *StudentService) Delete(ctx context.Context, id int) error {
	err := s.studentRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	return nil
}

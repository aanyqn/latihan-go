package student

import (
	"context"
	"errors"
	"strings"

	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
	"latihan-fiber/helper"
)

var (
	ErrNotFound       = errors.New("student not found")
	ErrDuplicate      = errors.New("username already used")
	ErrInvalidInput   = errors.New("input isn't valid")
	ErrNoFieldsChange = errors.New("no changes input")
)

type StudentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) List(ctx context.Context, q helper.ListQuery) ([]model.Student, *helper.Meta, error) {
	students, total, err := s.repo.FindAll(ctx, q)
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
	student, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, err
	}
	return student, nil
}

func (s *StudentService) Create(ctx context.Context, req model.CreateStudentRequest) (model.Student, error) {
	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.NIM == "" {
		return model.Student{}, ErrInvalidInput
	}

	newStudent, err := s.repo.Create(ctx, model.Student{
		Username: req.Username,
		NIM:      req.NIM,
		IsActive: true,
	})

	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, err
	}

	return newStudent, nil
}

func (s *StudentService) Replace(ctx context.Context, id int, req model.ReplaceStudentRequest) (model.Student, error) {
	if req.Username == nil || req.NIM == nil || req.Grade == nil || req.IsActive == nil {
		return model.Student{}, ErrInvalidInput
	}

	username := strings.TrimSpace(*req.Username)
	nim := strings.TrimSpace(*req.NIM)

	if username == "" || nim == "" {
		return model.Student{}, ErrInvalidInput
	}

	updated, err := s.repo.Update(ctx, model.Student{
		ID:       id,
		Username: username,
		NIM:      nim,
		Grade:    *req.Grade,
		IsActive: *req.IsActive,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Student{}, ErrNotFound
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, err
	}

	return updated, nil
}

func (s *StudentService) Patch(ctx context.Context, id int, req model.PatchStudentRequest) (model.Student, error) {
	if req.Username == nil && req.NIM == nil && req.Grade == nil && req.IsActive == nil {
		return model.Student{}, ErrNoFieldsChange
	}

	current, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, err
	}

	if req.Username != nil {
		val := strings.TrimSpace(*req.Username)
		if val == "" {
			return model.Student{}, ErrInvalidInput
		}
		current.Username = val
	}
	if req.NIM != nil {
		val := strings.TrimSpace(*req.NIM)
		if val == "" {
			return model.Student{}, ErrInvalidInput
		}
		current.NIM = val
	}
	if req.Grade != nil {
		current.Grade = *req.Grade
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	updated, err := s.repo.Update(ctx, current)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, err
	}

	return updated, nil
}

func (s *StudentService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
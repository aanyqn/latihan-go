package course

import (
	"context"
	"errors"
	"siakad/app/model"
	"siakad/app/repository"
	"siakad/helper"
)

var (
	ErrNotFound = errors.New("student not found")
)

type CourseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (c *CourseService) List(ctx context.Context, q helper.ListQuery) ([]model.Course, *helper.Meta, error) {
	courses, total, err := c.repo.FindAll(ctx, q)
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

	return courses, meta, nil
}
package handler

import (
	"context"
	"errors"
	"strconv"
	"latihan-fiber/app/model"
	"latihan-fiber/app/service/student"
	"latihan-fiber/helper"

	"github.com/gofiber/fiber/v2"
)

type StudentHandler struct {
	svc *student.StudentService
}

func NewStudentHandler(svc *student.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

func (h *StudentHandler) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := helper.ParseListQuery(c)

	students, meta, err := h.svc.List(ctx, q)
	if err != nil {
		return translateError(c, err, "Error fetching students")
	}

	return helper.SuccessList(c, "Successful fetching Students data", students, meta)
}

func (h *StudentHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid")
	}

	student, err := h.svc.Get(ctx, id)
	if err != nil {
		return translateError(c, err, "Fail to get student data")
	}

	return helper.Success(c, fiber.StatusOK, "Student found!", student)
}

func (h *StudentHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format!")
	}

	student, err := h.svc.Create(ctx, req)
	if err != nil {
		return translateError(c, err, "Failed to save student data")
	}

	return helper.Created(c, "Student succesfully Created", student, "/api/v1/students/"+strconv.Itoa(student.ID))
}

func (h *StudentHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Body must be valid in JSON format")
	}

	student, err := h.svc.Replace(ctx, id, req)
	if err != nil {
		return translateError(c, err, "Failed to update student")
	}

	return helper.Success(c, fiber.StatusOK, "Student data succesfully replaced", student)
}

func (h *StudentHandler) Patch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format")
	}

	student, err := h.svc.Patch(ctx, id, req)
	if err != nil {
		return translateError(c, err, "Failed to update student")
	}

	return helper.Success(c, fiber.StatusOK, "student successfuly patch updated", student)
}

func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		return translateError(c, err, "Failed to delete student")
	}

	return helper.NoContent(c)
}

func translateError(c *fiber.Ctx, err error, generalMessage string) error {
	switch {
	case errors.Is(err, student.ErrNotFound):
		return helper.Fail(c, fiber.StatusNotFound, err.Error())
	case errors.Is(err, student.ErrDuplicate):
		return helper.Fail(c, fiber.StatusConflict, err.Error())
	case errors.Is(err, student.ErrInvalidInput), errors.Is(err, student.ErrNoFieldsChange):
		return helper.Fail(c, fiber.StatusBadRequest, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return helper.Fail(c, fiber.StatusGatewayTimeout, "Request timeout")
	default:
		return helper.Fail(c, fiber.StatusInternalServerError, generalMessage)
	}
}
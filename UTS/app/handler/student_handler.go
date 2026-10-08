package handler

import (
	"context"
	"errors"
	"siakad/app/model"
	"siakad/app/service/authz"
	"siakad/app/service/student"
	"siakad/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type StudentHandler struct {
	svc   *student.StudentService
	perms *helper.PermissionSet
}

func NewStudentHandler(svc *student.StudentService, perms *helper.PermissionSet) *StudentHandler {
	return &StudentHandler{svc: svc, perms: perms}
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

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID must be valid")
	}

	student, err := h.svc.Get(ctx, id)
	if err != nil {
		return translateError(c, err, "Fail to get student data")
	}

	if !authz.CanAccessStudent(current, student.UserID, h.perms, "student:read:any") {
		return helper.Forbidden("You don't have access")
	}

	return helper.Success(c, fiber.StatusOK, "Student found!", student)
}

func (h *StudentHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	_, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	var req model.CreateStudentRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body must be in valid JSON format!")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Internal(err)
	}

	req.Password = hashed

	student, err := h.svc.Create(ctx, req)
	if err != nil {
		return translateError(c, err, "Failed to save student data")
	}

	return helper.Created(c, "Student succesfully Created", student, "/api/v1/students/"+strconv.Itoa(student.ID))
}

func (h *StudentHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	if current.Role == "user" {
		return helper.Forbidden("You doesn't have right to do this.")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID must be valid!")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body must be valid in JSON format")
	}

	student, err := h.svc.Replace(ctx, id, req)
	if err != nil {
		return translateError(c, err, "Failed to update student")
	}

	if !authz.CanAccessStudent(current, student.UserID, h.perms, "student:update:any") {
		return helper.Forbidden("You don't have access")
	}

	return helper.Success(c, fiber.StatusOK, "Student data succesfully replaced", student)
}

func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID must be valid!")
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

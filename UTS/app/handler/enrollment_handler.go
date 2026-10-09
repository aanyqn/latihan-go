package handler

import (
	"siakad/app/model"
	"siakad/app/service/enrollment"
	"siakad/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EnrollmentHandler struct {
	svc   *enrollment.EnrollmentService
	perms *helper.PermissionSet
}

func NewEnrollmentHandler(svc *enrollment.EnrollmentService, perms *helper.PermissionSet) *EnrollmentHandler {
	return &EnrollmentHandler{svc: svc, perms: perms}
}

func (e *EnrollmentHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	var req model.CreateEnrollmentRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body must be in valid JSON format!")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	enrollment, err := e.svc.Create(ctx, req, current.UserID)
	if err != nil {
		return helper.TranslateError(c, err, "Failed to save enrollment data")
	}

	return helper.Created(c, "Enrollment successfully created", enrollment, "/api/v1/enrollments/"+strconv.Itoa(enrollment.ID))
}

func (e *EnrollmentHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID must be valid!")
	}

	if err := e.svc.Delete(ctx, id); err != nil {
		return helper.TranslateError(c, err, "Failed to delete enrollment")
	}

	return helper.NoContent(c)
}
package handler

import (
	"siakad/app/service/course"
	"siakad/helper"

	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	svc   *course.CourseService
	perms *helper.PermissionSet
}

func NewCourseHandler(svc *course.CourseService, perms *helper.PermissionSet) *CourseHandler {
	return &CourseHandler{svc:svc, perms: perms}
}

func (h *CourseHandler) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := helper.ParseListQuery(c)

	courses, meta, err := h.svc.List(ctx, q)
	if err != nil {
		return helper.TranslateError(c, err, "Error fetching course")
	}

	return helper.SuccessList(c, "Successful fetching course data", courses, meta)
}
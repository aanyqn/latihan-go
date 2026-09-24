package user

import (
	"errors"
	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
	"latihan-fiber/helper"
	"strconv"
	"strings"

	"latihan-fiber/app/service/authz"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	repo  repository.UserRepository
	perms *helper.PermissionSet
}

func NewUserService(repo repository.UserRepository, perms *helper.PermissionSet) *UserService {
	return &UserService{repo: repo, perms: perms}
}

func (h *UserService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := helper.ParseListQuery(c)

	users, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return translateError(err, "users")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	return helper.SuccessList(c, "Successful fetching Users data", users, &helper.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func (h *UserService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.BadRequest("Invalid ID")
	}

	if current.Role == "user" && current.UserID != targetID {
		return helper.Forbidden("You doesn't have right to do this.")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID must be valid")
	}

	if !authz.CanAccessUser(current, id, h.perms, "user:read:any") {
		return helper.Forbidden(
			"tidak berhak mengakses data user lain")
	}

	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(err, "user")
	}

	return helper.Success(c, fiber.StatusOK, "User found!", user)
}

func (h *UserService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body must be in valid JSON format!")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	baru, err := h.repo.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: true,
	})
	if err != nil {
		return translateError(err, "user")
	}
	return helper.Created(c, "User succesfully Created", baru,
		"/api/v1/users/"+strconv.Itoa(baru.ID))
}

func (h *UserService) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.BadRequest("Invalid ID")
	}

	if (current.Role == "user" || current.Role == "staff") && current.UserID != targetID {
		return helper.Forbidden("You doesn't have right to do this.")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID must be valid!")
	}
	var req model.ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body must be valid in JSON format")
	}
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	hasil, err := h.repo.Update(ctx, model.User{
		ID: id, Username: req.Username, Email: req.Email, IsActive: req.IsActive,
	})
	if err != nil {
		return translateError(err, "users")
	}

	return helper.Success(c, fiber.StatusOK, "User data succesfully replaced", hasil)
}

func (h *UserService) Patch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("ID must be valid!")
	}

	if (current.Role == "user" || current.Role == "staff") && current.UserID != id {
		return helper.Forbidden("You don't have right to do this.")
	}

	var req model.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body must be in valid JSON format")
	}

	empty := IsEmptyPatch(req)
	if empty == true {
		return helper.BadRequest("Body shouldn't be empty")
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	saatIni, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(err, "users")
	}

	saatIni = ApplyPatch(saatIni, req)

	hasil, err := h.repo.Update(ctx, saatIni)
	if err != nil {
		return translateError(err, "user")
	}

	return helper.Success(c, fiber.StatusOK, "user successfuly patch updated", hasil)
}

func (h *UserService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	id, valid := helper.ParamID(c)

	if !valid {
		return helper.BadRequest("Failed to delete user")
	}

	if current.UserID == id {
		return helper.Forbidden(
			"Can't delete yourself")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return translateError(err, "user")
	}
	return helper.NoContent(c)
}

func (s *UserService) AssignRole(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id must be valid")
	}
	var req model.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("JSON body must be valid")
	}
	if errs := authz.ValidateAssignRole(current, id, req, s.perms); len(errs) > 0 {
		return helper.Validation(errs)
	}
	result, err := s.repo.UpdateRole(ctx, id, strings.TrimSpace(req.Role))
	if err != nil {
		return translateError(err, "user")
	}
	return helper.Success(c, fiber.StatusOK, "role user berhasil diubah", result)
}

func translateError(err error, entity string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound(entity + " tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Conflict("username sudah dipakai")
	default:
		return helper.Internal(err)
	}
}

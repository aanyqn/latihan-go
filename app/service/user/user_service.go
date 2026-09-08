package user

import (
	"errors"
	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
	"latihan-fiber/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (h *UserService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := helper.ParseListQuery(c)

	users, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "Error fetching users")
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

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid")
	}

	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "Fail to get user data")
	}

	return helper.Success(c, fiber.StatusOK, "User found!", user)
}

func (h *UserService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format!")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	errs := map[string]string{}
	if req.Username == "" {
		errs["username"] = "must be filled"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "format isn't valid"
	}
	if len(req.Password) < 8 {
		errs["password"] = "minimum in 8 characters"
	}
	if len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	baru, err := h.repo.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: true,
	})
	if err != nil {
		return translateError(c, err, "Failed to save user data")
	}
	return helper.Created(c, "User succesfully Created", baru,
		"/api/v1/users/"+strconv.Itoa(baru.ID))
}

func (h *UserService) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}
	var req model.ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Body must be valid in JSON format")
	}
	errs := map[string]string{}
	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "must be filled"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "must be filled with email format"
	}
	if len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	hasil, err := h.repo.Update(ctx, model.User{
		ID: id, Username: req.Username, Email: req.Email, IsActive: req.IsActive,
	})
	if err != nil {
		return translateError(c, err, "Failed to update user")
	}

	return helper.Success(c, fiber.StatusOK, "User data succesfully replaced", hasil)
}

func (h *UserService) Patch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid!")
	}

	var req model.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Body must be in valid JSON format")
	}
	if req.Username == nil && req.Email == nil && req.IsActive == nil {
		return helper.Fail(c, fiber.StatusBadRequest, "No field is changed")
	}

	saatIni, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "Failed to fetch user")
	}
	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return helper.FailValidation(c, map[string]string{"username": "can't be empty"})
		}
		saatIni.Username = *req.Username
	}
	if req.Email != nil {
		if !strings.Contains(*req.Email, "@") {
			return helper.FailValidation(c, map[string]string{"email": "invalid email format"})
		}
		saatIni.Email = *req.Email
	}
	if req.IsActive != nil {
		saatIni.IsActive = *req.IsActive
	}
	hasil, err := h.repo.Update(ctx, saatIni)
	if err != nil {
		return translateError(c, err, "Failed to update user")
	}
	return helper.Success(c, fiber.StatusOK, "user successfuly patch updated", hasil)
}

func (h *UserService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "Failed to delete user")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return translateError(c, err, "Failed to delete user")
	}
	return helper.NoContent(c)
}

func translateError(c *fiber.Ctx, err error, generalMessage string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.Fail(c, fiber.StatusNotFound, "user tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Fail(c, fiber.StatusConflict, "username sudah dipakai")
	default:
		return helper.Fail(c, fiber.StatusInternalServerError, generalMessage)
	}
}

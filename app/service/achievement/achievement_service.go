package achievement

// import (
// 	"errors"
// 	"latihan-fiber/app/repository"
// 	"latihan-fiber/helper"

// 	"github.com/gofiber/fiber/v2"
// )

// type AchievementService struct {
// 	repo repository.AchievementRepository
// }

// func NewAchievementService(repo repository.AchievementRepository) *AchievementService {
// 	return &AchievementService{repo: repo}
// }

// func (h *AchievementService) Get(c *fiber.Ctx) error {
// 	ctx, cancel := helper.RequestContext(c)
// 	defer cancel()

// 	id, valid := helper.ParamID(c)
// 	if !valid {
// 		return helper.Fail(c, fiber.StatusBadRequest, "ID must be valid")
// 	}

// 	achievement, err := h.repo.FindByID(ctx, id)
// 	if err != nil {
// 		return translateError(c, err, "Fail to get achievement data")
// 	}

// 	return helper.Success(c, fiber.StatusOK, "Achieveent found!", achievement)
// }

// func translateError(c *fiber.Ctx, err error, generalMessage string) error {
// 	switch {
// 	case errors.Is(err, repository.ErrNotFound):
// 		return helper.Fail(c, fiber.StatusNotFound, "Achievement tidak ditemukan")
// 	case errors.Is(err, repository.ErrDuplicate):
// 		return helper.Fail(c, fiber.StatusConflict, "username sudah dipakai")
// 	default:
// 		return helper.Fail(c, fiber.StatusInternalServerError, generalMessage)
// 	}
// }
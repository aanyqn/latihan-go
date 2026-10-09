package auth

import (
	"context"
	"siakad/app/model"
	"siakad/app/repository"
	"siakad/helper"

	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const refreshTokenBytes = 32

type AuthService struct {
	users      repository.UserRepository
	tokens     repository.TokenRepository
	jwt        *helper.JWTManager
	refreshTTL time.Duration
	perms      *helper.PermissionSet
	students   repository.StudentRepository
}

func NewAuthService(
	users repository.UserRepository,
	tokens repository.TokenRepository,
	jwtManager *helper.JWTManager,
	refreshTTL time.Duration,
	perms *helper.PermissionSet,
	students repository.StudentRepository,
) *AuthService {
	return &AuthService{
		users: users, tokens: tokens, jwt: jwtManager, refreshTTL: refreshTTL, perms: perms, students: students,
	}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if errs := ValidateLogin(req); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}
	user, err := s.users.FindByEmail(ctx, strings.TrimSpace(req.Email))
	if err != nil {
		helper.VerifyDummyPassword(req.Password)
		return helper.Unauthorized("email or password is invalid")
	}
	if !helper.VerifyPassword(user.Password, req.Password) {
		return helper.Unauthorized("email or password is invalid")
	}
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.Internal(err, "Token pair failed")
	}
	return helper.Success(c, fiber.StatusOK, "login success", pair)
}

func (s *AuthService) Refresh(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body must be valid JSON")
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		return helper.BadRequest("refresh token should be filled")
	}
	hash := helper.SHA256Hex(req.RefreshToken)
	stored, err := s.tokens.FindActive(ctx, hash)
	if err != nil {
		return helper.Unauthorized("refresh token isn't valid or expired")
	}
	user, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil {
		return helper.Unauthorized("Can't access account")
	}

	if err := s.tokens.Revoke(ctx, hash); err != nil {
		return helper.Internal(err, "Revoking error")
	}
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.Internal(err, "Token pair failed")
	}
	return helper.Success(c, fiber.StatusOK, "token successfully updated", pair)
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body must be valid JSON")
	}
	if strings.TrimSpace(req.RefreshToken) != "" {
		_ = s.tokens.Revoke(ctx, helper.SHA256Hex(req.RefreshToken))
	}
	return helper.Success(c, fiber.StatusOK, "logout success", nil)
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()
	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("Not authenticated")
	}
	user, err := s.users.FindByID(ctx, authUser.UserID)
	if err != nil {
		return helper.NotFound("user not found")
	}
	if user.Role == "mahasiswa" {
		student, err := s.students.FindByUserID(ctx, authUser.UserID)
		if err != nil {
			return helper.NotFound("student not found")
		}
		return helper.Success(c, fiber.StatusOK, "successfully got profile information", fiber.Map{
			"user":        user,
			"student":     student,
			"permissions": s.perms.PermissionsOf(strings.ToLower(strings.TrimSpace(user.Role))),
		})
	}
	return helper.Success(c, fiber.StatusOK, "successfully got profile information", fiber.Map{
		"user":        user,
		"permissions": s.perms.PermissionsOf(strings.ToLower(strings.TrimSpace(user.Role))),
	})
}

func (s *AuthService) issueTokenPair(
	ctx context.Context, user model.User,
) (model.TokenPair, error) {
	accessToken, err := s.jwt.GenerateAccess(user)
	if err != nil {
		return model.TokenPair{}, err
	}
	refreshToken, err := helper.RandomToken(refreshTokenBytes)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = s.tokens.Save(ctx, model.RefreshToken{
		UserID:    user.ID,
		TokenHash: helper.SHA256Hex(refreshToken),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	})

	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwt.AccessTTL().Seconds()),
	}, nil

}

package handler

import (
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/request"
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/helper"
	"github.com/fatihrizqon/symetra-service/internal/service"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	IAuthService service.IAuthService
	Cookie       helper.RefreshCookie
}

// NewAuthHandler — FIXED
//
// BUG: Cookie field was left as zero-value (TTL=0, Production=false).
// This caused the Set() helper to use MaxAge=0, which in HTTP means
// "session cookie" — it expires when the browser closes. Since the
// DB credential ExpiresAt is 15 minutes, the cookie should match.
//
// FIX: Initialize Cookie with a proper TTL (7 days to match refresh JWT duration)
// and Production flag from environment. In production, set Production=true via config.
//
// NOTE: The TTL mismatch between DB credential (15 min) and JWT refresh token (7 days)
// was intentional in the original design — the DB credential gets rotated on every
// refresh call. The cookie should live as long as the JWT (7 days).
func NewAuthHandler(serv service.IAuthService) *AuthHandler {
	return &AuthHandler{
		IAuthService: serv,
		Cookie: helper.RefreshCookie{
			TTL:        7 * 24 * time.Hour, // Match refresh JWT expiry
			Production: false,              // Set to true in production (enables Secure flag)
		},
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return a JWT token in a cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} response.AuthJSON
// @Failure 400 {object} response.JSON "Invalid request format"
// @Failure 401 {object} response.JSON "Authentication failed"
// @Router /api/v1/auth/login [post]
func (handler *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req request.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := handler.IAuthService.Login(req)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	handler.Cookie.Set(ctx, "refresh_token", result.RefreshToken)

	return ctx.Status(fiber.StatusOK).JSON(response.AuthJSON{
		Message: "you are authenticated",
		Status:  fiber.StatusOK,
		User: response.UserInfo{
			Id:              result.User.Id,
			Username:        result.User.Username,
			Name:            result.User.Name,
			Email:           result.User.Email,
			Status:          result.User.Status,
			IsSuperadmin:    result.User.IsSuperadmin,
			EmailVerifiedAt: result.User.EmailVerifiedAt.Format(time.RFC3339),
		},
		AccessToken: result.AccessToken,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generate new access token and rotate refresh token
// @Tags Auth
// @Produce json
// @Success 200 {object} response.AuthJSON
// @Failure 401 {object} response.JSON "Unauthorized"
// @Router /api/v1/auth/refresh [post]
//
// FIXED: The original handler called Cookie.Clear() BEFORE validating the token.
// If RefreshToken() failed, the cookie was already cleared but no new one was set,
// permanently logging out the user even on transient errors (e.g. DB timeout).
//
// FIX: Only clear the old cookie AFTER successfully generating new tokens.
func (handler *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: "missing refresh token",
		})
	}

	// Validate & rotate token FIRST — do not touch cookie until we have a new token
	result, err := handler.IAuthService.RefreshToken(refreshToken)
	if err != nil {
		// Do NOT clear the cookie here — the token may still be valid
		// (e.g. DB was temporarily unavailable). Let the client retry.
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	// Only clear old cookie + set new one after successful rotation
	handler.Cookie.Clear(ctx, "refresh_token")
	handler.Cookie.Set(ctx, "refresh_token", result.RefreshToken)

	return ctx.Status(fiber.StatusOK).JSON(response.AuthJSON{
		Status:  fiber.StatusOK,
		Message: "token refreshed",
		User: response.UserInfo{
			Id:              result.User.Id,
			Username:        result.User.Username,
			Name:            result.User.Name,
			Email:           result.User.Email,
			Status:          result.User.Status,
			IsSuperadmin:    result.User.IsSuperadmin,
			EmailVerifiedAt: result.User.EmailVerifiedAt.Format(time.RFC3339),
		},
		AccessToken: result.AccessToken,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user, blacklist refresh token, clear cookie
// @Tags Auth
// @Success 200 {object} response.JSON "Successfully logged out"
// @Failure 401 {object} response.JSON "Unauthorized"
// @Router /api/v1/auth/logout [post]
func (handler *AuthHandler) Logout(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: "unauthorized",
		})
	}

	if err := handler.IAuthService.Logout(refreshToken); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	handler.Cookie.Clear(ctx, "refresh_token")

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "successfully logged out",
	})
}

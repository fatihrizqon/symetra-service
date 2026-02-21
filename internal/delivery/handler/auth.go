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

func NewAuthHandler(serv service.IAuthService) *AuthHandler {
	return &AuthHandler{IAuthService: serv}
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
func (handler *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: "missing refresh token",
		})
	}

	handler.Cookie.Clear(ctx, "refresh_token")

	result, err := handler.IAuthService.RefreshToken(refreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

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
			EmailVerifiedAt: result.User.EmailVerifiedAt.Format(time.RFC3339),
		},
		AccessToken: result.AccessToken,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user dengan menghapus access_token dan refresh_token dari cookie,
// @Description serta memasukkan refresh token ke blacklist.
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

/*

// Refresh Token godoc
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token from HttpOnly cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.JSON "Access token refreshed"
// @Failure 400 {object} response.JSON "Invalid request format"
// @Failure 401 {object} response.JSON "Invalid or missing refresh token"
// @Router /api/v1/auth/refresh [post]
func (handler *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("refresh token required"))
		return nil
	}

	claims, err := util.ParseToken(refreshToken)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("invalid refresh token"))
		return nil
	}

	refreshToken, _ = util.CreateToken(entity.User{
		Id: uuid.MustParse(claims.UserID),
	})

	// setAuthCookies(ctx, accessToken, refreshToken)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Access token refreshed",
	})
}


// Get User Info godoc
// @Summary Get authenticated user info
// @Description Retrieve the current authenticated user's information using the access token stored in HttpOnly cookie or Authorization header.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.AuthJSON "User info retrieved"
// @Failure 401 {object} response.JSON "Invalid or missing access token"
// @Router /api/v1/auth/me [get]
func (handler *AuthHandler) Me(ctx *fiber.Ctx) error {
	accessToken := ctx.Cookies("access_token")

	if accessToken == "" {
		authHeader := ctx.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if accessToken == "" {
		return errorResponse(ctx, fiber.StatusUnauthorized, "access token required")
	}

	claims, err := util.ParseToken(accessToken)
	if err != nil {
		return errorResponse(ctx, fiber.StatusUnauthorized, "invalid access token")
	}

	userID, err := parseUserIDFromClaims(claims)
	if err != nil {
		return errorResponse(ctx, fiber.StatusUnauthorized, "invalid user ID")
	}

	username, _ := claims["username"].(string)
	name, _ := claims["name"].(string)
	email, _ := claims["email"].(string)

	statusFloat, _ := claims["status"].(float64)
	status := int(statusFloat)

	return ctx.Status(fiber.StatusOK).JSON(response.AuthJSON{
		Message: "user info retrieved",
		Status:  fiber.StatusOK,
		User: response.UserInfo{
			Id:       userID,
			Username: username,
			Name:     name,
			Email:    email,
			Status:   status,
		},
	})
}
*/

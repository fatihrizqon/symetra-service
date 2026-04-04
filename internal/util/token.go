package util

import (
	"time"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var accessSecret []byte
var refreshSecret []byte

// Claims is the access token payload.
// Note: company_id is NOT embedded in the token — it is read from
// the X-Company-ID request header and validated per-request by the
// company middleware. This keeps tokens stateless across company switches.
type Claims struct {
	UserID    uuid.UUID `json:"uid"`
	SessionID uuid.UUID `json:"sid"`
	jwt.RegisteredClaims
}

func NewJWT(config *viper.Viper) {
	jwtSecret := config.GetString("jwt.secret")
	jwtRefreshSecret := config.GetString("jwt.refresh_secret")
	if jwtSecret == "" || jwtRefreshSecret == "" {
		panic("jwt secret missing")
	}
	accessSecret = []byte(jwtSecret)
	refreshSecret = []byte(jwtRefreshSecret)
}

func CreateAccessToken(user entity.User, sessionID uuid.UUID) (string, error) {
	return create(user, sessionID, accessSecret, 15*time.Minute)
}

func CreateRefreshToken(user entity.User, sessionID uuid.UUID) (string, error) {
	return create(user, sessionID, refreshSecret, 7*24*time.Hour)
}

func ParseAccessToken(token string) (*Claims, error) {
	return parse(token, accessSecret)
}

func ParseRefreshToken(token string) (*Claims, error) {
	return parse(token, refreshSecret)
}

func parse(tokenString string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fiber.ErrUnauthorized
			}
			return secret, nil
		},
	)
	if err != nil || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}
	return claims, nil
}

func create(user entity.User, sessionID uuid.UUID, secret []byte, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:    user.Id,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ─── Context helpers ──────────────────────────────────────────────────────────

// GetCallerID extracts the authenticated user ID from Fiber locals.
func GetCallerID(ctx *fiber.Ctx) (uuid.UUID, error) {
	claims, ok := ctx.Locals("auth").(*Claims)
	if !ok || claims == nil {
		return uuid.Nil, fiber.ErrUnauthorized
	}
	return claims.UserID, nil
}

// GetCompanyID extracts the validated company ID from Fiber locals.
// Set by the company middleware after verifying X-Company-ID header membership.
func GetCompanyID(ctx *fiber.Ctx) (uuid.UUID, error) {
	cid, ok := ctx.Locals("company_id").(uuid.UUID)
	if !ok || cid == uuid.Nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "X-Company-ID header is required")
	}
	return cid, nil
}

// GetCallerRole extracts the caller's role within the current company from locals.
// Set by the company middleware.
func GetCallerRole(ctx *fiber.Ctx) entity.CompanyRole {
	role, _ := ctx.Locals("company_role").(entity.CompanyRole)
	return role
}

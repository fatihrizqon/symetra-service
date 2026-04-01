package helper

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type RefreshCookie struct {
	Name       string
	TTL        time.Duration
	Production bool
}

func (c RefreshCookie) Set(ctx *fiber.Ctx, name string, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Production,
		// SameSite=Lax (bukan Strict) agar cookie terkirim pada cross-origin request
		// dari frontend (localhost:5173) ke backend (localhost:3000).
		// SameSite=Strict memblokir cookie pada request cross-origin meskipun
		// masih same-site (localhost), sehingga /api/v1/auth/refresh selalu 401.
		// Di production (Secure=true + same domain), Lax sudah cukup aman.
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   int(c.TTL.Seconds()),
	})
}

func (c RefreshCookie) Clear(ctx *fiber.Ctx, name string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})
}

package middleware

import (
	"fmt"

	"github.com/fatihrizqon/symetra-service/internal/entity"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/fatihrizqon/symetra-service/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// NewCompany returns a middleware that:
//  1. Reads the X-Company-ID header
//  2. Validates that the authenticated user is a member of that company
//  3. Injects company_id and company_role into Fiber locals
//  4. Allows superadmins to access any company without membership check
//
// This middleware must run AFTER NewAuth (which sets ctx.Locals("auth")).
func NewCompany(memberRepo repository.ICompanyMemberRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 1. Get authenticated user from auth middleware
		callerID, err := util.GetCallerID(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  401,
				"message": "unauthorized",
			})
		}

		// 2. Read X-Company-ID header
		companyIDStr := ctx.Get("X-Company-ID")
		if companyIDStr == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  400,
				"message": "X-Company-ID header is required",
			})
		}

		companyID, err := uuid.Parse(companyIDStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  400,
				"message": "X-Company-ID is not a valid UUID",
			})
		}

		// 3. Check membership (superadmin bypass)
		member, err := memberRepo.FindMembership(callerID, companyID)
		if err != nil {
			// Not a member — check if superadmin via platform role
			// Superadmin is stored on the User record directly (manual DB set)
			isSuperadmin, _ := memberRepo.IsSuperadmin(callerID)
			if !isSuperadmin {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"status":  403,
					"message": fmt.Sprintf("you do not have access to company %s", companyIDStr),
				})
			}
			// Superadmin gets owner-level permissions for any company
			ctx.Locals("company_id", companyID)
			ctx.Locals("company_role", entity.RoleSuperadmin)
			return ctx.Next()
		}

		// 4. Inject into locals for downstream handlers
		ctx.Locals("company_id", companyID)
		ctx.Locals("company_role", member.Role)

		return ctx.Next()
	}
}

// NewRequirePermission returns a middleware that gates a route behind a specific permission.
// Must run after NewCompany (which sets company_role).
//
// Usage in route setup:
//   rc.App.Delete("/api/v1/companies/:id", rc.CompanyMiddleware, middleware.NewRequirePermission("company:delete"), rc.CompanyHandler.Delete)
func NewRequirePermission(permission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role := util.GetCallerRole(ctx)
		if !entity.HasPermission(role, permission) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  403,
				"message": fmt.Sprintf("permission denied: %s required", permission),
			})
		}
		return ctx.Next()
	}
}

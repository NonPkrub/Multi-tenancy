package middleware

import (
	"go-multi-tenancy/internals/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		// Check if the Authorization header starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
		}

		// Extract the token
		tokenString := parts[1]

		// Validate and parse the JWT token
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Store user information in locals
		c.Locals("user_id", claims.UserID)
		c.Locals("company_id", claims.CompanyID)
		c.Locals("branch_id", claims.BranchID)

		return c.Next()
	}
}

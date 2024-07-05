package middleware

import (
	"go-multi-tenancy/internals/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTAuth is a middleware function that handles JWT authentication.
// It extracts the JWT token from the Authorization header, validates it,
// and stores user information in the locals of the context.
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// Check if the Authorization header is missing
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		// Extract the token from the Authorization header
		tokenString := parts[1]

		// Parse and validate the JWT token
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store user information in the locals of the context
		c.Locals("username", claims.ID)
		c.Locals("company", claims.Company)
		c.Locals("branch", claims.Branch)
		c.Locals("role", claims.Role)

		// Continue to the next handler in the chain
		return c.Next()
	}
}

// AuthorizeRole is a middleware function that checks if the user's role
// matches the required role. If the user's role does not match the required
// role, it returns a 403 Forbidden response.
// Parameters:
// - requiredRole: the required role for access.
// Returns:
// - fiber.Handler: a function that handles the request and returns an error.
func AuthorizeRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user's role from the context
		userRole, ok := c.Locals("role").(string)

		// Check if the user's role is missing or does not match the required role
		if !ok || userRole != requiredRole {
			// Return a 403 Forbidden response
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
		}

		// Continue to the next handler in the chain
		return c.Next()
	}
}

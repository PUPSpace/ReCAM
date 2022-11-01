package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func ValidateAuthToken(c *fiber.Ctx) map[string]interface{} {
	// Get now time.
	now := time.Now().Unix()
	returnMap := make(map[string]interface{})
	returnMap["error"] = false

	// Get claims from JWT.
	claims, err := ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		returnMap["status"] = fiber.StatusInternalServerError
		returnMap["error"] = true
		returnMap["msg"] = err.Error()
		return returnMap
	}
	// Set expiration time from JWT data of current route.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 	"error": true,
		// 	"msg":   "unauthorized, check expiration time of your token",
		// })
		returnMap["status"] = fiber.StatusUnauthorized
		returnMap["error"] = true
		returnMap["msg"] = "unauthorized, check expiration time of your token"
		return returnMap
	}

	return returnMap
}

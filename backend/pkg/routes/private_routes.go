package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kaleemubarok/recam/backend/app/controllers"
	"github.com/kaleemubarok/recam/backend/pkg/middleware"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	// route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)           // create a new book
	route.Post("/route", middleware.JWTProtected(), controllers.CreateRoute)         // create a new route
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)   // renew Access & Refresh tokens

	// Routes for PUT method:
	route.Put("/route", middleware.JWTProtected(), controllers.UpdateRoute) // update one route by ID

	// Routes for DELETE method:
	// route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID
}

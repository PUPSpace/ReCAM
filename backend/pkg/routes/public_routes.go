package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kaleemubarok/recam/backend/app/controllers"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	// route.Get("/books", controllers.GetBooks)   // get list of all books
	// route.Get("/book/:id", controllers.GetBook) // get one book by ID
	route.Get("/routes", controllers.GetRoutes)          // get list of all  routes
	route.Get("/route/chart", controllers.GetRouteChart) // get route logs chart for dashboard
	route.Get("/route/:id", controllers.GetRoute)        // get one book by ID

	// Routes for POST method:
	route.Post("/user/sign/up", controllers.UserSignUp) // register a new user
	route.Post("/user/sign/in", controllers.UserSignIn) // auth, return Access & Refresh tokens

	// PublicRoutes func for describe group of ReCAM public routes.
	route.All("/go/:token/:slug/*", controllers.RecamControl) // handle get request

	// https://recam.djbk.org/go/7MDZWAYNC/google?q=kenapa+ayam+berkokok+jam+12+malam
}

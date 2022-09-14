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
	route.Post("/route", controllers.CreateRoute)                                    // create a new route
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)   // renew Access & Refresh tokens

	// Routes for PUT method:
	route.Put("/route", controllers.UpdateRoute)                     // update one route by ID
	route.Put("/route/:id/renewtoken", controllers.UpdateRouteToken) // update route token by ID

	// Routes for DELETE method:
	// route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID

	route.Get("/logs", controllers.GetLogs)                       /*TODO: add jwt middleware*/
	route.Get("/logs/:err", controllers.GetLogsSpecial)           /*TODO: add jwt middleware*/
	route.Get("/logs/count-u5xx", controllers.CountUnresolved5XX) /*TODO: add jwt middleware*/
	route.Get("/log/:id", controllers.GetLog)                     // get one log details by ID

}

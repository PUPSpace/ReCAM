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
	route.Post("/route", middleware.JWTProtected(), controllers.CreateRoute)         // create a new route
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)   // renew Access & Refresh tokens

	// Routes for PUT method:
	route.Put("/route", middleware.JWTProtected(), controllers.UpdateRoute)                     // update one route by ID
	route.Put("/route/:id/renewtoken", middleware.JWTProtected(), controllers.UpdateRouteToken) // update route token by ID
	route.Put("/log/:id/resolved", middleware.JWTProtected(), controllers.UpdateResolvedStatus) // resolve a log details by ID

	// Routes for DELETE method:
	route.Delete("/route/:id", middleware.JWTProtected(), controllers.DeleteRoute) // delete one route by ID

	route.Get("/logs", middleware.JWTProtected(), controllers.GetLogs)                       /*TODO: add jwt middleware*/
	route.Get("/logs/:err", middleware.JWTProtected(), controllers.GetLogsSpecial)           /*TODO: add jwt middleware*/
	route.Get("/logs/count-u5xx", middleware.JWTProtected(), controllers.CountUnresolved5XX) /*TODO: add jwt middleware*/
	route.Get("/log/:id", middleware.JWTProtected(), controllers.GetLog)                     // get one log details by ID
	route.Get("/route/:id", middleware.JWTProtected(), controllers.GetRoute)                 // get one route by ID
	route.Get("/routes", middleware.JWTProtected(), controllers.GetRoutes)                   // get list of all  routes

}

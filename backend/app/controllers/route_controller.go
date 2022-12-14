package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/kaleemubarok/recam/backend/app/models"
	"github.com/kaleemubarok/recam/backend/pkg/utils"
	"github.com/kaleemubarok/recam/backend/platform/cache"
	"github.com/kaleemubarok/recam/backend/platform/database"
)

// GetRoutes func gets all exists routes.
// @Description Get all exists routes.
// @Summary get all exists routes
// @Tags Routes
// @Accept json
// @Produce json
// @Success 200 {array} models.Route
// @Router /v1/routes [get]
func GetRoutes(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all routes.
	routes, err := db.GetRoutes()
	if err != nil {
		// Return, if routes not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  true,
			"msg":    "routes were not found " + err.Error(), /*TODO: delete this debug purpose error print snippet*/
			"count":  0,
			"routes": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":  false,
		"msg":    nil,
		"count":  len(routes),
		"routes": routes,
	})
}

// GetRoute func gets route by given ID or 404 error.
// @Description Get route by given ID.
// @Summary get route by given ID
// @Tags Route
// @Accept json
// @Produce json
// @Param id path string true "Route ID"
// @Success 200 {object} models.Route
// @Router /v1/route/{id} [get]
func GetRoute(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error() + err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get route by ID.
	route, err := db.GetRoute(id)
	if err != nil {
		// Return, if route not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "route with the given ID is not found" + err.Error(),
			"route": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"route": route,
	})
}

// CreateRoute func for creates a new route.
// @Description Create a new route.
// @Summary create a new route
// @Tags Route
// @Accept json
// @Produce json
// @Param title body string true "Title" /*TODO: Benerin nih Paramnya*/
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param route_attrs body models.RouteAttrs true "Route attributes"
// @Success 200 {object} models.Route
// @Security ApiKeyAuth
// @Router /v1/route [post]
func CreateRoute(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}

	// Create new Route struct
	route := &models.Route{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(route); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set initialized default data for route:
	route.ID = uuid.New()
	route.CreatedAt = time.Now()
	route.UserID = uuid.MustParse("70872edf-974f-4c75-8f23-c0929268e041") //TODO: Update with dynamic ID

	if route.IsRetryable == "" {
		route.IsRetryable = "N"
		route.MaxRetry = 0
		route.RetryPeriod = 0
	}

	// Set slug from route name
	generatedSlug, err := utils.GenerateSlug(route.Name)
	if err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// check duplicate slug
	if slugFounded, err := db.CountSlug(generatedSlug); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	} else {
		if slugFounded > 0 {
			str := strconv.Itoa(slugFounded)
			route.Slug = generatedSlug + "-" + str
		} else {
			route.Slug = generatedSlug
		}
	}

	// Generate token
	token, err := utils.GenerateToken()
	if err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	route.Token = token

	// Create route by given model.
	if err := db.CreateRoute(route); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"route": route,
	})
}

// UpdateRoute func for updates route by given ID.
// @Description Update route.
// @Summary update route
// @Tags Route
// @Accept json
// @Produce json
// @Param id body string true "Book ID" /*TODO: Benerin nih Paramnya*/
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param route_status body integer true "Book status"
// @Param route_attrs body models.BookAttrs true "Book attributes"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/route [put]
func UpdateRoute(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}
	/* if !credential {
	// 	// Return status 403 and permission denied error message.
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   "permission denied, check credentials of your token",
	// 	})
	// }
	*/

	// Create new Route struct
	route := &models.Route{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(route); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if route with given ID is exists.
	foundedRoute, err := db.GetRoute(route.ID)
	if err != nil {
		// Return status 404 and route not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "route with this ID not found",
		})
	}

	// Update route by given ID.
	if err := db.UpdateRoute(foundedRoute.ID, route); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// delete existing redis key when route updated
	_ = connRedis.Del(context.Background(), foundedRoute.Slug).Err()

	// Return status 201.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}

/*
// DeleteBook func for deletes route by given ID.
// @Description Delete route by given ID.
// @Summary delete route by given ID
// @Tags Book
// @Accept json
// @Produce json
// @Param id body string true "Book ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/route [delete]
func DeleteBook(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current route.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `route:delete` from JWT data of current route.
	credential := claims.Credentials[repository.BookDeleteCredential]

	// Only route creator with `route:delete` credential can delete his route.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new Book struct
	route := &models.Book{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(route); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Book model.
	validate := utils.NewValidator()

	// Validate route fields.
	if err := validate.StructPartial(route, "id"); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if route with given ID is exists.
	foundedBook, err := db.GetRoute(route.ID)
	if err != nil {
		// Return status 404 and route not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "route with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his route.
	if foundedBook.UserID == userID {
		// Delete route by given ID.
		if err := db.DeleteBook(foundedBook.ID); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, only the creator can delete his route",
		})
	}
}
*/

// UpdateRoute func for updates route by given ID.
// @Description Update route.
// @Summary update route
// @Tags Route
// @Accept json
// @Produce json
// @Param id body string true "Book ID" /*TODO: Benerin nih Paramnya*/
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param route_status body integer true "Book status"
// @Param route_attrs body models.BookAttrs true "Book attributes"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/route [put]
func UpdateRouteToken(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}

	// Create new Route struct
	route := &models.Route{}

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error() + err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if route with given ID is exists.
	var dbRoute models.Route
	dbRoute, err = db.GetRoute(id)
	if err != nil {
		// Return status 404 and route not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "route with this ID not found",
		})
	}

	// Generate token
	token, err := utils.GenerateToken()
	if err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	route.Token = token
	route.ID = id

	// Update route by given ID.
	if err := db.UpdateRouteToken(id, token); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// get route from Redis.
	var routeData models.Route
	_, err = connRedis.Get(context.Background(), dbRoute.Slug).Result()

	if err != redis.Nil {
		// delete existing key
		_ = connRedis.Del(context.Background(), dbRoute.Slug)

		// get route from DB
		routeData, err = db.GetRouteSlug(dbRoute.Slug)
		if err != nil {
			// Return status 400 and error message.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   fmt.Sprintf("no route matching with your requested %s, please try again", dbRoute.Slug),
			})
		}
		routeDataJson, err := json.Marshal(routeData)
		if err != nil {
			return errors.New("there is an error on marshaling json data from redis")
		}

		_ = connRedis.Set(context.Background(), routeData.Slug, string(routeDataJson), 60*60*time.Second).Err()

	}

	// Return status 201.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":     false,
		"msg":       "updated with new token " + token,
		"new_token": token,
	})
}

func GetRouteChart(c *fiber.Ctx) error {
	var chartData []models.Chart

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// get chart cache from Redis.
	_, err = connRedis.Get(context.Background(), "dashboardChart").Result()

	if err == redis.Nil {
		log.Println("processing chart from DB")
		// Create database connection.
		db, err := database.OpenDBConnection()
		if err != nil {
			// Return status 500 and database connection error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Get cache from DB.
		chartData, err = db.GetRouteChart()
		if err != nil {
			// Return status 400 and error message.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "error on getting chart data from db," + err.Error(),
				"chart": nil,
			})
		}
		rchartDataJson, err := json.Marshal(chartData)
		if err != nil {
			return errors.New("there is an error on marshaling chart data to redis")
		}

		err = connRedis.Set(context.Background(), "dashboardChart", string(rchartDataJson), 60*time.Second).Err()
		if err != nil {
			return errors.New(err.Error())
		}

	} else {
		// reading chart from redis
		chartRedis, _ := connRedis.Get(context.Background(), "dashboardChart").Result()
		err = json.Unmarshal([]byte(chartRedis), &chartData)
		if err != nil {
			return errors.New("there is an error on unmarshaling json data from redis")
		}
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"chart": chartData,
	})
}

// DeleteRoute func for deletes route by given ID.
// @Description Delete route by given ID.
// @Summary delete route by given ID
// @Tags Route
// @Accept json
// @Produce json
// @Param id body string true "Route ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/route/:id [delete]
func DeleteRoute(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error() + err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get route by ID.
	route, err := db.GetRoute(id)
	if err != nil {
		// Return, if route not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "route with the given ID is not found" + err.Error(),
			"route": nil,
		})
	}

	// Delete route by given ID.
	if err := db.DeleteRoute(route.ID); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}

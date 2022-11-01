package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/kaleemubarok/recam/backend/app/models"
	"github.com/kaleemubarok/recam/backend/pkg/utils"
	"github.com/kaleemubarok/recam/backend/platform/database"
)

// GetLogs func gets all exists routes.
// @Description Get all exists routes.
// @Summary get all exists routes
// @Tags Routes
// @Accept json
// @Produce json
// @Success 200 {array} models.Route
// @Router /v1/routes [get]
func GetLogs(c *fiber.Ctx) error {
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
	logs, err := db.GetLogs()
	if err != nil {
		// Return, if routes not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "log were not found " + err.Error(), /*TODO: delete this debug purpose error print snippet*/
			"count": 0,
			"logs":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(logs),
		"logs":  logs,
	})
}

// GetRoute func gets log by given ID or 404 error.
// @Description Get route by given ID.
// @Summary get route by given ID
// @Tags Route
// @Accept json
// @Produce json
// @Param id path string true "Route ID"
// @Success 200 {object} models.Route
// @Router /v1/route/{id} [get]
func GetLog(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}

	// Catch log ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	// fmt.Println(id)
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

	// Get log by ID.
	rlog, err := db.GetLog(id)
	if err != nil {
		// Return, if log not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "log with the given ID is not found" + err.Error(),
			"log":   nil,
		})
	}

	//decrypt start
	key := os.Getenv("ENCRYPT_KEY")
	decText, err := utils.Decrypt(rlog.Data, key)
	if err != nil {
		fmt.Println("error decrypting your encrypted log data: ", err)
	}

	decOthers, err := utils.Decrypt(rlog.Others, key)
	if err != nil {
		fmt.Println("error decrypting your encrypted log data others: ", err)
	}
	// fmt.Println("DECRYPTED---> " + decText)

	others := models.RetryLog{}
	json.Unmarshal([]byte(decOthers), &others)

	out, _ := utils.ExtractLogData(decText)
	// print log
	// log.Println(out["query"].(string))
	// log.Println(out["rqHeader"].(string))
	// log.Println(out["rqBody"].(string))
	// log.Println(out["rsBody"].(string))

	data := fiber.Map{
		"req_header":      out["rqHeader"].(string),
		"req_query":       out["query"].(string),
		"req_body":        out["rqBody"].(string),
		"res_body":        out["rsBody"].(string),
		"host_addr":       rlog.HostAddr,
		"name":            rlog.Name,
		"code":            rlog.ResponseCode,
		"connection_type": rlog.Type,
		"created_at":      rlog.CreatedAt,
		"slug":            rlog.Slug,
		"trial_attempt":   rlog.TrialAttempt,
		"others":          others,
		"is_resolved":     rlog.IsResolved,
	}
	//decrypt end

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"log":   data,
	})
}

func GetLogsSpecial(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}

	// Catch err cd from URL.
	ec := c.Params("err")

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
	logs, err := db.GetLogsSpecial(ec)
	if err != nil {
		// Return, if routes not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "log were not found " + err.Error(), /*TODO: delete this debug purpose error print snippet*/
			"count": 0,
			"logs":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(logs),
		"logs":  logs,
	})
}

func CountUnresolved5XX(c *fiber.Ctx) error {
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
	count, err := db.CountUnresolved5XX()
	if err != nil {
		// Return, if routes not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(), /*TODO: delete this debug purpose error print snippet*/
			"count": 0,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": count,
	})
}

// UpdateLogSolved func for updates communication details by given ID.
// @Description Update route.
// @Summary update route
// @Tags Route
// @Accept json
// @Produce json
// @Param id body string true "Book ID" /*TODO: Benerin nih Paramnya*/
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/route [put]
func UpdateResolvedStatus(c *fiber.Ctx) error {
	authValidationResult := utils.ValidateAuthToken(c)
	if authValidationResult["error"] == true {
		return c.Status(authValidationResult["status"].(int)).JSON(fiber.Map{
			"error": true,
			"msg":   authValidationResult["msg"],
		})
	}

	// Catch log ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	// fmt.Println(id)
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

	// Get log by ID.
	_, err = db.GetLog(id)
	if err != nil {
		// Return, if log not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "log with the given ID is not found" + err.Error(),
			"log":   nil,
		})
	}

	// Update route by given ID.
	if err := db.UpdateResolvedLog(id); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 201.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "updated " + id.String() + " as resolved",
	})
}

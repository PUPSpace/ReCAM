package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flowchartsman/retry"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kaleemubarok/recam/backend/app/models"
	"github.com/kaleemubarok/recam/backend/pkg/utils"
	"github.com/kaleemubarok/recam/backend/platform/database"
	"github.com/valyala/fasthttp"
)

func RecamControl(c *fiber.Ctx) error {
	/*TODO: Check redis get all route config before procceed*/

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create new Route struct for logging purpose
	routeLog := &models.RouteLog{}

	log.Println("AllParams")
	log.Println(c.AllParams())

	// get route from DB
	allParams := c.AllParams()
	routeData, err := db.GetRouteSlug(allParams["slug"])
	if err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//kick out if method different from db
	if !strings.EqualFold(routeData.Token, allParams["token"]) {
		// Return status 404 and error message.
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "your token is not authorized, please contact admin",
		})
	}
	//kick out if method different from db
	if !strings.EqualFold(routeData.CommType, string(c.Request().Header.Method())) {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("error while matching request method %s, expected %s", string(c.Request().Header.Method()), routeData.CommType),
		})
	}

	//update context
	rqData := reqData{
		reqUri:         routeData.HostAddr,
		isRetryable:    routeData.IsRetryable,
		maxRetry:       routeData.MaxRetry,
		retryPeriod:    routeData.RetryPeriod,
		retryMaxPeriod: 9999,
		reqURI:         routeData.HostAddr,
		reqQuery:       string(c.Request().URI().QueryString()),
	}

	err = RetryFastHttp(c, rqData)
	if err != nil {
		return err
	}

	code, body := c.Response().StatusCode(), c.Response().Body()

	// try := 3
	// period := 10

	//Declare map to store request response data
	rqrs := jwt.MapClaims{}
	rqrs["query"] = string(c.Request().URI().QueryString())
	rqrs["rqHeader"] = string(c.Request().Header.Header())
	rqrs["rqBody"] = string(c.Request().Body())
	rqrs["rsBody"] = string(body)

	// Set correct content-type
	contentType := string(c.Response().Header.ContentType())
	c.Set(fiber.HeaderContentType, contentType)

	log.Println("contentType: " + contentType)

	// Set RouteLog data
	routeLog.ID = uuid.New()

	// // Check, if received JSON data is valid.
	// if err := c.BodyParser(route); err != nil {
	// 	// Return status 400 and error message.
	// 	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   err.Error(),
	// 	})
	// }

	// Set user ID from JWT data of current user.
	// userID := claims.UserID

	//encrypt part start
	encrypted, err := utils.GenerateReqResLog(rqrs)
	if err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	key := os.Getenv("ENCRYPT_KEY")
	encrypted, err = utils.Encrypt(encrypted, key)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
	}
	//encrypt part end

	// set logs data
	routeLog.Data = string(encrypted)
	routeLog.ResponseCode = strconv.Itoa(code)
	routeLog.RouteID = routeData.ID /*TODO: cek ini*/
	routeLog.Type = string(c.Request().Header.Method())
	routeLog.CreatedAt = time.Now()

	// // Create a new validator for a Route model.
	// validate := utils.NewValidator()

	// // Validate route fields.
	// if err := validate.Struct(routeLog); err != nil {
	// 	// Return, if some fields are not valid.
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   utils.ValidatorErrors(err),
	// 	})
	// }

	// log.Println("routeLog before insertion: ")
	// log.Println(routeLog)

	// create a communication log
	if err := db.CreateLog(routeLog); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// //decrypt start
	// decText, err := utils.Decrypt(encrypted, key)
	// if err != nil {
	// 	fmt.Println("error decrypting your encrypted text: ", err)
	// }
	// // fmt.Println("DECRYPTED---> " + decText)

	// out, _ := utils.ExtractLogData(decText)
	// // print log
	// log.Println(out["query"].(string))
	// log.Println(out["rqHeader"].(string))
	// log.Println(out["rqBody"].(string))
	// log.Println(out["rsBody"].(string))
	// //decrypt end

	// Send response body
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
	c.Response().Header.SetStatusCode(code)
	return c.Send(body)

	// return nil
}

func RetryFastHttp(c *fiber.Ctx, rqData reqData) error {
	retrier := retry.NewRetrier(1, 100*time.Millisecond, time.Second)
	if strings.EqualFold("Y", c.Get("isRetryable")) {
		retrier = retry.NewRetrier(rqData.maxRetry, time.Duration(rqData.retryPeriod)*time.Millisecond, time.Duration(rqData.retryMaxPeriod)*time.Second)
	}
	log.Println("running retry")

	//prepare request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(string(c.Request().Header.Method()))
	req.Header.SetContentType(string(c.Request().Header.ContentType()))
	req.SetBody(c.Request().Body())
	log.Println(string(c.Request().Body()))
	log.Println(string(c.Request().URI().QueryString()))
	// req.SetRequestURI(string(c.Request().RequestURI()))
	// req.SetRequestURI("https://reqres.in/api/register")
	// req.SetRequestURI("http://httpstat.us/502")
	log.Println(rqData.reqURI + "?" + rqData.reqQuery)
	req.SetRequestURI(rqData.reqURI + "?" + rqData.reqQuery)
	//prepare response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := retrier.Run(func() error {
		// Perform the request
		err := fasthttp.Do(req, resp)
		if err != nil {
			log.Printf("Client get failed: %s\n", err)
			return err
		}

		switch {
		case err != nil:
			// request error - return it
			return err
		case resp.StatusCode() == 0 || resp.StatusCode() >= 500:
			// retryable StatusCode - return it
			log.Println("retryable HTTP status:", resp.StatusCode())
			return fmt.Errorf("retryable HTTP status: %s", http.StatusText(resp.StatusCode()))
		case resp.StatusCode() != 200:
			// non-retryable error - stop now
			log.Println("non-retryable HTTP statu", resp.StatusCode())
			return retry.Stop(fmt.Errorf("non-retryable HTTP status: %s", http.StatusText(resp.StatusCode())))
		}
		return nil
	})
	if err != nil {
		log.Println("error on retryFastHttpCoba", err)
	}
	// Verify the content type
	contentType := resp.Header.Peek("Content-Type")
	if bytes.Index(contentType, []byte("application/json")) != 0 {
		log.Printf("Expected content type application/json but got %s\n", contentType)
		// return err
	}

	// Do we need to decompress the response?
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		log.Println("Unzipping...")
		body, _ = resp.BodyGunzip()
	} else {
		body = resp.Body()
	}

	c.Response().SetBody(body)

	log.Printf("Response body is: %s", body)
	c.Response().Header.SetContentType(string(resp.Header.ContentType()))
	c.Response().Header.SetStatusCode(resp.StatusCode())
	// return c.Send(body)
	return nil
}

type reqData struct {
	reqUri         string
	isRetryable    string
	maxRetry       int
	retryPeriod    int
	retryMaxPeriod int
	reqURI         string
	reqQuery       string
}

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flowchartsman/retry"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kaleemubarok/recam/backend/app/models"
	"github.com/kaleemubarok/recam/backend/pkg/utils"
	"github.com/kaleemubarok/recam/backend/platform/cache"
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
	// Generate new RouteLog ID
	routeLog.ID = uuid.New()

	allParams := c.AllParams()

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
	_, err = connRedis.Get(context.Background(), allParams["slug"]).Result()

	if err == redis.Nil {
		routeData, err = db.GetRouteSlug(allParams["slug"])
		if err != nil {
			// Return status 400 and error message.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   fmt.Sprintf("no route matching with your requested %s, please try again", allParams["slug"]),
			})
		}
		routeDataJson, err := json.Marshal(routeData)
		if err != nil {
			return errors.New("there is an error on marshaling json data from redis")
		}

		_ = connRedis.Set(context.Background(), routeData.Slug, string(routeDataJson), 60*60*time.Second).Err()

	} else {
		// route datanya full from redis
		routeRedis, _ := connRedis.Get(context.Background(), allParams["slug"]).Result()
		err = json.Unmarshal([]byte(routeRedis), &routeData)
		if err != nil {
			return errors.New("there is an error on unmarshaling json data from redis")
		}
	}

	// kick out if token different from db
	if !strings.EqualFold(routeData.Token, allParams["token"]) {
		// Return status 404 and error message.
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "your token is not authorized, please contact admin",
		})
	}
	// kick out if method different from db
	if !strings.EqualFold(routeData.CommType, string(c.Request().Header.Method())) {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("error while matching request method %s, expected %s", string(c.Request().Header.Method()), routeData.CommType),
		})
	}

	// update context
	rqData := ReqData{
		isRetryable:    routeData.IsRetryable,
		maxRetry:       routeData.MaxRetry,
		retryPeriod:    routeData.RetryPeriod,
		retryMaxPeriod: 9999,
		reqURI:         routeData.HostAddr,
		reqQuery:       string(c.Request().URI().QueryString()),
		slug:           routeData.Slug,
		ID:             routeLog.ID,
		IPAddr:         c.IP(),
	}

	trialAttempt, _, encryptedRData := RetryFastHttp(c, &rqData)

	code, body := c.Response().StatusCode(), c.Response().Body()

	// Declare map to store request response data
	rqrs := jwt.MapClaims{}
	rqrs["query"] = string(c.Request().URI().QueryString())
	rqrs["rqHeader"] = string(c.Request().Header.Header())
	rqrs["rqBody"] = string(c.Request().Body())
	rqrs["rsBody"] = string(body)

	// Set correct content-type
	contentType := string(c.Response().Header.ContentType())
	c.Set(fiber.HeaderContentType, contentType)

	// encrypt part start
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
	// encrypt part end

	// set logs data
	routeLog.Data = string(encrypted)
	routeLog.ResponseCode = strconv.Itoa(code)
	routeLog.RouteID = routeData.ID /*TODO: cek ini*/
	routeLog.Type = string(c.Request().Header.Method())
	routeLog.CreatedAt = time.Now()
	routeLog.TrialAttempt = trialAttempt
	routeLog.Others = encryptedRData
	if code >= 500 || code == 0 {
		routeLog.IsResolved = "N"
	}

	// create a communication log
	go func() {
		if err := db.CreateLog(routeLog); err != nil {
			// Print error message on log console.
			log.Println("There is an error while inserting log to DB: ", err.Error())
		}
	}()

	/* //decrypt start
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
	// c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
	// c.Response().Header.SetStatusCode(code)
	// log.Println("sampe bawah")*/
	return c.Send(body)

	// return nil
}

func RetryFastHttp(c *fiber.Ctx, rqData *ReqData) (errorCode int, errResult error, strResult string) {
	// if 0 means deflault, those are 5, 1ms, 1s in order
	retrier := retry.NewRetrier(1, 0*time.Millisecond, 30*time.Second)
	if strings.EqualFold("Y", rqData.isRetryable) {
		retrier = retry.NewRetrier(rqData.maxRetry, time.Duration(rqData.retryPeriod)*time.Millisecond, 24*time.Hour)
	}

	// prepare request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(string(c.Request().Header.Method()))
	req.Header.SetContentType(string(c.Request().Header.ContentType()))
	req.SetBody(c.Request().Body())

	/*dummy host alternative*/
	// req.SetRequestURI(string(c.Request().RequestURI()))
	// req.SetRequestURI("https://reqres.in/api/register")
	// req.SetRequestURI("http://httpstat.us/502")

	// uri params preparation
	completeURI := c.Request().URI().String()
	aURISlug := strings.Split(completeURI, rqData.slug)
	aURIOnly := strings.Split(aURISlug[1], "?")
	path := aURIOnly[0]

	req.Header.SetUserAgent("ReCAM-PUPR")
	req.SetRequestURI(rqData.reqURI + path + "?" + rqData.reqQuery)
	if rqData.reqQuery == "" {
		req.SetRequestURI(rqData.reqURI + path)
	}
	// log.Println("reqURI", req.URI())
	// prepare response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// prepare temp var to handle retry error response
	lData := []models.LogData{}

	retryCounter := 0
	err := retrier.Run(func() error {
		retryCounter++
		// log.Println("retryer ke-", strconv.Itoa(retryCounter))
		// Perform the request
		err := fasthttp.Do(req, resp)
		if err != nil {
			log.Printf("Client get failed: %s\n", err)
			// return err
		}

		// set log data
		lData = append(lData, models.LogData{
			No:        retryCounter,
			ReqHeader: req.Header.String(),
			ReqBody:   string(req.Body()),
			ResCode:   resp.StatusCode(),
			ResBody:   string(resp.Body()),
			IPAddr:    rqData.IPAddr,
		})

		switch {
		case err != nil:
			// request error - return it
			return err
		case resp.StatusCode() == 0 || resp.StatusCode() >= 500:
			// retryable StatusCode - return it
			log.Println("retryable HTTP status:", resp.StatusCode(), string(resp.Body()))
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

	// log.Println("lData: ", lData)
	// set data for retry log
	rLogs := models.RetryLog{
		Total: retryCounter,
		ID:    rqData.ID.String(),
		Data:  lData,
	}
	rLogJSON, err := json.Marshal(rLogs)
	if err != nil {
		return retryCounter, err, ""
	}

	// encrypt part start
	key := os.Getenv("ENCRYPT_KEY")
	encrypted, err := utils.Encrypt(string(rLogJSON), key)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
	}
	// encrypt part end

	c.Response().SetBody(body)

	c.Response().Header.SetContentType(string(resp.Header.ContentType()))
	c.Response().Header.SetStatusCode(resp.StatusCode())

	return retryCounter, nil, encrypted
}

type ReqData struct {
	isRetryable    string
	maxRetry       int
	retryPeriod    int
	retryMaxPeriod int
	reqURI         string
	reqQuery       string
	slug           string
	ID             uuid.UUID
	IPAddr         string
}

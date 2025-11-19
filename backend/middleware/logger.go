// middleware/logger.go → 100% JALAN DI FIBER v2.52+ (FINAL VERSION)
package middleware

import (
	"encoding/json"
	"strings"
	"time"

	"notes-app-backend/config"
	"notes-app-backend/models"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Copy request body
		reqBody := make([]byte, len(c.Body()))
		copy(reqBody, c.Body())

		payloadStr := ""
		if len(reqBody) > 0 {
			if len(reqBody) > 1000 {
				payloadStr = string(reqBody[:1000]) + "...(truncated)"
			} else {
				payloadStr = string(reqBody)
			}
		}

		// 2. Headers + mask Authorization
		headers := make(map[string]string)
		for k, values := range c.GetReqHeaders() {
			if len(values) > 0 {
				val := values[0]
				if strings.ToLower(k) == "authorization" && val != "" {
					if strings.HasPrefix(val, "Bearer ") {
						token := strings.TrimPrefix(val, "Bearer ")
						masked := "****"
						if len(token) > 6 {
							masked = token[len(token)-6:]
						}
						headers[k] = "Bearer " + masked + "...(masked)"
					} else {
						headers[k] = "****(masked)"
					}
				} else {
					headers[k] = val
				}
			}
		}
		headersJSON, _ := json.Marshal(headers)

		// 3. Jalankan semua handler dulu
		err := c.Next()

		// 4. Setelah semua selesai, ambil response body (ini sudah final)
		responseBody := string(c.Response().Body())
		if len(responseBody) > 2000 {
			responseBody = responseBody[:2000] + "...(truncated)"
		}

		// 5. User ID dari JWT
		var userID *uint
		if uid := c.Locals("user_id"); uid != nil {
			id := uid.(uint)
			userID = &id
		}

		// 6. Simpan log (async)
		logEntry := models.Log{
			Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
			Method:       c.Method(),
			Endpoint:     c.Path(),
			Headers:      string(headersJSON),
			Payload:      payloadStr,
			ResponseBody: responseBody,
			StatusCode:   c.Response().StatusCode(),
			UserID:       userID,
		}

		go config.DB.Create(&logEntry)

		return err
	}
}
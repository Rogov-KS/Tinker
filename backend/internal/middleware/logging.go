package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Логируем входящий запрос
		query := c.Request.URL.Query()
		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		// Маскируем пароль в body если есть
		var body interface{}
		if c.Request.Body != nil && c.ContentType() == "application/json" {
			// Для упрощения, просто логируем что body есть
			// В production можно парсить и маскировать
			body = "[body present]"
		}

		slog.Info("Incoming request",
			"method", method,
			"path", path,
			"ip", ip,
			"userAgent", userAgent,
			"query", query,
			"params", params,
			"body", body,
		)

		// Обрабатываем запрос
		c.Next()

		// Логируем ответ
		duration := time.Since(start)
		status := c.Writer.Status()

		slog.Info("Outgoing response",
			"method", method,
			"path", path,
			"status", status,
			"duration", duration.Milliseconds(),
		)

		// Логируем ошибки если есть
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				slog.Error("Request error",
					"method", method,
					"path", path,
					"error", err.Error(),
				)
			}
		}
	}
}


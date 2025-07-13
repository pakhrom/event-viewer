package main

import (
	"net/http"

	"github.com/billowdev/clog"
	"github.com/labstack/echo/v4"
)

func CustomErrorHandler(err error, c echo.Context) {
	var code int
	var message interface{}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	} else {
		code = http.StatusInternalServerError
		message = err.Error()
	}

	clog.Error("[%d] IP: %s | %s %s\n%v", code, c.RealIP(), c.Request().Method, c.Request().URL.Path, message)

	// Стандартный ответ клиенту
	if !c.Response().Committed {
		c.JSON(code, map[string]interface{}{
			"error": message,
		})
	}
}

type responseWriter struct {
	echo.Response
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.Response.WriteHeader(code)
}

// Middleware для логирования успешных запросов
func ClogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rw := &responseWriter{Response: *c.Response(), status: 0}
		c.SetResponse(&rw.Response)
		err := next(c)

		code := rw.status
		if code == 0 {
			code = c.Response().Status
		}

		method := c.Request().Method
		path := c.Request().URL.Path
		clientIP := c.RealIP() // Получаем IP клиента

		if code < 400 {
			clog.Info("[%d] IP: %s | %s %s ", code, clientIP, method, path)
		}
		return err
	}
}

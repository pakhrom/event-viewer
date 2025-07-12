package main

import (
	"errors"
	"help/db"
	"help/routers"
	"net/http"

	"github.com/billowdev/clog"
	"github.com/labstack/echo/v4"

	docs "help/docs" // импорт с именем для доступа к SwaggerInfo

	echoSwagger "github.com/swaggo/echo-swagger"
)

// pingResponse — пример простого обработчика
func pingResponse(g echo.Context) error {
	return g.JSON(http.StatusOK, "pong")
}

// Кастомный responseWriter для захвата статуса ответа
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

		// Получаем статус-код, который реально отправлен клиенту
		code := rw.status
		if code == 0 {
			code = c.Response().Status
		}

		method := c.Request().Method
		path := c.Request().URL.Path

		// Логируем только успешные ответы (коды < 400)
		if code < 400 {
			clog.Info("%s %s | %d", method, path, code)
		}
		return err
	}
}

func main() {
	DBConnection, err := db.InitDB("mongodb+srv://TestProjectEducationalEvents:AtoG0nw1fonvX6BR@timetableproject.mt2imbb.mongodb.net/?retryWrites=true&w=majority&appName=TimetableProject")
	if err != nil {
		panic(err)
	}
	defer DBConnection.CloseConnection()

	router := echo.New()
	clog.SetShowGoroutineID(false)
	router.Use(ClogMiddleware)

	// Кастомный обработчик ошибок для корректного логирования ошибок
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		var code int
		var message interface{}

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message
		} else {
			code = http.StatusInternalServerError
			message = err.Error()
		}

		clog.Error("%s %s | %d | %v", c.Request().Method, c.Request().URL.Path, code, message)

		// Стандартный ответ клиенту
		if !c.Response().Committed {
			c.JSON(code, map[string]interface{}{
				"error": message,
			})
		}
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/ping", pingResponse)
		}

		eventsGroup := v1.Group("/events")
		{
			eventsGroup.GET("/:id", routers.RequestEvent(&DBConnection.EducationalEvents))
			eventsGroup.GET("", routers.RequestEvents(&DBConnection.EducationalEvents))
			eventsGroup.POST("", routers.CreateEvent(&DBConnection.EducationalEvents))
			eventsGroup.DELETE("/:id", routers.DeleteEvent(&DBConnection.EducationalEvents))
			eventsGroup.PUT("/:id", routers.UpdateEvent(&DBConnection.EducationalEvents))
		}
	}

	router.GET("/swagger/*", echoSwagger.WrapHandler) // Swagger UI

	// Запуск сервера
	if err := router.Start("0.0.0.0:8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		clog.Error("failed to start server: %v", err)
	}
}

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
func main() {
	router := echo.New()
	clog.SetShowGoroutineID(false)
	router.Use(ClogMiddleware)

	// Кастомный обработчик ошибок для корректного логирования ошибок
	router.HTTPErrorHandler = CustomErrorHandler

	DBConnection, err := db.InitDB("mongodb+srv://TestProjectEducationalEvents:AtoG0nw1fonvX6BR@timetableproject.mt2imbb.mongodb.net/?retryWrites=true&w=majority&appName=TimetableProject")
	if err != nil {
		clog.Error("Failed to connect to MongoDB!\n%s", err.Error())
	}
	defer DBConnection.CloseConnection()

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

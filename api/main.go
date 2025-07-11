package main

import (
	"errors"
	"help/db"

	_ "help/docs" // ваш пакет с swagger-документацией
	"help/routers"
	"net/http"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	docs "help/docs" // импорт с именем для доступа к SwaggerInfo

	echoSwagger "github.com/swaggo/echo-swagger"
)

// gin-swagger middleware
// swagger embed files

func pingResponse(g echo.Context) error {
	return g.JSON(http.StatusOK, "pong")
}

// @title       My App API
// @version     1.0
// @description This is a sample Gin server.
// @termsOfService http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://example.com/support
// @contact.email  support@example.com

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http
func main() {

	DBConnection, err := db.InitDB("mongodb+srv://TestProjectEducationalEvents:AtoG0nw1fonvX6BR@timetableproject.mt2imbb.mongodb.net/?retryWrites=true&w=majority&appName=TimetableProject")
	if err != nil {
		panic(err)
	}
	defer DBConnection.CloseConnection()

	router := echo.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/ping", pingResponse)
		}

		eventsGroup := v1.Group("/events")
		{
			eventsGroup.GET("", routers.RequestEvents(&DBConnection.EducationalEvents))
			// eventsGroup.POST("")
		}
	}

	router.GET("/swagger/*", echoSwagger.WrapHandler) // Swagger UI
	// Start server
	if err := router.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

}

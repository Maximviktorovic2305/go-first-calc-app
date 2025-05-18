package main

import (
	"fmt"
	"log"
	"server/internal/calculationService"
	"server/internal/db"
	"server/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Инициализация базы данных
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	e := echo.New()         

	calcRepo := calculationService.NewCalculationRepository(database)
	calcService := calculationService.NewCalculationService(calcRepo)
	calculationsHandlers := handlers.NewCalculationHandler(calcService)

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", calculationsHandlers.GetCalculations)
	e.POST("/calculations", calculationsHandlers.PostCalculations)
	e.PATCH("/calculations/:id", calculationsHandlers.PatchCalculations)
	e.DELETE("/calculations/:id", calculationsHandlers.DeleteCalculations)

	fmt.Println("Server is running on http://localhost:8080")
	e.Start(":8080")
}

package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Модель для хранения вычислений в базе данных
type Calculation struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

// Структура для запроса вычисления
type CalculationRequest struct {
	Expression string `json:"expression"`
}

var db *gorm.DB

// Инициализация базы данных
func initDB() error {
	dsn := "host=localhost user=postgres password=admin dbname=go-calc port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Автоматическая миграция таблицы
	return db.AutoMigrate(&Calculation{})
}

// Функция для вычисления выражения
func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

// Получение всех расчетов из базы
func getCalculations(c echo.Context) error {
	var calculations []Calculation
	if err := db.Find(&calculations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch calculations"})
	}
	return c.JSON(http.StatusOK, calculations)
}

// Создание нового расчета
func postCalculations(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Expression"})
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}

	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save calculation"})
	}

	return c.JSON(http.StatusCreated, calc)
}

// Обновление существующего расчета
func patchCalculations(c echo.Context) error {
	id := c.Param("id")
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	var calc Calculation
	if err := db.First(&calc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Calculation not found"})
	}

	calc.Expression = req.Expression
	calc.Result = result

	if err := db.Save(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update calculation"})
	}

	return c.JSON(http.StatusOK, calc)
}

// Удаление расчета
func deleteCalculations(c echo.Context) error {
	id := c.Param("id")
	if err := db.Delete(&Calculation{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete calculation"})
	}
	return c.NoContent(http.StatusNoContent)
}

func main() {
	// Инициализация базы данных
	if err := initDB(); err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	fmt.Println("Server is running on http://localhost:8080")
	e.Start(":8080")
}

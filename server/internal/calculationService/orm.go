package calculationService

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
package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ByteNode1/GolangCalculator/internal/calculator"
	"github.com/gorilla/mux"
)

// Обработчик HTTP запросов для вычисления выражений
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := calculator.Calc(req.Expression)
	if err != nil {
		var status int
		var message string

		switch {
		case errors.Is(err, calculator.ErrInvalidExpression),
			errors.Is(err, calculator.ErrInvalidCharacter),
			errors.Is(err, calculator.ErrConsecutiveOperators),
			errors.Is(err, calculator.ErrMismatchedParentheses):
			status = http.StatusUnprocessableEntity
			message = "Expression is not valid"
		case errors.Is(err, calculator.ErrDivisionByZero):
			status = http.StatusBadRequest
			message = "Division by zero"
		default:
			status = http.StatusInternalServerError
			message = "Internal server error"
		}

		http.Error(w, `{"error": "`+message+`"}`, status)
		return
	}

	response := map[string]float64{"result": result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Запуск HTTP-сервера
func Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", calculateHandler).Methods("POST")

	return http.ListenAndServe(":8080", r)
}

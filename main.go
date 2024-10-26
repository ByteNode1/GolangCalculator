package main

import (
	"errors"
	"strconv"
	"unicode"
)

// Предопределенные ошибки
var (
	ErrMismatchedParentheses = errors.New("несоответствие открывающих и закрывающих скобок")
	ErrInvalidExpression     = errors.New("неверное выражение")
	ErrInvalidCharacter      = errors.New("недопустимый символ в выражении")
	ErrConsecutiveOperators  = errors.New("два оператора подряд")
	ErrDivisionByZero        = errors.New("деление на ноль")
)

// Определение приоритета оператора
func precedence(op byte) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

// Выполнение арифметической операции
func ParseOp(a, b float64, op byte) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	}
	return 0, ErrInvalidExpression
}

// Функция для вычисления выражения
func Calc(expression string) (float64, error) {
	var values []float64 // Стек для чисел
	var ops []byte       // Стек для операторов

	balance := 0      // Счетчик для проверки баланса скобок
	lastWasOp := true // Флаг для проверки последовательности символов
	numBuffer := ""   // Буфер для накопления чисел

	for _, char := range expression {
		switch {
		case char == ' ':
			continue

		case unicode.IsDigit(char) || char == '.':
			// Если символ - часть числа, добавляем его в буфер
			numBuffer += string(char)
			lastWasOp = false

		case char == '(':
			if len(numBuffer) > 0 { // Если скобка после числа - ошибка
				return 0, ErrInvalidExpression
			}
			balance++
			lastWasOp = true
			ops = append(ops, byte(char))

		case char == ')':
			if balance == 0 {
				return 0, ErrMismatchedParentheses
			}
			if len(numBuffer) > 0 { // Если число перед скобкой - добавляем его в стек
				val, err := strconv.ParseFloat(numBuffer, 64)
				if err != nil {
					return 0, ErrInvalidExpression
				}
				values = append(values, val)
				numBuffer = ""
			}
			balance--
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if len(values) < 2 {
					return 0, ErrInvalidExpression
				}
				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				result, err := ParseOp(val1, val2, op)
				if err != nil {
					return 0, err
				}
				values = append(values, result)
			}
			if len(ops) == 0 {
				return 0, ErrMismatchedParentheses
			}
			ops = ops[:len(ops)-1]
			lastWasOp = false

		case char == '+' || char == '-' || char == '*' || char == '/':
			if lastWasOp { // Проверка на два оператора подряд
				return 0, ErrConsecutiveOperators
			}
			if len(numBuffer) > 0 { // Если число перед оператором - добавляем его в стек
				val, err := strconv.ParseFloat(numBuffer, 64)
				if err != nil {
					return 0, ErrInvalidExpression
				}
				values = append(values, val)
				numBuffer = ""
			}
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(byte(char)) {
				if len(values) < 2 {
					return 0, ErrInvalidExpression
				}
				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				result, err := ParseOp(val1, val2, op)
				if err != nil {
					return 0, err
				}
				values = append(values, result)
			}
			ops = append(ops, byte(char))
			lastWasOp = true

		default:
			return 0, ErrInvalidCharacter
		}
	}

	// Обработка последнего числа в буфере
	if len(numBuffer) > 0 {
		val, err := strconv.ParseFloat(numBuffer, 64)
		if err != nil {
			return 0, ErrInvalidExpression
		}
		values = append(values, val)
	}

	// Проверка на несбалансированные скобки
	if balance != 0 {
		return 0, ErrMismatchedParentheses
	}

	// Выполняем оставшиеся операции
	for len(ops) > 0 {
		if len(values) < 2 {
			return 0, ErrInvalidExpression
		}
		val2 := values[len(values)-1]
		values = values[:len(values)-1]
		val1 := values[len(values)-1]
		values = values[:len(values)-1]
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		result, err := ParseOp(val1, val2, op)
		if err != nil {
			return 0, err
		}
		values = append(values, result)
	}

	// Проверка, что в стеке осталось одно значение - результат
	if len(values) != 1 {
		return 0, ErrInvalidExpression
	}
	return values[0], nil
}

// func main() {
// 	result, err := Calc("1/2")
// 	if err != nil {
// 		fmt.Println("Ошибка:", err)
// 	} else {
// 		fmt.Println("Результат:", result)
// 	}
// }

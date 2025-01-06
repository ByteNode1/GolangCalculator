package calculator

import (
	"strconv"
	"unicode"
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
func parseOp(a, b float64, op byte) (float64, error) {
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
				result, err := parseOp(val1, val2, op)
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
				result, err := parseOp(val1, val2, op)
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

	if len(numBuffer) > 0 {
		val, err := strconv.ParseFloat(numBuffer, 64)
		if err != nil {
			return 0, ErrInvalidExpression
		}
		values = append(values, val)
	}

	if balance != 0 {
		return 0, ErrMismatchedParentheses
	}

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
		result, err := parseOp(val1, val2, op)
		if err != nil {
			return 0, err
		}
		values = append(values, result)
	}

	if len(values) != 1 {
		return 0, ErrInvalidExpression
	}
	return values[0], nil
}

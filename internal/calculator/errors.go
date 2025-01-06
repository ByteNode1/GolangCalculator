package calculator

import "errors"

// ашипки
var (
	ErrMismatchedParentheses = errors.New("несоответствие открывающих и закрывающих скобок")
	ErrInvalidExpression     = errors.New("неверное выражение")
	ErrInvalidCharacter      = errors.New("недопустимый символ в выражении")
	ErrConsecutiveOperators  = errors.New("два оператора подряд")
	ErrDivisionByZero        = errors.New("деление на ноль")
)

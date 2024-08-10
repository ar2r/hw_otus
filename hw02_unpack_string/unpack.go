package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var builder strings.Builder

	originalLen := len(str)
	// Дополнить строку лишним пробелом для удобства обработки
	str += " "

	for i := 0; i < originalLen; i++ {
		currentIsDigit := isDigit(rune(str[i]))
		currentIsChar := !currentIsDigit

		nextIsDigit := isDigit(rune(str[i+1]))
		nextIsChar := !nextIsDigit

		if currentIsChar && nextIsDigit {
			// Сценарий 1 - буква + цифра
			repeatedChar := strings.Repeat(string(str[i]), int(str[i+1]-'0'))
			builder.WriteString(repeatedChar)
			i++
		}

		if currentIsChar && nextIsChar {
			// Сценарий 2 - буква + буква
			builder.WriteString(string(str[i]))
		}

		if currentIsDigit {
			// Сценарий 3 - цифра + любой символ - ошибка
			// Например цифра в конце строки или две цифры подряд
			return "", ErrInvalidString
		}
	}

	return builder.String(), nil
}

func isDigit(r rune) bool {
	_, e := strconv.Atoi(string(r))
	return e == nil
}

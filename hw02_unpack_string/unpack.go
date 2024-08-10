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
		currentIsChar := isChar(rune(str[i]))
		nextIsChar := isChar(rune(str[i+1]))

		if currentIsChar && !nextIsChar {
			// Сценарий 1 - буква + цифра
			repeatCount, _ := strconv.Atoi(string(str[i+1]))
			repeatedChar := strings.Repeat(string(str[i]), repeatCount)
			builder.WriteString(repeatedChar)
			i++
		}

		if currentIsChar && nextIsChar {
			// Сценарий 2 - буква + буква
			builder.WriteString(string(str[i]))
		}

		if !currentIsChar {
			// Сценарий 3 - цифра + любой символ - ошибка
			// Например цифра в конце строки или две цифры подряд
			return "", ErrInvalidString
		}
	}

	return builder.String(), nil
}

func isChar(r rune) bool {
	_, err := strconv.Atoi(string(r))
	return err != nil
}

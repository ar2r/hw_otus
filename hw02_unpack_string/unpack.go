package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var builder strings.Builder
	runes := []rune(str)
	// Дополнить строку лишним пробелом для удобства обработки
	runes = append(runes, ' ')

	for i := 0; i < len(runes)-1; i++ {
		currentIsChar := isChar(runes[i])
		nextIsChar := isChar(runes[i+1])

		if currentIsChar && !nextIsChar { //nolint:gocritic
			// Сценарий 1 - буква + цифра
			repeatCount, err := strconv.Atoi(string(runes[i+1]))
			if err != nil {
				return "", ErrInvalidString
			}
			repeatedChar := strings.Repeat(string(runes[i]), repeatCount)
			builder.WriteString(repeatedChar)
			i++
		} else if currentIsChar {
			// Сценарий 2 - буква + буква
			builder.WriteString(string(runes[i]))
		} else {
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

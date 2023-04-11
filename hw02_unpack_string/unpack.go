package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	//Case: Пустая строка
	if len(s) == 0 {
		return "", nil
	}

	//Case: 1-ый символ - цифра
	if unicode.IsDigit(rune(s[0])) {
		return "", ErrInvalidString
	}

	var result strings.Builder

	backSlashFlag := false
	digitFlag := false

	for _, elem := range s {

		switch {
		//Case: в строке c back-tricks`` после обратного слеша идет не цифра или не обратный слеш - ошибка
		case backSlashFlag && !(unicode.IsDigit(elem) || elem == 92):
			return "", ErrInvalidString
		//Case: в любой строке после цифры встречаем очередню цифру
		case unicode.IsDigit(elem) && digitFlag:
			return "", ErrInvalidString
		//Case: в строке c back-tricks`` после обратного слеша идет цифра или обратный слеш - добавляем в Builder
		case (unicode.IsDigit(elem) || elem == 92) && backSlashFlag:
			backSlashFlag = false
			fallthrough
		//Case: если буква или пробел (или табуляция или переход на следующую строку) - добавляем в Builder
		case unicode.IsLetter(elem) || unicode.IsSpace(elem):
			result.WriteRune(elem)
			digitFlag = false
		//Case: если цифра:
		case unicode.IsDigit(elem):
			repeatCount, _ := strconv.Atoi(string(elem))
			// Проверка для теста с цифрой 0. По сути, n > 0 - добавляем n-1 раз, иначе - удаляем из Builder добавленную букву
			if repeatCount > 0 {
				bulderRunes := []rune(result.String())
				lastInsertedLetter := rune(bulderRunes[(len(bulderRunes) - 1)])
				result.WriteString(strings.Repeat(string(lastInsertedLetter), repeatCount-1))
				digitFlag = true
			} else {
				temp_s := []rune(result.String())
				result.Reset()
				result.WriteString(string(temp_s[:len(temp_s)-1]))
			}
		//Case: в строке с bak-tricks`` если обратный слеш
		case elem == 92:
			backSlashFlag = true
			digitFlag = false
		}

	}
	return result.String(), nil
}

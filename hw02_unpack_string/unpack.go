package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}

	if unicode.IsDigit(rune(s[0])) {
		return "", ErrInvalidString
	}

	var result strings.Builder

	backSlashFlag := false
	digitFlag := false

	for _, elem := range s {
		switch {
		case backSlashFlag && !(unicode.IsDigit(elem) || elem == 92):
			return "", ErrInvalidString
		case unicode.IsDigit(elem) && digitFlag:
			return "", ErrInvalidString
		case (unicode.IsDigit(elem) || elem == 92) && backSlashFlag:
			backSlashFlag = false
			fallthrough
		case unicode.IsLetter(elem) || unicode.IsSpace(elem):
			result.WriteRune(elem)
			digitFlag = false
		case unicode.IsDigit(elem):
			repeatCount, _ := strconv.Atoi(string(elem))
			if repeatCount > 0 {
				bulderRunes := []rune(result.String())
				lastInsertedLetter := bulderRunes[(len(bulderRunes) - 1)]
				result.WriteString(strings.Repeat(string(lastInsertedLetter), repeatCount-1))
				digitFlag = true
			} else {
				tempS := []rune(result.String())
				result.Reset()
				result.WriteString(string(tempS[:len(tempS)-1]))
			}
		case elem == 92:
			backSlashFlag = true
			digitFlag = false
		}
	}
	return result.String(), nil
}

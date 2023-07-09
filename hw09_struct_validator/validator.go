package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	// global errors.
	errInvalidType     = errors.New("inserted argument is not a struct")
	errValidationError = errors.New("validation completed with errors")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	for _, elem := range v {
		b.WriteString(elem.Field + "-->" + elem.Err.Error() + "\n")
	}

	return b.String()
}

func validateInt(key, field string, value int, ve *ValidationErrors) {
	keys := strings.Split(key, "|")
	for _, k := range keys {
		subkeys := strings.Split(k, ":")
		switch subkeys[0] {
		case "min":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid key value, key min=" + subkeys[1]),
				})
			} else if value < res {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid value, ocured value is less than min=" + subkeys[1]),
				})
			}
		case "max":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid key value, key max=" + subkeys[1]),
				})
			} else if value > res {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid value, ocured value is greater than max=" + subkeys[1]),
				})
			}
		case "in":
			values := strings.Split(subkeys[1], ",")
			found := false
			for _, v := range values {
				res, err := strconv.Atoi(v)
				if err != nil {
					*ve = append(*ve, ValidationError{
						Field: field,
						Err:   errors.New("invalid key value, key in=" + v),
					})
				} else if res == value {
					found = true
					break
				}
			}
			if !found {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid value, value is not in \"in\" list"),
				})
			}
		default:
			*ve = append(*ve, ValidationError{
				Field: field,
				Err:   errors.New("invalid validate subtag"),
			})
		}
	}
}

func validateString(key, field, value string, ve *ValidationErrors) {
	keys := strings.Split(key, "|")
	for _, k := range keys {
		subkeys := strings.Split(k, ":")
		switch subkeys[0] {
		case "len":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid key value, key len=" + subkeys[1]),
				})
			} else if len(value) != res {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid value, value's length is greater than len=" + subkeys[1]),
				})
			}
		case "regexp":
			rx, err := regexp.Compile(subkeys[1])
			if err != nil {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid key value, key regexp = " + subkeys[1]),
				})
			} else {
				res := rx.FindString(value)
				if res != value {
					*ve = append(*ve, ValidationError{
						Field: field,
						Err:   errors.New("invalid value, value doesn't match regular expression"),
					})
				}
			}
		case "in":
			if !strings.Contains(subkeys[1], value) {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errors.New("invalid value, value is not in \"in\" list"),
				})
			}
		default:
			*ve = append(*ve, ValidationError{
				Field: field,
				Err:   errors.New("invalid validate subtag"),
			})
		}
	}
}

func Validate(v interface{}) (string, error) {
	ve := ValidationErrors{}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ve.Error(), errInvalidType
	}

	fieldsList := reflect.VisibleFields(rv.Type())

	for _, field := range fieldsList {
		if param, ok := field.Tag.Lookup("validate"); ok {
			switch field.Type.Kind() {
			case reflect.Slice:
				switch field.Type.Elem().Kind() {
				case reflect.Int:
					slice := rv.FieldByName(field.Name).Interface().([]int)
					for _, elem := range slice {
						validateInt(param, field.Name, elem, &ve)
					}
				case reflect.String:
					slice := rv.FieldByName(field.Name).Interface().([]string)
					for _, elem := range slice {
						validateString(param, field.Name, elem, &ve)
					}
				default:
					ve = append(ve, ValidationError{Field: field.Name, Err: errors.New("unknown slice field type")})
				}
			case reflect.Int:
				validateInt(param, field.Name, rv.FieldByName(field.Name).Interface().(int), &ve)
			case reflect.String:
				validateString(param, field.Name, rv.FieldByName(field.Name).Interface().(string), &ve)
			default:
				ve = append(ve, ValidationError{Field: field.Name, Err: errors.New("unknown field type")})
			}
		}
	}

	if len(ve) > 0 {
		return ve.Error(), errValidationError
	}
	return ve.Error(), nil
}

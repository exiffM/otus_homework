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
	globalVe           ValidationErrors
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

func Validate(v interface{}) error {
	ve := make(ValidationErrors, 0)
	defer func() {
		globalVe = ve
	}()
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return errInvalidType
	}

	fieldsList := reflect.VisibleFields(rv.Type())

	for _, field := range fieldsList {
		if param, ok := field.Tag.Lookup("validate"); ok {
			value := reflect.ValueOf(rv.FieldByName(field.Name).Interface())
			switch field.Type.Kind() { //nolint:exhaustive
			case reflect.Slice:
				switch field.Type.Elem().Kind() { //nolint:exhaustive
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					for i := 0; i < value.Len(); i++ {
						validateInt(param, field.Name, int(value.Index(i).Int()), &ve)
					}
				case reflect.String:
					for i := 0; i < value.Len(); i++ {
						validateString(param, field.Name, value.Index(i).String(), &ve)
					}
				default:
					ve = append(ve, ValidationError{Field: field.Name, Err: errors.New("unknown slice field type")})
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				validateInt(param, field.Name, int(value.Int()), &ve)

			case reflect.String:
				validateString(param, field.Name, value.String(), &ve)
			default:
				ve = append(ve, ValidationError{Field: field.Name, Err: errors.New("unknown field type")})
			}
		}
	}

	if len(ve) > 0 {
		return errValidationError
	}
	return nil
}

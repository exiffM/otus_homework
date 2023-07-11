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

// type SafeErrors struct {
// 	globalVe ValidationErrors
// 	m        sync.RWMutex
// }

// func (se *SafeErrors) Assign(ve ValidationErrors) {
// 	se.m.Lock()
// 	se.globalVe = ve
// 	se.m.Unlock()
// }

// func (se *SafeErrors) Show() ValidationErrors {
// 	se.m.RLock()
// 	defer se.m.RUnlock()
// 	return se.globalVe
// }

var (
	// global errors.
	errInvalidType = errors.New("inserted argument is not a struct")
	errMinKey      = errors.New("invalid key value, key min")
	errMaxKey      = errors.New("invalid key value, key max")
	errInKey       = errors.New("invalid key value, key in")
	errLenKey      = errors.New("invalid key value, key len")
	errRegexpKey   = errors.New("invalid key value, key regexp")
	errInvelidTag  = errors.New("invalid validate tag")
	// validation errors.
	errValidationError = errors.New("validation completed with errors")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	return errValidationError.Error()
}

func validateInt(key, field string, value int, ve *ValidationErrors) error {
	keys := strings.Split(key, "|")
	for _, k := range keys {
		subkeys := strings.Split(k, ":")
		switch subkeys[0] {
		case "min":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				return errMinKey
			} else if value < res {
				vErr := ValidationError{
					Field: field,
					Err:   errors.New("at " + field + ": invalid value, ocured value is less than min"),
				}
				*ve = append(*ve, vErr)
				errValidationError = errors.Join(errValidationError, vErr.Err)
			}
		case "max":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				return errMaxKey
			} else if value > res {
				vErr := ValidationError{
					Field: field,
					Err:   errors.New("at " + field + ": invalid value, ocured value is greater than max"),
				}
				*ve = append(*ve, vErr)
				errValidationError = errors.Join(errValidationError, vErr.Err)
			}
		case "in":
			values := strings.Split(subkeys[1], ",")
			found := false
			for _, v := range values {
				res, err := strconv.Atoi(v)
				if err != nil {
					return errInKey
				} else if res == value {
					found = true
					break
				}
			}
			if !found {
				vErr := ValidationError{
					Field: field,
					Err:   errors.New("at " + field + ": invalid value, value is not in \"in\" list"),
				}
				*ve = append(*ve, vErr)
				errValidationError = errors.Join(errValidationError, vErr.Err)
			}
		default:
			return errInvelidTag
		}
	}
	return nil
}

func validateString(key, field, value string, ve *ValidationErrors) error {
	keys := strings.Split(key, "|")
	for _, k := range keys {
		subkeys := strings.Split(k, ":")
		switch subkeys[0] {
		case "len":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				return errLenKey
			} else if len(value) != res {
				vErr := ValidationError{
					Field: field,
					Err:   errors.New("at " + field + ": invalid value, value's length is greater than len"),
				}
				*ve = append(*ve, vErr)
				errValidationError = errors.Join(errValidationError, vErr.Err)
			}
		case "regexp":
			rx, err := regexp.Compile(subkeys[1])
			if err != nil {
				return errRegexpKey
			} else {
				res := rx.FindString(value)
				if res != value {
					vErr := ValidationError{
						Field: field,
						Err:   errors.New("at " + field + ": invalid value, value doesn't match regular expression"),
					}
					*ve = append(*ve, vErr)
					errValidationError = errors.Join(errValidationError, vErr.Err)
				}
			}
		case "in":
			if !strings.Contains(subkeys[1], value) {
				vErr := ValidationError{
					Field: field,
					Err:   errors.New("at " + field + ": invalid value, value is not in \"in\" list"),
				}
				*ve = append(*ve, vErr)
				errValidationError = errors.Join(errValidationError, vErr.Err)
			}
		default:
			return errInvelidTag
		}
	}
	return nil
}

func Validate(v interface{}) error {
	ve := make(ValidationErrors, 0)
	// errValidationError = errors.New("validation completed with errors")
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
						err := validateInt(param, field.Name, int(value.Index(i).Int()), &ve)
						if err != nil {
							return err
						}
					}
				case reflect.String:
					for i := 0; i < value.Len(); i++ {
						err := validateString(param, field.Name, value.Index(i).String(), &ve)
						if err != nil {
							return err
						}
					}
				default:
					v := ValidationError{Field: field.Name, Err: errors.New("at " + field.Name + ": unknown slice field type")}
					ve = append(ve, v)
					errValidationError = errors.Join(v.Err)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				err := validateInt(param, field.Name, int(value.Int()), &ve)
				if err != nil {
					return err
				}

			case reflect.String:
				err := validateString(param, field.Name, value.String(), &ve)
				if err != nil {
					return err
				}
			default:
				v := ValidationError{Field: field.Name, Err: errors.New("at " + field.Name + ": unknown field type")}
				ve = append(ve, v)
				errValidationError = errors.Join(v.Err)
			}
		}
	}

	if len(ve) > 0 {
		return errValidationError
	}
	return nil
}

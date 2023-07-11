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
	errInvalidType = errors.New("inserted argument is not a struct")
	errMinKey      = errors.New("invalid key value, key min")
	errMaxKey      = errors.New("invalid key value, key max")
	errInKey       = errors.New("invalid key value, key in")
	errLenKey      = errors.New("invalid key value, key len")
	errRegexpKey   = errors.New("invalid key value, key regexp")
	errInvelidTag  = errors.New("invalid validate tag")
	// validation errors.
	errValidationMinError     = errors.New("invalid value, ocured value is less than min")
	errValidationMaxError     = errors.New("invalid value, ocured value is greater than max")
	errValidationInError      = errors.New("invalid value, value is not in \"in\" list")
	errValidationRegexpError  = errors.New("invalid value, value's length is greater than len")
	errValidationLenError     = errors.New("invalid value, value doesn't match regular expression")
	errValidationUnknownType  = errors.New("unknown field type")
	errValidationUnknownSlice = errors.New("unknown slice field type")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var err error
	for _, elem := range v {
		err = errors.Join(err, elem.Err)
	}
	return err.Error()
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
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errValidationMinError,
				})
			}
		case "max":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				return errMaxKey
			} else if value > res {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errValidationMaxError,
				})
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
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errValidationInError,
				})
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
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errValidationLenError,
				})
			}
		case "regexp":
			rx, err := regexp.Compile(subkeys[1])
			if err != nil {
				return errRegexpKey
			} else {
				res := rx.FindString(value)
				if res != value {
					*ve = append(*ve, ValidationError{
						Field: field,
						Err:   errValidationRegexpError,
					})
				}
			}
		case "in":
			if !strings.Contains(subkeys[1], value) {
				*ve = append(*ve, ValidationError{
					Field: field,
					Err:   errValidationInError,
				})
			}
		default:
			return errInvelidTag
		}
	}
	return nil
}

func Validate(v interface{}) error {
	ve := make(ValidationErrors, 0)
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
					v := ValidationError{Field: field.Name, Err: errValidationUnknownSlice}
					ve = append(ve, v)
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
				v := ValidationError{Field: field.Name, Err: errValidationUnknownType}
				ve = append(ve, v)
			}
		}
	}

	if len(ve) > 0 {
		var err error
		for _, elem := range ve {
			err = errors.Join(err, elem.Err)
		}
		return err
	}
	return nil
}

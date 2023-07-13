package hw09structvalidator

import (
	"errors"
	"fmt"
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
	errValidationMin          = errors.New("invalid value, ocured value is less than min")
	errValidationMax          = errors.New("invalid value, ocured value is greater than max")
	errValidationIn           = errors.New("invalid value, value is not in \"in\" list")
	errValidationRegexp       = errors.New("invalid value, value's length is greater than len")
	errValidationLen          = errors.New("invalid value, value doesn't match regular expression")
	errValidationUnknownType  = errors.New("unknown field type")
	errValidationUnknownSlice = errors.New("unknown slice field type")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	msg := ""
	for idx, elem := range v {
		if idx == 0 {
			msg = "Validation finished with errors"
		}
		msg = fmt.Sprintf("%s\n%s --> %s", msg, elem.Field, elem.Err.Error())
	}
	return msg
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
				*ve = append(*ve, ValidationError{Field: field, Err: errValidationMin})
			}
		case "max":
			res, err := strconv.Atoi(subkeys[1])
			if err != nil {
				return errMaxKey
			} else if value > res {
				*ve = append(*ve, ValidationError{Field: field, Err: errValidationMax})
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
				*ve = append(*ve, ValidationError{Field: field, Err: errValidationIn})
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
				*ve = append(*ve, ValidationError{Field: field, Err: errValidationLen})
			}
		case "regexp":
			rx, err := regexp.Compile(subkeys[1])
			if err != nil {
				return errRegexpKey
			}
			res := rx.FindString(value)
			if res != value {
				*ve = append(*ve, ValidationError{Field: field, Err: errValidationRegexp})
			}
		case "in":
			if !strings.Contains(subkeys[1], value) {
				*ve = append(*ve, ValidationError{Field: field, Err: errValidationIn})
			}
		default:
			return errInvelidTag
		}
	}
	return nil
}

func validateSlice(value reflect.Value, field reflect.StructField, param string, ve *ValidationErrors) error {
	switch field.Type.Elem().Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for i := 0; i < value.Len(); i++ {
			err := validateInt(param, field.Name, int(value.Index(i).Int()), ve)
			if err != nil {
				return err
			}
		}
	case reflect.String:
		for i := 0; i < value.Len(); i++ {
			err := validateString(param, field.Name, value.Index(i).String(), ve)
			if err != nil {
				return err
			}
		}
	default:
		*ve = append(*ve, ValidationError{Field: field.Name, Err: errValidationUnknownSlice})
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
				if err := validateSlice(value, field, param, &ve); err != nil {
					return err
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
				ve = append(ve, ValidationError{Field: field.Name, Err: errValidationUnknownType})
			}
		}
	}
	if len(ve) > 0 {
		return fmt.Errorf(ve.Error())
	}
	return nil
}

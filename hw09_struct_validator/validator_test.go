package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   string          `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
	ExtraStruct struct {
		ArifmeticVals []int  `validate:"max:7"`          // should be ok
		NotInSet      string `validate:"in:bar,foo,D12"` // not in set
		NotAccordMin  int    `validate:"min:9"`          // will be 5
		NotAccordMax  int    `validate:"max:9"`          // will be 10
	}
	InvalidTags struct {
		BadRegexp string `validate:"len:5zc|regexp:a(*"` // bad regexp and incorrect len
		BadVal    int    `validate:"min:4rt|max:6u"`     // incorrect tag values
		BadTag    string `validate:"lan:10"`             // bad tag
		BadInTag  int    `validate:"in:7,5,g10"`
		BadMaxTag int    `validate:"maxf:15"`
	}
	UnknownTypes struct {
		wideCharacters []rune  `validate:"max:10"`
		radius         float64 `validate:"in:3.14,6.28,14.2"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name                     string
		in                       interface{}
		expectedErr              error
		expectedValidationErrors string
	}{
		{
			name:                     "int",
			in:                       10,
			expectedErr:              errInvalidType,
			expectedValidationErrors: "",
		},
		{
			name: "ExtraStruct",
			in: ExtraStruct{
				ArifmeticVals: []int{1, 2, 3, 4, 5},
				NotInSet:      "any",
				NotAccordMin:  5,
				NotAccordMax:  10,
			},
			expectedErr: errValidationError,
			expectedValidationErrors: "NotInSet-->invalid value, value is not in \"in\" list\n" +
				"NotAccordMin-->invalid value, ocured value is less than min=9\n" +
				"NotAccordMax-->invalid value, ocured value is greater than max=9\n",
		},
		{
			name: "InvalidTags",
			in: InvalidTags{
				BadRegexp: "any string",
				BadVal:    155,
				BadTag:    "any",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: errValidationError,
			expectedValidationErrors: "BadRegexp-->invalid key value, key len=5zc\n" +
				"BadRegexp-->invalid key value, key regexp = a(*\n" +
				"BadVal-->invalid key value, key min=4rt\n" +
				"BadVal-->invalid key value, key max=6u\n" +
				"BadTag-->invalid validate subtag\n" +
				"BadInTag-->invalid key value, key in=g10\n" +
				"BadInTag-->invalid value, value is not in \"in\" list\n" +
				"BadMaxTag-->invalid validate subtag\n",
		},
		{
			name: "Response",
			in: Response{
				Code: 200,
				Body: "Any body",
			},
			expectedErr:              nil,
			expectedValidationErrors: "",
		},
		{
			name: "User",
			in: User{
				ID:     "43257405nfyelfyr6493jfy936781130ik",
				Name:   "Any name",
				Age:    16,
				Email:  "somemaildomenyanix.ru",
				Role:   "admin",
				Phones: []string{"281-330-800", "1-800-nmber"},
			},
			expectedErr: errValidationError,
			expectedValidationErrors: "ID-->invalid value, value's length is greater than len=36\n" +
				"Age-->invalid value, ocured value is less than min=18\n" +
				"Email-->invalid value, value doesn't match regular expression\n",
		},
		{
			name: "Token",
			in: Token{
				Header:    []byte{'c', '1', 'b'},
				Payload:   []byte{'p', 'a', 'y'},
				Signature: []byte{'b', 'y', 't', 'e', 's'},
			},
			expectedErr:              nil,
			expectedValidationErrors: "",
		},
		{
			name: "UnknownTypes",
			in: UnknownTypes{
				wideCharacters: []rune{432, 123, 654},
				radius:         3.14,
			},
			expectedErr:              errValidationError,
			expectedValidationErrors: "wideCharacters-->unknown slice field type\nradius-->unknown field type\n",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d %v", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr, "Errors missmatch. Actual error is %v", err)
			require.Equal(t, tt.expectedValidationErrors, ve.Error(), "Errors missmatch. Actual errors are %v", ve.Error())
		})
	}
}

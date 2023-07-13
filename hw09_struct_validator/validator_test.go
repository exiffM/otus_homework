package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	Counter  int
	UserRole string
)

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
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
	InvalidLenTag struct {
		BadRegexp string `validate:"len:5zc|regexp:a(*"` // incorrect len
		BadVal    int    `validate:"min:4rt|max:6u"`     // incorrect tag values
		BadTag    string `validate:"lan:10"`             // bad tag
		BadInTag  int    `validate:"in:7,5,g10"`
		BadMaxTag int    `validate:"maxf:15"`
	}
	InvalidRegexpTag struct {
		BadRegexp string `validate:"regexp:a(*"`     // bad regexp
		BadVal    int    `validate:"min:4rt|max:6u"` // incorrect tag values
		BadTag    string `validate:"lan:10"`         // bad tag
		BadInTag  int    `validate:"in:7,5,g10"`
		BadMaxTag int    `validate:"maxf:15"`
	}
	InvalidMinTag struct {
		BadVal    int    `validate:"min:4rt|max:6u"` // incorrect tag values
		BadTag    string `validate:"lan:10"`         // bad tag
		BadInTag  int    `validate:"in:7,5,g10"`
		BadMaxTag int    `validate:"maxf:15"`
	}
	InvalidMaxTag struct {
		BadVal    int    `validate:"max:6u"` // incorrect tag values
		BadTag    string `validate:"lan:10"` // bad tag
		BadInTag  int    `validate:"in:7,5,g10"`
		BadMaxTag int    `validate:"maxf:15"`
	}
	InvalidTag struct {
		BadVal    int    `validate:"mex:6"` // bad tag
		BadTag    string `validate:"len:10"`
		BadInTag  int    `validate:"in:7,5,10"`
		BadMaxTag int    `validate:"max:15"`
	}
	NotIn struct {
		BadVal    int    `validate:"in:6,16,184"`
		BadTag    string `validate:"len:10"`
		BadInTag  int    `validate:"in:7,5,10"`
		BadMaxTag int    `validate:"max:15"`
	}
	InvalidIn struct {
		BadVal    int    `validate:"on:6,16,184"`
		BadTag    string `validate:"len:10"`
		BadInTag  int    `validate:"in:7,5,10"`
		BadMaxTag int    `validate:"max:15"`
	}
	UnknownTypes struct {
		FloatingPoint []float32 `validate:"max:10"`
		Radius        float64   `validate:"in:3.14,6.28,14.2"`
	}
	TypeAlias struct {
		WideCharacters []rune  `validate:"in:435,324,546"`
		Count          Counter `validate:"in:10,15,20"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr []error
	}{
		{
			name:        "int",
			in:          10,
			expectedErr: []error{errInvalidType},
		},
		{
			name: "ExtraStruct",
			in: ExtraStruct{
				ArifmeticVals: []int{1, 2, 3, 4, 5},
				NotInSet:      "any",
				NotAccordMin:  5,
				NotAccordMax:  10,
			},
			expectedErr: []error{errValidationIn, errValidationMin, errValidationMax},
		},
		{
			name: "InvalidLenTag",
			in: InvalidLenTag{
				BadRegexp: "any string",
				BadVal:    155,
				BadTag:    "any",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: []error{errLenKey},
		},
		{
			name: "InvalidRegexpTag",
			in: InvalidRegexpTag{
				BadRegexp: "any string",
				BadVal:    155,
				BadTag:    "any",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: []error{errRegexpKey},
		},
		{
			name: "InvalidMinTag",
			in: InvalidMinTag{
				BadVal:    155,
				BadTag:    "any",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: []error{errMinKey},
		},
		{
			name: "InvalidMaxTag",
			in: InvalidMaxTag{
				BadVal:    155,
				BadTag:    "any",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: []error{errMaxKey},
		},
		{
			name: "InvalidTag",
			in: InvalidTag{
				BadVal:    155,
				BadTag:    "any",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: []error{errInvelidTag},
		},
		{
			name: "NotIn",
			in: NotIn{
				BadVal:    155,
				BadTag:    "anylentens",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: []error{errValidationIn},
		},
		{
			name: "InvalidIn",
			in: InvalidIn{
				BadVal:    155,
				BadTag:    "anylentens",
				BadInTag:  10,
				BadMaxTag: 15,
			},
			expectedErr: []error{errInvelidTag},
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
			expectedErr: []error{errValidationLen, errValidationMin, errValidationRegexp},
		},
		{
			name: "Token",
			in: Token{
				Header:    []byte{'c', '1', 'b'},
				Payload:   []byte{'p', 'a', 'y'},
				Signature: []byte{'b', 'y', 't', 'e', 's'},
			},
			expectedErr: nil,
		},
		{
			name: "UnknownTypes",
			in: UnknownTypes{
				FloatingPoint: []float32{4.4, 12.4, 213.4},
				Radius:        3.14,
			},
			expectedErr: []error{errValidationUnknownSlice, errValidationUnknownType},
		},
		{
			name: "TypeAlias",
			in: TypeAlias{
				WideCharacters: []rune{435, 324, 546},
				Count:          15,
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d %v", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			// expErrStr := errors.Join(tt.expectedErr...)
			// if expErrStr != nil {
			// 	require.Equal(t, expErrStr.Error(), err.Error(), "Errorrs message mismatch. Actual error is %v", err.Error())
			// }
			for _, expErr := range tt.expectedErr {
				require.ErrorIs(t, err, expErr, "Errors missmatch. Actual error is %v", err)
			}
		})
	}
}

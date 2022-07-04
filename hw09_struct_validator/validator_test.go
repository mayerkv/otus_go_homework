package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "123",
				Name:   "John Doe",
				Age:    10,
				Email:  "bad_email",
				Role:   "user",
				Phones: []string{"899988877660", "89998887766"},
				meta:   nil,
			},
			ValidationErrors{
				{"ID", ErrStrLengthRuleIsInvalid},
				{"Age", ErrIntMinRuleIsInvalid},
				{"Email", ErrStrRegexpRuleIsInvalid},
				{"Role", ErrStrInRuleIsInvalid},
				{"Phones[0]", ErrStrLengthRuleIsInvalid},
			},
		},
		{
			User{
				ID:     "111111111100000000001111111111000000",
				Name:   "",
				Age:    18,
				Email:  "foo@bar.com",
				Role:   "admin",
				Phones: []string{"89998887766"},
				meta:   nil,
			},
			nil,
		},
		{
			[]string{},
			ErrOnlyStructValidationAllowed,
		},
		{
			App{"123"},
			ValidationErrors{
				{"Version", ErrStrLengthRuleIsInvalid},
			},
		},
		{
			App{"12345"},
			nil,
		},
		{
			Response{200, ""},
			nil,
		},
		{
			Response{201, ""},
			ValidationErrors{
				{"Code", ErrIntInRuleIsInvalid},
			},
		},
		{
			struct {
				A string `validate:"min:invalid_value"`
			}{},
			ErrIntMinRuleWrongFormat,
		},
		{
			struct {
				Items []string `validate:"len:11"`
			}{
				Items: nil,
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			errs := Validate(tt.in)

			if errors.Is(tt.expectedErr, nil) {
				require.NoError(t, errs)
				return
			}

			var expErrs ValidationErrors
			var vErrs ValidationErrors
			if errors.As(tt.expectedErr, &expErrs) {
				require.ErrorAs(t, errs, &vErrs)
				for idx, expErr := range expErrs {
					require.Equal(t, expErr.Field, vErrs[idx].Field)
					require.ErrorIs(t, expErr.Err, errors.Unwrap(vErrs[idx].Err))
				}
				return
			}

			require.ErrorIs(t, errs, tt.expectedErr)
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name string
		v    ValidationErrors
		want string
	}{
		{
			name: "when no errors, then returns empty string",
			v:    ValidationErrors{},
			want: "",
		},
		{
			name: "simple test",
			v: ValidationErrors{
				{"Id", fmt.Errorf("invalid id")},
				{"Name", fmt.Errorf("invalid name")},
			},
			want: "Id: invalid id\nName: invalid name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

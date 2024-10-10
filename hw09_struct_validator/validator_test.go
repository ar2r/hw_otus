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
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Digits struct {
		Values []int `validate:"min:0|max:9"`
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
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "mixed valid",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "hello@world.ru",
				Role:   "admin",
				Phones: []string{"12345678901", "12345678902"},
			},
			expectedErr: nil,
		},
		{
			name: "string len error",
			in: User{
				ID:    "too-short",
				Name:  "John",
				Age:   25,
				Email: "hello@world.ru",
				Role:  "admin",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("%w: expected 36, got 9", ErrStringLengthMismatch),
				},
			},
		},
		{
			name: "int min max error",
			in: User{
				ID:    "123456789012345678901234567890123456",
				Name:  "John",
				Age:   100500,
				Email: "hello@world.ru",
				Role:  "admin",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("%w: 50", ErrValueIsMoreThanMaxValue),
				},
			},
		},
		{
			name: "string slice error",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "hello@world.ru",
				Role:   "admin",
				Phones: []string{"123456789012"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("%w: expected 11, got 12", ErrStringLengthMismatch),
				},
			},
		},
		{
			name: "int slice valid",
			in: Digits{
				Values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			},
			expectedErr: nil,
		},
		{
			name: "int slice error",
			in: Digits{
				Values: []int{100},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Values",
					Err:   fmt.Errorf("%w: 9", ErrValueIsMoreThanMaxValue),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
			fmt.Println(err)
		})
	}
}

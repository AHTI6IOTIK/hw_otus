package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
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
		meta   json.RawMessage //nolint:nolintlint,unused
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
	t.Parallel()
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "236a2eda-990b-40e1-be94-2b8f1514e6a1",
				Name:   "",
				Age:    50,
				Email:  "sss@ss.ss",
				Role:   "stuff",
				Phones: []string{"11122233344"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "",
				Name:   "",
				Age:    50,
				Email:  "sss@ss.ss",
				Role:   "stuff",
				Phones: []string{"11122233344"},
				meta:   nil,
			},
			expectedErr: ErrLen,
		},
		{
			in: User{
				ID:     "236a2eda-990b-40e1-be94-2b8f1514e6a1",
				Name:   "",
				Age:    50,
				Email:  "",
				Role:   "stuff",
				Phones: []string{"11122233344"},
				meta:   nil,
			},
			expectedErr: ErrRegexp,
		},
		{
			in: User{
				ID:     "236a2eda-990b-40e1-be94-2b8f1514e6a1",
				Name:   "",
				Age:    50,
				Email:  "sss@ss.ss",
				Role:   "stuff",
				Phones: []string{"32"},
				meta:   nil,
			},
			expectedErr: ErrLen,
		},
		{
			in: User{
				ID:     "236a2eda-990b-40e1-be94-2b8f1514e6a1",
				Name:   "",
				Age:    1000,
				Email:  "sss@ss.ss",
				Role:   "stuff",
				Phones: []string{"11122233344"},
				meta:   nil,
			},
			expectedErr: ErrNumber,
		},
		{
			in: User{
				ID:     "236a2eda-990b-40e1-be94-2b8f1514e6a1",
				Name:   "",
				Age:    50,
				Email:  "sss@ss.ss",
				Role:   "test",
				Phones: []string{"11122233344"},
				meta:   nil,
			},
			expectedErr: ErrList,
		},
		{
			in: App{
				Version: "123",
			},
			expectedErr: ErrLen,
		},
		{
			in: App{
				Version: "",
			},
			expectedErr: ErrLen,
		},
		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 12,
			},
			expectedErr: ErrNumberList,
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 500,
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr != nil {
				cErr, ok := err.(*ValidationErrors) //nolint:errorlint,nolintlint
				if !ok {
					t.Errorf("actual = %v, is fail cast to ValidationErrors", err)
				}

				for _, validationError := range *cErr {
					if !errors.Is(validationError, tt.expectedErr) {
						t.Errorf("actual = %v, want = %v", err, tt.expectedErr)
					}
				}
			} else if err != nil {
				t.Errorf("the error was not expected, actual = %v", err)
			}
		})
	}
}

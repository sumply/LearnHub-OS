package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTeacher_CheckGroupsAllowed(t *testing.T) {
	var (
		ID1 = uuid.New()
		ID2 = uuid.New()
		ID3 = uuid.New()
		ID4 = uuid.New()
		ID5 = uuid.New()
		ID6 = uuid.New()
	)
	tests := []struct {
		name          string
		allowed       []uuid.UUID
		checked       []uuid.UUID
		expectedError bool
		idsInError    []uuid.UUID
	}{
		{
			name:          "checked is empty",
			allowed:       []uuid.UUID{ID1, ID2, ID3},
			checked:       nil,
			expectedError: false,
			idsInError:    nil,
		},
		{
			name:          "all checked are not allowed",
			allowed:       []uuid.UUID{ID1, ID2, ID3},
			checked:       []uuid.UUID{ID4, ID5, ID6},
			expectedError: true,
			idsInError:    []uuid.UUID{ID4, ID5, ID6},
		},
		{
			name:          "some checked are not allowed",
			allowed:       []uuid.UUID{ID1, ID2},
			checked:       []uuid.UUID{ID1, ID2, ID3, ID4},
			expectedError: true,
			idsInError:    []uuid.UUID{ID3, ID4},
		},
		{
			name:          "correct",
			allowed:       []uuid.UUID{ID1, ID2, ID3},
			checked:       []uuid.UUID{ID1, ID2, ID3},
			expectedError: false,
			idsInError:    nil,
		},
		{
			name:          "allowed is empty",
			allowed:       nil,
			checked:       []uuid.UUID{ID1, ID2},
			expectedError: true,
			idsInError:    []uuid.UUID{ID1, ID2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teacher := Teacher{
				Groups: tt.allowed,
			}
			err := teacher.CheckGroupsAllowed(tt.checked)
			if tt.expectedError {
				assert.Error(t, err)
				if err, ok := errors.AsType[*Error](err); ok {
					assert.Len(t, err.Data(), len(tt.idsInError))
					for _, d := range err.Data() {
						id, err := uuid.Parse(d.Field)
						assert.NoError(t, err)
						assert.Contains(t, tt.idsInError, id)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

/*
func TestNewUserRole(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		expectedReturn UserRole
		expectedError  bool
	}{
		{
			name:           "admin",
			param:          "admin",
			expectedReturn: RoleAdmin,
			expectedError:  false,
		},
		{
			name:           "teacher",
			param:          "teacher",
			expectedReturn: RoleTeacher,
			expectedError:  false,
		},
		{
			name:           "student",
			param:          "student",
			expectedReturn: RoleStudent,
			expectedError:  false,
		},
		{
			name:           "invalid",
			param:          "invalid",
			expectedReturn: RoleInvalid,
			expectedError:  true,
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			role, err := NewUserRole(tests[i].param)
			if tests[i].expectedError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tests[i].expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tests[i].expectedReturn != role {
				t.Errorf("uncorrect conversion; expect: %d; got: %d", tests[i].expectedReturn, role)
			}
		})
	}
}
*/

/*
func TestUserRoleString(t *testing.T) {
	tests := []struct {
		name           string
		role           UserRole
		expectedReturn string
	}{
		{
			name:           "admin",
			role:           RoleAdmin,
			expectedReturn: "admin",
		},
		{
			name:           "teacher",
			role:           RoleTeacher,
			expectedReturn: "teacher",
		},
		{
			name:           "student",
			role:           RoleStudent,
			expectedReturn: "student",
		},
		{
			name:           "invalid",
			role:           123312,
			expectedReturn: "invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := test.role.String()
			if test.expectedReturn != str {
				t.Errorf("invalid conversion; expected: %s; got: %s", test.expectedReturn, str)
			}
		})
	}
}
*/

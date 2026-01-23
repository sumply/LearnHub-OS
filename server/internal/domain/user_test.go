package domain

import "testing"

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

package domain

import (
	"strings"
	"testing"
)

// TESTING EMAIL

func newEmailTestInvalidRegexp(t *testing.T, email string) {
	_email, err := NewEmail(email)
	if err == nil {
		t.Errorf("expected error, but got email = %s", _email)
	}
}

func TestNewEmailInvalidRegexp1(t *testing.T) {
	newEmailTestInvalidRegexp(t, "invalid")
}

func TestNewEmailInvalidRegexp2(t *testing.T) {
	newEmailTestInvalidRegexp(t, "invalid@invalid")
}

func TestNewEmailInvalidRegexp3(t *testing.T) {
	newEmailTestInvalidRegexp(t, "@invalid")
}

func TestNewEmailInvalidRegexp4(t *testing.T) {
	newEmailTestInvalidRegexp(t, "invalid.invalid")
}

func TestNewEmailInvalidRegexp5(t *testing.T) {
	newEmailTestInvalidRegexp(t, "invalid@invalid.")
}

func TestNewEmailInvalidRegexp6(t *testing.T) {
	newEmailTestInvalidRegexp(t, "@.")
}

func TestNewEmailEmpty(t *testing.T) {
	email, err := NewEmail("")
	if err == nil {
		t.Errorf("expected error, but got email = %s", email)
	}
}

func TestNewEmailCorrect(t *testing.T) {
	_, err := NewEmail("test@test.test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TESTING CREDENTIAL

type CredentialMakerInline struct {
	ResultLogin   string
	ResultPwd     string
	ResultPwdHash string
}

func NewStubCredentialMaker() *CredentialMakerInline {
	return &CredentialMakerInline{
		ResultLogin:   "",
		ResultPwd:     "",
		ResultPwdHash: "",
	}
}

func (s *CredentialMakerInline) GenerateLogin() Login {
	return Login(s.ResultLogin)
}

func (s *CredentialMakerInline) GeneratePassword() Password {
	return Password(s.ResultPwd)
}

func (s *CredentialMakerInline) Hash(Password) PwdHash {
	return PwdHash(s.ResultPwdHash)
}

func TestNewCredentialInvalidEmail(t *testing.T) {
	cred, _, err := NewCredential("", NewStubCredentialMaker())
	if err == nil {
		t.Errorf("expected error, but got credential= %v", cred)
	}
}

func TestNewCredentialCheckMaker(t *testing.T) {
	const (
		LOGIN    = "login"
		PASSWORD = "pwd"
		HASH     = "hash"
		EMAIL    = "test@test.test"
	)
	maker := &CredentialMakerInline{
		ResultLogin:   LOGIN,
		ResultPwd:     PASSWORD,
		ResultPwdHash: HASH,
	}
	cred, pwd, err := NewCredential(EMAIL, maker)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cred.Email != EMAIL {
		t.Errorf("unexpected email; wanted: %s got: %s", EMAIL, cred.Email)
	}
	if cred.Login != LOGIN {
		t.Errorf("unexpected login; wanted: %s got: %s", LOGIN, cred.Login)
	}
	if pwd != PASSWORD {
		t.Errorf("unexpected password; wanted: %s got: %s", PASSWORD, pwd)
	}
	if cred.PwdHash != HASH {
		t.Errorf("unexpected password hash; wanted: %s got: %s", HASH, cred.PwdHash)
	}
}

// TEST USER NAME

func TestNewUserEmpty(t *testing.T) {
	_, err := newUserName(firstNamePiece, "")
	if err == nil {
		t.Error("expected error, but got nil")
	}
}

func TestNewUserShort(t *testing.T) {
	_, err := newUserName(firstNamePiece, "a")
	if err == nil {
		t.Error("expected error, but got nil")
	}
}

func TestNewUserLong(t *testing.T) {
	_, err := newUserName(firstNamePiece, strings.Repeat("a", 101))
	if err == nil {
		t.Error("expected error, but got nil")
	}
}

func TestNewUserWithAnotherChars1(t *testing.T) {
	_, err := newUserName("test", "asdfj23s")
	if err == nil {
		t.Error("expected error, but got nil")
	}
}

func TestNewUserValid(t *testing.T) {
	_, err := newUserName("test", "alex")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewUserFormatting1(t *testing.T) {
	newUserTestFormatting(t, " alex    ", "Alex")
}

func TestNewUserFormatting2(t *testing.T) {
	newUserTestFormatting(t, " алекс    ", "Алекс")
}

// TEST PROFILE

func TestNewProfileInvalidFirstName(t *testing.T) {
	_, err := NewProfile("", "valid", "valid", RoleNone, AccessUser, 0)
	if err == nil {
		t.Errorf("profile created with invalid first name")
	}
}
func TestNewProfileInvalidLastName(t *testing.T) {
	_, err := NewProfile("valid", "", "valid", RoleNone, AccessUser, 0)
	if err == nil {
		t.Errorf("profile created with invalid last name")
	}
}

func TestNewProfileEmptyMiddleName(t *testing.T) {
	_, err := NewProfile("valid", "valid", "", RoleNone, AccessUser, 0)
	if err != nil {
		t.Errorf("profile can not be created with empty middle name")
	}
}

func TestNewProfileInvalidRole(t *testing.T) {
	_, err := NewProfile("valid", "valid", "", 123, AccessUser, 0)
	if err == nil {
		t.Errorf("profile created with invalid role")
	}
}

func TestNewProfileValid(t *testing.T) {
	_, err := NewProfile("valid", "valid", "valid", RoleNone, AccessUser, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewProfileFields1(t *testing.T) {
	NewProfileWithTestFields(t, "Ivan", "Lomonosov", "Georgievich", RoleNone)
}

func TestNewProfileFields2(t *testing.T) {
	NewProfileWithTestFields(t, "Ivan", "Lomonosov", "Georgievich", RoleStudent)
}

func TestNewProfileFields3(t *testing.T) {
	NewProfileWithTestFields(t, "Ivan", "Lomonosov", "Georgievich", RoleTeacher)
}

func TestNewProfileFields4(t *testing.T) {
	NewProfileWithTestFields(t, "Ivan", "Lomonosov", "", RoleTeacher)
}

func NewProfileWithTestFields(t *testing.T, firstName, lastName, middleName string, role UserRole) {
	profile, err := NewProfile(firstName, lastName, middleName, role, AccessUser, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if profile.FirstName != UserName(firstName) {
		t.Errorf("first name field is incorrect; want: %s; got: %s", firstName, profile.FirstName)
	}
	if profile.LastName != UserName(lastName) {
		t.Errorf("first name field is incorrect; want: %s; got: %s", lastName, profile.LastName)
	}
	if profile.MiddleName != UserName(middleName) {
		t.Errorf("first name field is incorrect; want: %s; got: %s", middleName, profile.MiddleName)
	}
	if profile.Role != role {
		t.Errorf("role is incorrect; want: %d; got: %d", role, profile.Role)
	}
}

func newUserTestFormatting(t *testing.T, original, wanted string) {
	name, err := newUserName("test", original)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if name != UserName(wanted) {
		t.Errorf("expected formatting name; wanted: %s; got: %s", wanted, name)
	}
}

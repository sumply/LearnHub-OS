package auth

type MockToken struct {
	Subject int64  `json:"subject"`
	Role    string `json:"role"`
}

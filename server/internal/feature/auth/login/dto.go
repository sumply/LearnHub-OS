package login

type Input struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Output struct {
	JWT OutputJWT `json:"jwt"`
}

type OutputJWT struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

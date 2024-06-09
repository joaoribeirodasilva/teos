package requests

type ForgotPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Password      string `json:"password"`
	CheckPassword string `json:"checkPassword"`
}

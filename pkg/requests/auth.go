package requests

type SignUp struct {
	Username             string `json:"username" validate:"required,min=3,max=20"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6,max=20"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Email                string `json:"email" validate:"required,email"`
	Token                string `json:"token" validate:"required"`
	Password             string `json:"password" validate:"required,min=6,max=20"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

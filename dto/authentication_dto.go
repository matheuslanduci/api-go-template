package dto

type CreateSessionWithPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Remember bool   `json:"remember" validate:"boolean"`
}

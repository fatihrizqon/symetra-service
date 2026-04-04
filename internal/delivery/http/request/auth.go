package request

// LoginRequest — FIXED
//
// BUG: Email field had validate:"required,min=1,max=20" — a maximum of 20 characters
// is far too short for a valid email address. RFC 5321 allows up to 254 characters.
// Any email longer than 20 chars (e.g. "firstname.lastname@company.com") would be
// rejected with a cryptic validation error.
//
// FIX: Changed max=20 to max=254 to match the RFC 5321 standard.
// Also added the `email` validator tag to enforce proper email format.
type LoginRequest struct {
	Email    string `validate:"required,email,min=1,max=254" json:"email" example:"johndoe@example.com"`
	Password string `validate:"required,min=8" json:"password" example:"yoursecretpassword"`
}

type RegisterRequest struct {
	Username string `validate:"required,min=1,max=20" json:"username"`
	Email    string `validate:"required,email,min=1,max=254" json:"email"`
	Password string `validate:"required,min=8" json:"password"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

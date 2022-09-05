package models

// SignUp struct to describe register a new user.
type SignUp struct {
	Name     string `json:"name" validate:"required,name,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
	UserRole string `json:"role" validate:"required,lte=25"`
}

// SignIn struct to describe login user.
type SignIn struct {
	Name     string `json:"name" validate:"required,name,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}

package domain

type UserInput struct {
	FirstName string `json:"firstName" xml:"firstName" bson:"firstName" validate:"required"`
	LastName  string `json:"lastName" xml:"lastName" bson:"lastName" validate:"required"`
	Email     string `json:"email" xml:"email" bson:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" xml:"password,omitempty" bson:"password" validate:"required"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequestEncrypt struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequestEncrypt struct {
	ID       string `bson:"_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
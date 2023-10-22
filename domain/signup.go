package domain

type UserInput struct {
	FirstName string `json:"firstName" xml:"firstName" bson:"firstName" validate:"required"`
	LastName  string `json:"lastName" xml:"lastName" bson:"lastName" validate:"required"`
	Email     string `json:"email" xml:"email" bson:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" xml:"password,omitempty" bson:"password" validate:"required"`
}


type SignupRequestEncrypt struct {
	Email    []byte `json:"email"`
	Password []byte `json:"password"`
}
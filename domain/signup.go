package domain

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequestEncrypt struct {
	Email    []byte `json:"email"`
	Password []byte `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequestEncrypt struct {
	Email    []byte `json:"email"`
	Password []byte `json:"password"`
}
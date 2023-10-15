package domain

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequestEncrypt struct {
	Email    []byte `json:"email"`
	Password []byte `json:"password"`
}
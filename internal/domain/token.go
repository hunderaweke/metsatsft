package domain

type Token struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type TokenRepository interface {
	CreateToken(token Token) error
	GetTokenByEmail(email string) (Token, error)
	DeleteToken(email string) error
}

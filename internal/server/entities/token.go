package entities

type Token struct {
	AccessToken  string
	RefreshToken string
}

func NewToken(at, rt string) *Token {
	return &Token{
		AccessToken:  at,
		RefreshToken: rt,
	}
}

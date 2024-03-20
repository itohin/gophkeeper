package entities

const (
	TextType = iota + 1
	PasswordType
)

type Secret struct {
	Name       string
	SecretType int64
	Notes      string
	Data       string
}

type Password struct {
	Login    string
	Password string
}

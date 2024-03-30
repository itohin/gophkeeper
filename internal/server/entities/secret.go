package entities

const (
	TypeText = iota + 1
	TypePassword
)

type Secret struct {
	ID         string
	Name       string
	SecretType uint32
	Notes      string
	Data       interface{}
	UserID     string
}

type Text struct {
	Text string
}

type Password struct {
	Login    string
	Password string
}

type SecretDTO struct {
	ID         string
	Name       string
	SecretType uint32
	Notes      string
	Data       []byte
	UserID     string
}

package entities

const (
	TypeText = iota + 1
	TypePassword
	TypeBinary
	TypeCard

	TextLabel     = "Текстовые данные"
	PasswordLabel = "Данные для входа(логин/пароль)"
	BinaryLabel   = "Бинарные данные"
	CardLabel     = "Данные банковских карт"
)

type Secret struct {
	ID         string
	Name       string
	SecretType uint32
	Notes      string
	Data       interface{}
}

type Password struct {
	Login    string
	Password string
}

type Text struct {
	Text string
}

type Binary struct {
	Binary []byte
}

type Card struct {
	Number     string
	Expiration string
	Code       string
	Pin        string
	OwnerName  string
}

func (s *Secret) GetLabel() string {
	switch s.SecretType {
	case TypeText:
		return TextLabel
	case TypePassword:
		return PasswordLabel
	case TypeBinary:
		return BinaryLabel
	case TypeCard:
		return CardLabel
	default:
		return ""
	}
}

package entities

const (
	TypeText = iota + 1
	TypePassword

	TextLabel     = "Текстовые данные"
	PasswordLabel = "Данные для входа(логин/пароль)"
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

func (s *Secret) GetLabel() string {
	switch s.SecretType {
	case TypeText:
		return TextLabel
	case TypePassword:
		return PasswordLabel
	default:
		return ""
	}
}

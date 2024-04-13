package storage

import (
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	pb "github.com/itohin/gophkeeper/proto"
)

type SecretsHydrator struct {
}

func NewSecretsHydrator() *SecretsHydrator {
	return &SecretsHydrator{}
}

func (h *SecretsHydrator) FromProto(v *pb.Secret) (*entities.Secret, error) {
	secret := &entities.Secret{
		ID:         v.Id,
		Name:       v.Name,
		SecretType: v.SecretType,
		Notes:      v.Notes,
	}
	switch d := v.Data.(type) {
	case *pb.Secret_Password:
		secret.Data = entities.Password{
			Login:    d.Password.Login,
			Password: d.Password.Password,
		}
	case *pb.Secret_Text:
		secret.Data = d.Text
	default:
		return nil, fmt.Errorf("unknown secret data type")
	}
	return secret, nil
}

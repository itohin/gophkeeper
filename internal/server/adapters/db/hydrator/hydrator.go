package hydrator

import (
	"encoding/json"
	"fmt"

	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/events"
	pb "github.com/itohin/gophkeeper/proto"
)

type SecretsHydrator struct {
}

func NewSecretsHydrator() *SecretsHydrator {
	return &SecretsHydrator{}
}

func (h *SecretsHydrator) FromProto(in *pb.Secret, userID string) (*entities.Secret, error) {
	secret := &entities.Secret{
		ID:         in.Id,
		Name:       in.Name,
		Notes:      in.Notes,
		SecretType: in.SecretType,
		UserID:     userID,
	}
	data, err := getProtoSecretData(in)
	if err != nil {
		return nil, fmt.Errorf("failed to convert proto data: %v", err)
	}
	secret.Data = data
	return secret, nil
}

func (h *SecretsHydrator) ToProto(in *events.SecretDTO) (*pb.Secret, error) {
	var t entities.Text
	var p entities.Password
	var b entities.Binary
	var c entities.Card
	secret := pb.Secret{
		Id:         in.ID,
		Name:       in.Name,
		SecretType: in.SecretType,
		Notes:      in.Notes,
	}
	switch in.SecretType {
	case entities.TypeText:
		err := json.Unmarshal(in.Data, &t)
		if err != nil {
			return nil, err
		}
		secret.Data = &pb.Secret_Text{Text: t.Text}
	case entities.TypePassword:
		err := json.Unmarshal(in.Data, &p)
		if err != nil {
			return nil, err
		}
		secret.Data = &pb.Secret_Password{
			Password: &pb.Password{Login: p.Login, Password: p.Password},
		}
	case entities.TypeCard:
		err := json.Unmarshal(in.Data, &c)
		if err != nil {
			return nil, err
		}
		secret.Data = &pb.Secret_Card{
			Card: &pb.Card{
				Number:     c.Number,
				Expiration: c.Expiration,
				Code:       c.Code,
				Pin:        c.Pin,
				OwnerName:  c.OwnerName,
			},
		}
	case entities.TypeBinary:
		err := json.Unmarshal(in.Data, &b)
		if err != nil {
			return nil, err
		}
		secret.Data = &pb.Secret_Binary{Binary: b.Binary}
	default:
		return nil, fmt.Errorf("unknown secret type")
	}
	return &secret, nil
}

func getProtoSecretData(in *pb.Secret) ([]byte, error) {
	switch d := in.Data.(type) {
	case *pb.Secret_Text:
		return json.Marshal(&entities.Text{
			Text: d.Text,
		})
	case *pb.Secret_Password:
		return json.Marshal(&entities.Password{
			Login:    d.Password.Login,
			Password: d.Password.Password,
		})
	case *pb.Secret_Card:
		return json.Marshal(&entities.Card{
			Number:     d.Card.Number,
			Expiration: d.Card.Expiration,
			Code:       d.Card.Code,
			Pin:        d.Card.Pin,
			OwnerName:  d.Card.OwnerName,
		})
	case *pb.Secret_Binary:
		return json.Marshal(&entities.Binary{
			Binary: d.Binary,
		})
	default:
		return nil, fmt.Errorf("unknown secret data type")
	}
}

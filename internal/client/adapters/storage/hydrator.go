package storage

import (
	"encoding/json"
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/pkg/events"
	pb "github.com/itohin/gophkeeper/proto"
)

type SecretsHydrator struct {
}

func NewSecretsHydrator() *SecretsHydrator {
	return &SecretsHydrator{}
}

func (h *SecretsHydrator) ToProto(s *entities.Secret) (*pb.Secret, error) {
	ps := &pb.Secret{
		Id:         s.ID,
		Name:       s.Name,
		SecretType: s.SecretType,
		Notes:      s.Notes,
	}
	switch d := s.Data.(type) {
	case *entities.Password:
		ps.Data = &pb.Secret_Password{
			Password: &pb.Password{
				Login:    d.Login,
				Password: d.Password,
			},
		}
	case string:
		ps.Data = &pb.Secret_Text{
			Text: d,
		}
	default:
		return nil, fmt.Errorf("unknown secret data type")
	}

	return ps, nil
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
		secret.Data = &entities.Password{
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

func (h *SecretsHydrator) FromSecretEvent(event *events.SecretEvent) (*entities.Secret, error) {
	var s entities.Secret
	s.ID = event.Secret.ID
	s.Name = event.Secret.Name
	s.SecretType = event.Secret.SecretType
	s.Notes = event.Secret.Notes

	var t entities.Text
	var p entities.Password
	switch s.SecretType {
	case entities.TypeText:
		err := json.Unmarshal(event.Secret.Data, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to transform secret event: %v", err)
		}
		s.Data = t.Text
	case entities.TypePassword:
		err := json.Unmarshal(event.Secret.Data, &p)
		if err != nil {
			return nil, fmt.Errorf("failed to transform secret event: %v", err)
		}
		s.Data = &entities.Password{
			Login:    p.Login,
			Password: p.Password,
		}
	default:
		return nil, fmt.Errorf("failed to transform secret event: unknown secret type")
	}

	return &s, nil
}

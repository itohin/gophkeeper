package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/events"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

type Secrets interface {
	Save(ctx context.Context, secret *entities.Secret) (*entities.Secret, error)
	GetUserSecrets(ctx context.Context, userID string) ([]events.SecretDTO, error)
	GetUserSecret(ctx context.Context, userID, secretID string) (events.SecretDTO, error)
}

type StreamConnection struct {
	userID   string
	deviceID string
	stream   pb.Secrets_CreateStreamServer
	error    chan error
}
type devicesMap map[string]*StreamConnection
type clientsMap map[string]devicesMap

type SecretsServer struct {
	pb.UnimplementedSecretsServer
	secrets       Secrets
	log           logger.Logger
	eventCh       chan *events.SecretEvent
	streamClients clientsMap
	mx            *sync.RWMutex
}

func (s *SecretsServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	userSecrets, err := s.secrets.GetUserSecrets(ctx, ctx.Value("user_id").(string))
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var secrets []*pb.Secret
	for _, v := range userSecrets {
		secret, err := s.buildSecret(&v)
		if err != nil {
			s.log.Error(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		secrets = append(secrets, secret)
	}

	return &pb.SearchResponse{
		Secrets: secrets,
	}, nil
}

func (s *SecretsServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	sDTO, err := s.secrets.GetUserSecret(ctx, ctx.Value("user_id").(string), in.Id)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	secret, err := s.buildSecret(&sDTO)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetResponse{
		Secret: secret,
	}, nil
}

func (s *SecretsServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	data, err := s.getData(in.Secret)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	savedSecret, err := s.secrets.Save(
		ctx,
		&entities.Secret{
			Name:       in.Secret.Name,
			Notes:      in.Secret.Notes,
			SecretType: in.Secret.SecretType,
			Data:       data,
			UserID:     ctx.Value("user_id").(string),
		},
	)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.CreateResponse{
		Id: savedSecret.ID,
	}, nil
}

func (s *SecretsServer) CreateStream(pConn *pb.StreamConnect, stream pb.Secrets_CreateStreamServer) error {
	log.Println("stream connect", pConn)

	ch := make(chan error)

	err := s.addClient(&StreamConnection{
		userID:   pConn.UserId,
		deviceID: pConn.FingerPrint,
		stream:   stream,
		error:    make(chan error),
	})
	if err != nil {
		return err
	}

	return <-ch
}

func (s *SecretsServer) Broadcast() {
	for {
		select {
		case event := <-s.eventCh:
			secret, err := s.buildSecret(&event.Secret)
			if err != nil {
				log.Printf("failed to build message: %v", err)
			}
			message := &pb.SecretEvent{
				Type:   pb.SecretEventType_EVENT_CREATED,
				Secret: secret,
			}

			devices := s.getClientDevices(event.Secret.UserID)
			log.Printf("devices: %v", devices)
			if devices != nil {
				for _, conn := range devices {
					err := conn.stream.Send(message)
					log.Printf("sendin message: %v\n", message)

					if err != nil {
						log.Printf("error with stream: %v - error: %v/n", conn.stream, err)
						conn.error <- err
					}
				}
			}
		default:
		}
	}
}

func (s *SecretsServer) getClientDevices(clientID string) devicesMap {
	s.mx.Lock()
	defer s.mx.Unlock()

	devices, ok := s.streamClients[clientID]
	if !ok {
		return nil
	}

	return devices
}

func (s *SecretsServer) addClient(conn *StreamConnection) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	devices, ok := s.streamClients[conn.userID]
	if !ok {
		devices = make(devicesMap)
	}
	if _, ok := devices[conn.deviceID]; ok {
		return fmt.Errorf("client id %s, deviceID %s already connected", conn.userID, conn.deviceID)
	}
	devices[conn.deviceID] = conn
	s.streamClients[conn.userID] = devices
	return nil
}

func (s *SecretsServer) removeClient(clientID, deviceID string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	devices, ok := s.streamClients[clientID]
	if !ok {
		return fmt.Errorf("client id %s, deviceID %s not found", clientID, deviceID)
	}

	delete(devices, deviceID)
	if len(s.streamClients[clientID]) == 0 {
		delete(s.streamClients, clientID)
	}

	return nil
}

func (s *SecretsServer) buildSecret(in *events.SecretDTO) (*pb.Secret, error) {
	var t entities.Text
	var p entities.Password
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
	default:
		return nil, fmt.Errorf("unknown secret type")
	}
	return &secret, nil
}

func (s *SecretsServer) getData(in *pb.Secret) ([]byte, error) {
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
	default:
		return nil, fmt.Errorf("unknown secret data type")
	}
}

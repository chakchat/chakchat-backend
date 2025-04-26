package notifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chakchat/chakchat-backend/notification-service/internal/grpc_service"
	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	reader      *kafka.Reader
	apnsClient  *APNsClient
	parser      *Parser
	grpcClients grpc_service.GRPCClients
	shutdown    chan struct{}
}

type Config struct {
}

func NewService(reader *kafka.Reader, apnsClient *APNsClient, grpcClients grpc_service.GRPCClients) *Service {

	return &Service{
		reader:      reader,
		apnsClient:  apnsClient,
		grpcClients: grpcClients,
		parser:      NewParser(&grpcClients),
		shutdown:    make(chan struct{}),
	}
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-s.shutdown:
				return
			case <-ctx.Done():
				return
			default:
				msg, err := s.reader.ReadMessage(ctx)
				if err != nil {
					if err == context.DeadlineExceeded {
						return
					}
					log.Printf("Error reading message: %v\n", err)
					time.Sleep(1 * time.Second)
					return
				}

				if err := s.processMessage(ctx, msg); err != nil {
					log.Printf("Error processing message: %v\n", err)
					return
				}
			}
		}
	}()
	log.Println("Notification service started")
}

func (s *Service) processMessage(ctx context.Context, msg kafka.Message) error {
	var notification struct {
		Receivers []uuid.UUID `json:"receivers"`
	}

	if err := json.Unmarshal(msg.Value, &notification); err != nil {
		return fmt.Errorf("failed to parse notification receivers: %w", err)
	}

	tokens := make(map[uuid.UUID]string)
	for _, id := range notification.Receivers {
		token, err := s.grpcClients.GetDeviceToken(ctx, id)
		if err != nil {
			if errors.Is(err, grpc_service.ErrNotFound) {
				continue
			}
			return fmt.Errorf("failed to get device tokens: %w", err)
		}
		if token != nil {
			tokens[id] = *token
		}
	}

	if len(tokens) == 0 {
		return nil
	}

	receiver, title, err := s.parser.ParseNotification(ctx, msg.Value)
	if err != nil {
		return fmt.Errorf("failed to parse notification: %w", err)
	}

	var lastErr error
	for id, token := range tokens {
		if err := s.apnsClient.SendNotification(token, receiver, title); err != nil {
			log.Printf("failed to send to user with id %s: %v", id, err)
			lastErr = err
		}
	}

	return lastErr
}

func (s *Service) Stop() error {
	close(s.shutdown)

	if err := s.reader.Close(); err != nil {
		return fmt.Errorf("failed to close kafka reader: %w", err)
	}

	log.Println("Notification service stopped")
	return nil
}

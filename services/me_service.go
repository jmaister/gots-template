package services

import (
	"context"

	"github.com/jmaister/gots-template/session"
	client "github.com/jmaister/taronja-gateway-clients/go"
)

type MeServiceInterface interface {
	GetMe() (*client.GetCurrentUserResponse, error)
}

type MeService struct {
	client     *client.ClientWithResponses
	adminToken string // Token for server-level operations (admin)
}

func NewMeService(taronjaClient *client.ClientWithResponses, adminToken string) *MeService {
	return &MeService{
		client:     taronjaClient,
		adminToken: adminToken,
	}
}

func (s *MeService) GetMe(ctx context.Context) (*client.GetCurrentUserResponse, error) {
	resp, err := s.client.GetCurrentUserWithResponse(
		ctx,
		session.ForwardAuthorizationCookie,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

package service

import (
	"context"
	"database/sql"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/auth"
	"github.com/fernando8franco/attengo/internal/config"
	"github.com/fernando8franco/attengo/internal/repository"
)

type RefreshTokenInput struct {
	Token string
}

type RefreshTokenService interface {
	CreateAccessToken(ctx context.Context, input RefreshTokenInput) (RefreshTokenResponse, error)
	RevokeRefreshToken(ctx context.Context, input RefreshTokenInput) error
}

type refreshTokenService struct {
	queries *repository.Queries
	cfg     *config.Config
}

func NewRefreshTokenService(db *sql.DB, cfg *config.Config) RefreshTokenService {
	return &refreshTokenService{
		queries: repository.New(db),
		cfg:     cfg,
	}
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *refreshTokenService) CreateAccessToken(ctx context.Context, input RefreshTokenInput) (RefreshTokenResponse, error) {
	userId, err := s.queries.GetUserIdFromRefreshToken(ctx, input.Token)
	if err != nil {
		return RefreshTokenResponse{}, apperr.NewUnauthorizedRequest(err.Error())
	}

	accessToken, err := auth.MakeJWT(s.cfg.IssuerJWT, s.cfg.SecretJWT, userId, s.cfg.ExpirationTime)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	return RefreshTokenResponse{AccessToken: accessToken}, nil
}

func (s *refreshTokenService) RevokeRefreshToken(ctx context.Context, input RefreshTokenInput) error {
	err := s.queries.SetRevokedAt(ctx, input.Token)
	if err != nil {
		return err
	}

	return nil
}

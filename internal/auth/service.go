package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ijalandhika/tuto-api/internal/auth/db"
	"github.com/ijalandhika/tuto-api/pkg/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries *db.Queries
	cfg     config.JWTConfig
}

func NewService(queries *db.Queries, cfg config.JWTConfig) *Service {
	return &Service{
		queries: queries,
		cfg:     cfg,
	}
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (*AuthResponse, error) {
	hash := hashToken(req.RefreshToken)

	session, err := s.queries.GetSessionByTokenHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("Invalid or expired refresh token")
	}

	if err := s.queries.RevokeSession(ctx, hash); err != nil {
		return nil, fmt.Errorf("Revoke session: %w", err)
	}

	actorId, err := uuid.FromBytes(session.ActorID.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("parse actor ID: %w", err)
	}

	return s.generateTokens(ctx, actorId.String(), string(session.ActorType))
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	parent, err := s.queries.GetParentByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(parent.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}

	parentId, err := uuid.FromBytes(parent.ID.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("parse parent ID: %w", err)
	}

	return s.generateTokens(ctx, parentId.String(), string(db.ActorTypeParent))
}

func (s *Service) Logout(ctx context.Context, req LogoutRequest) error {
	hash := hashToken(req.RefreshToken)

	if err := s.queries.RevokeSession(ctx, hash); err != nil {
		return fmt.Errorf("Revoke session: %w", err)
	}

	return nil
}

func (s *Service) Signup(ctx context.Context, req SignupRequest) (*AuthResponse, error) {
	// check if email already taken
	_, err := s.queries.GetParentByEmail(ctx, req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already registered")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("check email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	parent, err := s.queries.CreateParent(ctx, db.CreateParentParams{
		Email:          req.Email,
		PasswordHash:   string(hash),
		DisplayName:    req.DisplayName,
		MarketingOptIn: req.MarketingOptIn,
	})
	if err != nil {
		return nil, fmt.Errorf("create parent: %w", err)
	}

	parentID, err := uuid.FromBytes(parent.ID.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("parse parent ID: %w", err)
	}

	return s.generateTokens(ctx, parentID.String(), string(db.ActorTypeParent))

}

func (s *Service) generateTokens(ctx context.Context, actorID, actorType string) (*AuthResponse, error) {
	accessExpiry := time.Now().Add(time.Duration(s.cfg.AccessExpiryMinutes) * time.Minute)
	accessToken, err := s.mintJWT(actorID, actorType, accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("mint access token: %w", err)
	}

	rawRefresh, err := generateRandomToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshExpiry := time.Now().Add(time.Duration(s.cfg.RefreshExpiryDays) * 24 * time.Hour)

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, fmt.Errorf("parse actor id: %w", err)
	}

	// Store refresh token hash in DB
	_, err = s.queries.CreateSession(ctx, db.CreateSessionParams{
		ActorType: db.ActorType(actorType),
		ActorID:   pgtype.UUID{Bytes: actorUUID, Valid: true},
		TokenHash: hashToken(rawRefresh),
		ExpiresAt: pgtype.Timestamptz{Time: refreshExpiry, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

func (s *Service) mintJWT(actorID, actorType string, expiry time.Time) (string, error) {
	claims := jwt.MapClaims{
		"actor_id":   actorID,
		"actor_type": actorType,
		"exp":        expiry.Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

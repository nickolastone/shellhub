package services

import (
	"context"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/api/request"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type SessionService interface {
	ListSessions(ctx context.Context, pagination paginator.Query) ([]models.Session, int, error)
	GetSession(ctx context.Context, uid models.UID) (*models.Session, error)
	CreateSession(ctx context.Context, session request.SessionCreate) (*models.Session, error)
	DeactivateSession(ctx context.Context, uid models.UID) error
	KeepAliveSession(ctx context.Context, uid models.UID) error
	SetSessionAuthenticated(ctx context.Context, uid models.UID, authenticated bool) error
}

func (s *service) ListSessions(ctx context.Context, pagination paginator.Query) ([]models.Session, int, error) {
	return s.store.SessionList(ctx, pagination)
}

func (s *service) GetSession(ctx context.Context, uid models.UID) (*models.Session, error) {
	return s.store.SessionGet(ctx, uid)
}

func (s *service) CreateSession(ctx context.Context, session request.SessionCreate) (*models.Session, error) {
	model := models.Session{
		UID:       session.UID,
		DeviceUID: models.UID(session.DeviceUID),
		Username:  session.Username,
		IPAddress: session.IPAddress,
		Type:      session.Type,
		Term:      session.Term,
	}

	return s.store.SessionCreate(ctx, model)
}

func (s *service) DeactivateSession(ctx context.Context, uid models.UID) error {
	err := s.store.SessionDeleteActives(ctx, uid)
	if err == store.ErrNoDocuments {
		return NewErrSessionNotFound(uid, err)
	}

	return err
}

func (s *service) KeepAliveSession(ctx context.Context, uid models.UID) error {
	return s.store.SessionSetLastSeen(ctx, uid)
}

func (s *service) SetSessionAuthenticated(ctx context.Context, uid models.UID, authenticated bool) error {
	return s.store.SessionSetAuthenticated(ctx, uid, authenticated)
}

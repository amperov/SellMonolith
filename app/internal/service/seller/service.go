package seller

import (
	"Selling/app/pkg/auth"
	"context"
	"github.com/sirupsen/logrus"
)

type SellerStore interface {
	SignUp(ctx context.Context, m map[string]interface{}) (int, error)
	SignIn(ctx context.Context, m map[string]interface{}) (int, error)
	UpdateData(ctx context.Context, m map[string]interface{}, UserID int) error
}
type SellerService struct {
	tm auth.TokenManager
	s  SellerStore
}

func NewSellerService(tm auth.TokenManager, s SellerStore) *SellerService {
	return &SellerService{tm: tm, s: s}
}

func (s *SellerService) UpdateData(ctx context.Context, m map[string]interface{}, UserID int) error {
	err := s.s.UpdateData(ctx, m, UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SellerService) SignUp(ctx context.Context, m map[string]interface{}) (int, error) {
	id, err := s.s.SignUp(ctx, m)
	if err != nil {
		logrus.Println(err)
		return 0, err
	}
	return id, nil
}

func (s *SellerService) SignIn(ctx context.Context, m map[string]interface{}) (string, error) {
	id, err := s.s.SignIn(ctx, m)
	if err != nil {
		return "", err
	}

	token, err := s.tm.GenerateToken(id)
	if err != nil {
		return "", err
	}
	return token, nil
}

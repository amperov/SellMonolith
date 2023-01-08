package history

import (
	"context"
)

type HistoryStore interface {
	GetAll(ctx context.Context, UserID int) (map[string]interface{}, error)
	GetOne(ctx context.Context, UserID, TransactID int) (map[string]interface{}, error)
}
type HistoryService struct {
	h HistoryStore
}

func NewHistoryService(h HistoryStore) *HistoryService {
	return &HistoryService{h: h}
}

func (h *HistoryService) GetAllTransactions(ctx context.Context, UserID int) (map[string]interface{}, error) {
	all, err := h.h.GetAll(ctx, UserID)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (h *HistoryService) GetOneTransaction(ctx context.Context, UserID, TransactID int) (map[string]interface{}, error) {
	one, err := h.h.GetOne(ctx, UserID, TransactID)
	if err != nil {
		return nil, err
	}
	return one, nil
}

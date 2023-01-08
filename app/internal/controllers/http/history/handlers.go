package history

import (
	"Selling/app/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type HistoryService interface {
	GetAllTransactions(ctx context.Context, UserID int) (map[string]interface{}, error)
	GetOneTransaction(ctx context.Context, UserID, TransactID int) (map[string]interface{}, error)
}

type HistoryHandler struct {
	ware auth.MiddleWare
	hs   HistoryService
}

func NewHistoryHandler(ware auth.MiddleWare, hs HistoryService) *HistoryHandler {
	return &HistoryHandler{ware: ware, hs: hs}
}

func (h *HistoryHandler) Register(r *httprouter.Router) {
	r.GET("/seller/history", h.ware.IsAuth(h.GetHistory))
	r.GET("/seller/history/:tran_id", h.ware.IsAuth(h.GetFullTransaction))
}

func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	transactions, err := h.hs.GetAllTransactions(r.Context(), UserID)
	if err != nil {
		return
	}

	transactionsMarshalled, err := json.Marshal(transactions)
	if err != nil {
		return
	}

	_, err = w.Write(transactionsMarshalled)
	if err != nil {
		return
	}
}

func (h *HistoryHandler) GetFullTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	tran_id := params.ByName("tran_id")
	tranID, err := strconv.Atoi(tran_id)
	if err != nil {
		return
	}

	transaction, err := h.hs.GetOneTransaction(r.Context(), UserID, tranID)
	if err != nil {
		return
	}

	marshal, err := json.Marshal(transaction)
	if err != nil {
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}

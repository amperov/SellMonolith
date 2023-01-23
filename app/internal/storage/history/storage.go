package history

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var table = "transactions"

type HistoryStorage struct {
	c *pgxpool.Pool
}

func NewHistoryStorage(c *pgxpool.Pool) *HistoryStorage {
	return &HistoryStorage{c: c}
}

func (h *HistoryStorage) GetAll(ctx context.Context, UserID int) (map[string]interface{}, error) {
	var transacts AllTransactInput
	m := make(map[string]interface{})
	query, args, err := squirrel.
		Select("id", "category_name", "subcategory_name", "unique_code",
			"content_key", "state", "amount", "date_check", "client_email").
		From(table).Where(squirrel.Eq{"user_id": UserID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := h.c.Query(ctx, query, args...)
	if err != nil {
		logrus.Debugf("error Query: %v", err)
		return nil, err
	}

	var arrayMap []map[string]interface{}

	for rows.Next() {
		err := rows.Scan(&transacts.ID, &transacts.Category, &transacts.Subcategory, &transacts.UniqueCode,
			&transacts.Content, &transacts.State, &transacts.AmountUSD, &transacts.DateCheck, &transacts.CLientEmail)
		if err != nil {
			logrus.Debugf("error scanning: %v", err)
			return nil, err
		}
		arrayMap = append(arrayMap, transacts.ToMap())
	}
	m["transactions"] = arrayMap

	return m, nil
}

// GetOne TODO Need testing
func (h *HistoryStorage) GetOne(ctx context.Context, UserID, TransactID int) (map[string]interface{}, error) {
	var transtact Transaction
	query, args, err := squirrel.
		Select("category_name", "subcategory_name", "unique_code",
			"client_email", "amount", "profit", "count",
			"unique_inv", "date_delivery", "date_confirmed",
			"content_key", "state", "amount_usd", "date_check").
		From(table).Where(squirrel.Eq{"user_id": UserID, "id": TransactID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	row := h.c.QueryRow(ctx, query, args...)

	err = row.Scan(&transtact.Category, &transtact.Subcategory, &transtact.UniqueCode.UniqueCode,
		&transtact.ClientEmail, &transtact.Amount, &transtact.Profit, &transtact.CountGoods,
		&transtact.UniqueInv, &transtact.UniqueCode.DateDelivery, &transtact.UniqueCode.DateConfirmed,
		&transtact.Content, &transtact.UniqueCode.State, &transtact.AmountUSD, &transtact.UniqueCode.DateCheck)
	if err != nil {
		logrus.Debugf("error scanning: %v", err)
		return nil, err
	}
	transtact.UserID = UserID
	transtact.ID = TransactID
	var mm = make(map[string]interface{})
	mm["transaction"] = transtact.ToMap()
	return mm, nil
}
func (h *HistoryStorage) GetOneByUniqueCode(ctx context.Context, UniqueCode string) (map[string]interface{}, error) {
	var transtact Transaction
	m := make(map[string]interface{})
	query, args, err := squirrel.
		Select("id", "category_name", "subcategory_name",
			"client_email", "amount", "profit", "count",
			"unique_inv", "date_delivery", "date_confirmed",
			"content_key", "state", "amount_usd", "date_check").
		From(table).Where(squirrel.Eq{"unique_code": UniqueCode}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := h.c.QueryRow(ctx, query, args...)

	err = row.Scan(&transtact.ID, &transtact.Category, &transtact.Subcategory, &transtact.UniqueCode.UniqueCode,
		&transtact.ClientEmail, &transtact.Amount, &transtact.Profit, &transtact.CountGoods,
		&transtact.UniqueInv, &transtact.UniqueCode.DateDelivery, &transtact.UniqueCode.DateConfirmed,
		&transtact.Content, &transtact.UniqueCode.State, &transtact.AmountUSD, &transtact.UniqueCode.DateCheck)
	if err != nil {
		logrus.Debugf("error scanning: %v", err)
		return nil, err
	}
	var mm = make(map[string]interface{})
	mm["transaction"] = m
	return mm, nil
}

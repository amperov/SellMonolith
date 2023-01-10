package seller

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	sellerTable      = "sellers"
	transactionTable = "transactions"
)

type SellerStorage struct {
	c *pgxpool.Pool
}

func NewSellerStorage(c *pgxpool.Pool) *SellerStorage {
	return &SellerStorage{c: c}
}

func (s *SellerStorage) UpdateData(ctx context.Context, m map[string]interface{}, SellerID int) error {
	query, args, err := squirrel.Update(sellerTable).SetMap(m).PlaceholderFormat(squirrel.Dollar).Where(squirrel.Eq{"id": SellerID}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil

}

func (s *SellerStorage) GetSellerInfo(ctx context.Context, Username string) (int, string, error) {
	var SellerID int
	var SellerKey string
	query, args, err := squirrel.Select("seller_id", "seller_key").PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"username": Username}).From(sellerTable).ToSql()
	if err != nil {
		return 0, "", err
	}

	row := s.c.QueryRow(ctx, query, args...)

	err = row.Scan(&SellerID, &SellerKey)
	if err != nil {
		return 0, "", err
	}

	return SellerID, SellerKey, nil
}

func (s *SellerStorage) SetTransaction(ctx context.Context, model map[string]interface{}) error {
	query, args, err := squirrel.Insert(transactionTable).SetMap(model).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = s.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *SellerStorage) SignUp(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Insert(sellerTable).SetMap(m).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING ID").ToSql()
	if err != nil {
		return 0, err
	}

	row := s.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SellerStorage) SignIn(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Select("id").Where(squirrel.Eq{"username": m["username"], "pass": m["password"]}).PlaceholderFormat(squirrel.Dollar).From(sellerTable).ToSql()
	if err != nil {
		return 0, err
	}

	row := s.c.QueryRow(ctx, query, args...)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SellerStorage) GetInfo(ctx context.Context, UserID int) (map[string]interface{}, error) {
	var user Seller

	query, args, err := squirrel.Select("username", "firstname", "lastname", "seller_id", "seller_key").From(sellerTable).
		Where(squirrel.Eq{"id": UserID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := s.c.QueryRow(ctx, query, args...)
	err = row.Scan(&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.SellerID,
		&user.SellerKey)
	if err != nil {
		return nil, err
	}
	return user.ToMap(), nil
}

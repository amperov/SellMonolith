package client

import (
	"context"
)

type ClientDigi interface {
	Auth(ctx context.Context, Username string) string
	GetProducts(ctx context.Context, UniqueCode, Token string) ([]map[string]interface{}, error)
}

type ProductStore interface {
	Check(ctx context.Context, ItemID int) (bool, error)
	SearchByUniqueCode(ctx context.Context, UniqueCode string) ([]map[string]interface{}, bool)
}

type ClientService struct {
	c ClientDigi
	p ProductStore
}

func (c *ClientService) Check(ctx context.Context, ItemID int) (bool, error) {
	return c.p.Check(ctx, ItemID)
}

func NewClientService(c ClientDigi, p ProductStore) *ClientService {
	return &ClientService{c: c, p: p}
}

func (c *ClientService) Get(ctx context.Context, UniqueCode string, Username string) ([]map[string]interface{}, error) {
	prods, ok := c.p.SearchByUniqueCode(ctx, UniqueCode)
	if ok == false {
		token := c.c.Auth(ctx, Username)
		get, err := c.c.GetProducts(ctx, UniqueCode, token)

		if err != nil {
			return nil, err
		}
		return get, nil
	}

	return prods, nil
}

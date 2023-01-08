package category

import (
	"context"
)

type CategoryStore interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, CatID int) (int, error)
	Delete(ctx context.Context, CatID int) error
	GetOne(ctx context.Context, CatID int) (map[string]interface{}, error)
	GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error)
}
type CategoryService struct {
	c CategoryStore
}

func NewCategoryService(c CategoryStore) *CategoryService {
	return &CategoryService{c: c}
}

func (c *CategoryService) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	id, err := c.c.Create(ctx, m)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CategoryService) Update(ctx context.Context, m map[string]interface{}, UserID int, CatID int) (int, error) {
	id, err := c.c.Update(ctx, m, CatID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CategoryService) Delete(ctx context.Context, UserID int, CatID int) error {
	err := c.c.Delete(ctx, CatID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryService) GetOne(ctx context.Context, UserID int, CatID int) (map[string]interface{}, error) {
	m, err := c.c.GetOne(ctx, CatID)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CategoryService) GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error) {
	all, err := c.c.GetAll(ctx, UserID)
	if err != nil {
		return nil, err
	}
	return all, nil
}

package product

import "context"

type ProductStore interface {
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, ProdID int) (int, error)
	Delete(ctx context.Context, SubCatID, ProdID int) error
	GetAll(ctx context.Context, SubCatID int) ([]map[string]interface{}, error)
	GetCount(ctx context.Context, SubCatID int) (int, error)
}
type ProductService struct {
	p ProductStore
}

func NewProductService(p ProductStore) *ProductService {
	return &ProductService{p: p}
}

func (p *ProductService) Create(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error) {
	create, err := p.p.Create(ctx, m)
	if err != nil {
		return 0, err
	}

	return create, nil
}

func (p *ProductService) Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID, ProdID int) (int, error) {
	update, err := p.p.Update(ctx, m, ProdID)
	if err != nil {
		return 0, err
	}
	return update, nil
}

func (p *ProductService) Delete(ctx context.Context, UserID, CatID, SubCatID, ProdID int) error {
	err := p.p.Delete(ctx, SubCatID, ProdID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductService) GetAll(ctx context.Context, UserID, CatID, SubCatID int) ([]map[string]interface{}, error) {
	all, err := p.p.GetAll(ctx, SubCatID)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (p *ProductService) GetCount(ctx context.Context, UserID, CatID, SubCatID int) (int, error) {
	count, err := p.p.GetCount(ctx, SubCatID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

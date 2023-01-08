package subcategory

import "context"

type SubcatStore interface {
	GetAll(ctx context.Context, CatID int) ([]map[string]interface{}, error)
	GetCount(ctx context.Context, CatID int) (int, error)
	Create(ctx context.Context, m map[string]interface{}) (int, error)
	Update(ctx context.Context, m map[string]interface{}, SubCatID int) (int, error)
	Delete(ctx context.Context, SubCatID int) error
	Get(ctx context.Context, SubCatID int) (map[string]interface{}, error)
}
type SubcategoryService struct {
	s SubcatStore
}

func NewSubcategoryService(s SubcatStore) *SubcategoryService {
	return &SubcategoryService{s: s}
}

func (s *SubcategoryService) GetAll(ctx context.Context, UserID, CatID int) ([]map[string]interface{}, error) {
	all, err := s.s.GetAll(ctx, CatID)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (s *SubcategoryService) GetCount(ctx context.Context, UserID, CatID int) (int, error) {
	count, err := s.s.GetCount(ctx, CatID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SubcategoryService) Create(ctx context.Context, m map[string]interface{}, UserID, CatID int) (int, error) {
	create, err := s.s.Create(ctx, m)
	if err != nil {
		return 0, err
	}
	return create, nil
}
func (s *SubcategoryService) Update(ctx context.Context, m map[string]interface{}, UserID, CatID, SubCatID int) (int, error) {
	update, err := s.s.Update(ctx, m, SubCatID)
	if err != nil {
		return 0, err
	}
	return update, nil
}

func (s *SubcategoryService) Delete(ctx context.Context, UserID, CatID, SubCatID int) error {
	err := s.s.Delete(ctx, SubCatID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubcategoryService) Get(ctx context.Context, UserID, CatID, SubCatID int) (map[string]interface{}, error) {
	get, err := s.s.Get(ctx, SubCatID)
	if err != nil {
		return nil, err
	}
	return get, nil
}

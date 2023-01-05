package entity

type SubcategoryModel struct {
	ID         int
	Title      string
	CategoryID int
}

func (s *SubcategoryModel) ToMap() map[string]interface{} {
	var ModelMap map[string]interface{}

	ModelMap["title"] = s.Title
	ModelMap["category_id"] = s.CategoryID

	return ModelMap
}

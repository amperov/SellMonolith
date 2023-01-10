package subcategory

type Subcategory struct {
	ID         int
	Title      string
	CategoryID int
}

func (s *Subcategory) ToMap() map[string]interface{} {
	ModelMap := make(map[string]interface{})
	ModelMap["id"] = s.ID
	ModelMap["title"] = s.Title
	ModelMap["category_id"] = s.CategoryID

	return ModelMap
}

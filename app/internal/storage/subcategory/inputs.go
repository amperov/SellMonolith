package subcategory

type Subcategory struct {
	ID         int
	Title      string
	CategoryID int
}

func (s *Subcategory) ToMap() map[string]interface{} {
	var ModelMap map[string]interface{}

	ModelMap["title"] = s.Title
	ModelMap["category_id"] = s.CategoryID

	return ModelMap
}

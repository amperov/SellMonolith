package category

type Category struct {
	ID          int
	TitleRu     string
	TitleEng    string
	Description string
	UserID      int
}

func (m *Category) ToMap() map[string]interface{} {
	ModelMap := make(map[string]interface{})
	ModelMap["id"] = m.ID
	ModelMap["title_ru"] = m.TitleRu
	ModelMap["title_eng"] = m.TitleEng
	ModelMap["description"] = m.Description
	ModelMap["user_id"] = m.UserID

	return ModelMap
}

package entity

type CategoryModel struct {
	ID          int
	TitleRu     string
	TitleEng    string
	Description string
	UserID      int
}

func (m *CategoryModel) ToMap() map[string]interface{} {
	var ModelMap map[string]interface{}
	ModelMap["title_ru"] = m.TitleRu
	ModelMap["title_eng"] = m.TitleEng
	ModelMap["description"] = m.Description
	ModelMap["user_id"] = m.UserID

	return ModelMap
}

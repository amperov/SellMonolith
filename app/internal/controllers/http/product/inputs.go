package product

type CreateProductInput struct {
	Content  string `json:"content,omitempty"`
	SubCatID int    `json:"subcategory_id,omitempty"`
}

func (c *CreateProductInput) ToMap() map[string]interface{} {
	var m map[string]interface{}
	m["content"] = c.Content
	return m
}

type UpdateProductInput struct {
	Content  string `json:"content,omitempty"`
	SubCatID int    `json:"subcategory_id,omitempty"`
}

func (c *UpdateProductInput) ToMap() map[string]interface{} {
	var m map[string]interface{}
	if c.Content != "" {
		m["content"] = c.Content
		return m
	}
	return nil
}

type CreateProductsInput struct {
	products []CreateProductInput `json:"products,omitempty"`
}

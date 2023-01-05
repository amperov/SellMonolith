package entity

type SellerModel struct {
	ID        int
	Username  string
	Firstname string
	Lastname  string
	Password  string
	SellerID  int
	SellerKey string
}

func (m *SellerModel) ToMap() map[string]interface{} {

	var SellerMap map[string]interface{}

	SellerMap["username"] = m.Username
	SellerMap["firstname"] = m.Firstname
	SellerMap["lastname"] = m.Lastname

	if m.Password != "" {
		SellerMap["password"] = m.Password
	}
	if m.SellerID != 0 {
		SellerMap["seller_id"] = m.SellerID
	}
	if m.Password != "" {
		SellerMap["seller_key"] = m.SellerKey
	}
	return SellerMap
}

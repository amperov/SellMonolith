package seller

type SignUpInput struct {
	Username  string `json:"username,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Password  string `json:"password,omitempty"`
	SellerID  int    `json:"seller_id,omitempty"`
	SellerKey string `json:"seller_key,omitempty"`
}

func (m *SignUpInput) ToMap() map[string]interface{} {

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

type SignInInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (m *SignInInput) ToMap() map[string]interface{} {

	var SellerMap map[string]interface{}

	SellerMap["username"] = m.Username

	if m.Password != "" {
		SellerMap["password"] = m.Password
	}
	return SellerMap
}

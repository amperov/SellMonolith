package digiseller

type DigiInput struct {
	Retval          int     `json:"retval"`
	Retdesc         string  `json:"retdesc"`
	Inv             int     `json:"inv"`
	IdGoods         int     `json:"id_goods"`
	Amount          float64 `json:"amount"`
	TypeCurr        string  `json:"type_curr"`
	Profit          string  `json:"profit"`
	AmountUsd       float64 `json:"amount_usd"`
	DatePay         string  `json:"date_pay"`
	Email           string  `json:"email"`
	AgentId         int     `json:"agent_id"`
	AgentPercent    int     `json:"agent_percent"`
	UnitGoods       int     `json:"unit_goods"`
	CntGoods        int     `json:"cnt_goods"`
	PromoCode       string  `json:"promo_code"`
	BonusCode       string  `json:"bonus_code"`
	CartUid         string  `json:"cart_uid"`
	UniqueCodeState struct {
		State         int    `json:"state"`
		DateCheck     string `json:"date_check"`
		DateDelivery  string `json:"date_delivery"`
		DateConfirmed string `json:"date_confirmed"`
		DateRefuted   string `json:"date_refuted"`
	} `json:"unique_code_state"`
	Options []struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"options"`
}

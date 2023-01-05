package entity

type TransactionModel struct {
	ID            int
	UniqueCode    string
	ClientEmail   string
	DateCheck     string
	Content       string
	SubcategoryID int
}

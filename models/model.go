package models

import "github.com/shopspring/decimal"

//INSERT INTO USERS(user_name, password, balance, activity, bank_country_code, bank_name) VALUES ('test1', 'pass', '10000',0, 'ru', 'sberbank');
type User struct {
	UserName        string `json:"user_name" db:"user_name"`
	Password        string `json:"password" db:"password"`
	Balance         string `json:"balance" db:"balance"`
	BankName        string `json:"bank_name" db:"bank_name"`
	BankCountryCode string `json:"bank_country_code" db:"bank_country_code"`
	Activity        int    `db:"activity"`
}
type Request struct {
	RecipientUsername string `json:"recipient_username"`
	Amount            string `json:"amount"`
}
type Users struct {
	UsersList map[string]User `json:"users_list"`
}

type Transaction struct {
	Sender         string `db:"sender_user_name"`
	SBalance       string `db:"sender_balance"`
	SResultBalance string `db:"sender_result_balance"`
	Recipient      string `db:"recipient_user_name"`
	RBalance       string `db:"recipient_balance"`
	RResultBalance string `db:"recipient_result_balance"`
	Amount         string `db:"amount"`
}

type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}
type TransactionInfo struct {
	// a UUID of transaction
	ID string
	// in USD, typically a value between "0.01" and "1000" USD.
	Amount string
	// bank name, e.g. "Bank of Scotland"
	BankName string
	// a 2-letter country code of where the bank is located
	BankCountryCode string

	TransactionRef *Transaction
}
type Transactions struct {
	TxList    []TransactionInfo
	Latencies map[string]int
}

func (tx Transactions) Len() int {
	return len(tx.TxList)
}
func (tx Transactions) Less(i, j int) bool {
	firstAmount, err := decimal.NewFromString(tx.TxList[i].Amount)
	if err != nil {
		return false
	}
	secondAmount, err := decimal.NewFromString(tx.TxList[j].Amount)
	if err != nil {
		return false
	}
	return firstAmount.Div(decimal.NewFromInt(int64(tx.Latencies[tx.TxList[i].BankCountryCode]))).LessThan(secondAmount.Div(decimal.NewFromInt(int64(tx.Latencies[tx.TxList[j].BankCountryCode]))))
}

type FraudDetectionResult struct {
	TransactionID string
	IsFraudulent  bool
}

func (tx Transactions) Swap(i, j int) {
	tx.TxList[i], tx.TxList[j] = tx.TxList[j], tx.TxList[i]
}

type FraudDetectionResults []FraudDetectionResult

func Sum(array []TransactionInfo) string {
	result := decimal.NewFromInt(0)
	for _, v := range array {
		amountVal, _ := decimal.NewFromString(v.Amount)
		result = result.Add(amountVal)
	}
	return result.String()
}

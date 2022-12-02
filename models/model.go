package models

type User struct {
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
	Balance  string `json:"balance" db:"balance"`
	Activity int    `db:"activity"`
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

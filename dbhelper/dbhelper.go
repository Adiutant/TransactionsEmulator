package dbhelper

import (
	"TestBank/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS Users (
    user_name text NOT NULL PRIMARY KEY,
    password text,
    balance text,
    activity int,
    bankCountryCode text,
    bankName text
);
CREATE TABLE IF NOT EXISTS Transactions(
    id serial NOT NULL PRIMARY KEY ,
    sender_user_name text,
    sender_balance text,
    sender_result_balance text,
    recipient_user_name text,
    recipient_balance text,
    recipient_result_balance text,
    amount text
                                       
);`

type DBHelper struct {
	dbConnection *sqlx.DB
}

func NewDbHelper() (DBHelper, error) {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=ts password=pass dbname=bank sslmode=disable")
	if err != nil {
		return DBHelper{}, err
	}
	db.MustExec(schema)
	helper := DBHelper{dbConnection: db}
	return helper, nil
}
func (dbHelper *DBHelper) GetUser(userName string) (models.User, error, bool) {
	user := models.User{}
	err := dbHelper.dbConnection.Get(&user, "SELECT * FROM Users WHERE user_name=$1 LIMIT 1", userName)
	if err != nil {
		return models.User{}, err, false
	}
	if user.UserName == "" {
		return user, nil, false
	}
	return user, nil, true
}
func (dbHelper *DBHelper) ExecuteTransaction(transaction models.Transaction) (bool, error) {
	tx, err := dbHelper.dbConnection.Begin()
	if err != nil {
		return false, err
	}
	result, err := tx.Exec("UPDATE Users SET balance=$1 WHERE user_name=$2", transaction.SResultBalance, transaction.Sender)
	if err != nil {
		tx.Rollback()
		return false, err

	}
	affected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if affected == 0 {
		tx.Rollback()
		return false, nil
	}
	result, err = tx.Exec("UPDATE Users SET balance=$1 WHERE user_name=$2", transaction.RResultBalance, transaction.Recipient)
	if err != nil {
		tx.Rollback()
		return false, err

	}
	affected, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if affected == 0 {
		tx.Rollback()
		return false, nil
	}

	result, err = dbHelper.dbConnection.NamedExec("INSERT INTO Transactions(sender_user_name,"+
		" sender_balance, sender_result_balance, recipient_user_name,"+
		" recipient_balance, recipient_result_balance,"+
		" amount) VALUES (:sender_user_name, :sender_balance, :sender_result_balance, :recipient_user_name, :recipient_balance, :recipient_result_balance, :amount)", &transaction)
	if err != nil {
		return false, err

	}
	affected, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if affected == 0 {
		tx.Rollback()
		return false, nil
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

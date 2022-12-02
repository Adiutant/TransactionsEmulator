package engine

import (
	"TestBank/dbhelper"
	"TestBank/models"
)

type BankEngine struct {
	Users    models.Users
	LastUser models.User
	DbHelper *dbhelper.DBHelper
}

func NewEngine() (*BankEngine, error) {
	dbHelper, err := dbhelper.NewDbHelper()
	if err != nil {
		return nil, err
	}
	return &BankEngine{DbHelper: &dbHelper, Users: models.Users{UsersList: map[string]models.User{}}}, nil
}

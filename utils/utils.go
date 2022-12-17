package utils

import (
	"TestBank/models"
	"encoding/json"
	"io/ioutil"
	"sort"
	"time"
)

func Prioritize(tx []models.TransactionInfo, totalTime time.Duration) ([]models.TransactionInfo, error) {
	jsonBytes, err := ioutil.ReadFile("data/api_latencies.json")
	if err != nil {
		return nil, err
	}
	var lat map[string]int
	json.Unmarshal(jsonBytes, &lat)
	transactions := models.Transactions{
		TxList:    tx,
		Latencies: lat,
	}
	sort.Sort(transactions)
	currentTime := 0
	resultSlice := make([]models.TransactionInfo, 0)
	for i := len(transactions.TxList) - 1; i >= 0; i-- {
		if int64(currentTime+transactions.Latencies[transactions.TxList[i].BankCountryCode]) > totalTime.Milliseconds() {
			break
		}
		currentTime += transactions.Latencies[transactions.TxList[i].BankCountryCode]
		resultSlice = append(resultSlice, transactions.TxList[i])
	}
	return resultSlice, nil

}

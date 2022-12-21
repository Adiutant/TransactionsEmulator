package http_server

import (
	"TestBank/engine"
	"TestBank/models"
	"TestBank/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type BankServer struct {
	HttpEngine   *gin.Engine
	Engine       *engine.BankEngine
	logger       *logrus.Logger
	queue        []models.TransactionInfo
	processTimer *time.Ticker
	mutex        sync.Mutex
}

func (s *BankServer) auth(c *gin.Context) {
	// A *model.User will eventually be added to context in middleware
	user, pass, ok := c.Request.BasicAuth()
	if !ok {
		c.Abort()
	}
	getUser, err, ok := s.Engine.DbHelper.GetUser(user)
	if err != nil {
		s.logger.Error(err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if /*hash()*/ pass != getUser.Password {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	s.Engine.Users.UsersList[user] = getUser
	c.Next()

}
func NewBankServer() (*BankServer, error) {
	bankEngine, err := engine.NewEngine()
	if err != nil {
		return nil, err
	}
	server := BankServer{
		HttpEngine:   gin.Default(),
		Engine:       bankEngine,
		logger:       logrus.New(),
		queue:        make([]models.TransactionInfo, 0),
		processTimer: time.NewTicker(time.Second),
	}
	return &server, nil
}

func remove(slice []models.TransactionInfo, s int) []models.TransactionInfo {
	return append(slice[:s], slice[s+1:]...)
}

func (s *BankServer) Run() error {
	err := s.HttpEngine.Run(":8080")
	s.processTimer.Reset(time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (s *BankServer) processTransactions() {
	for {
		if len(s.queue) == 0 {
			continue
		}
		select {
		case <-s.processTimer.C:
			copyQueue := s.queue[:]
			fmt.Println(s.queue)
			prioritized, err := utils.Prioritize(copyQueue, time.Second)
			fmt.Println(prioritized)
			if err != nil {
				panic(err)
			}
			s.mutex.Lock()
			for i, tx := range prioritized {
				ok, err := s.Engine.DbHelper.ExecuteTransaction(*tx.TransactionRef)
				s.queue = remove(s.queue, i)
				if err != nil || !ok {
					s.logger.Error(err)
					break
				}
			}
			s.mutex.Unlock()
			s.queue = make([]models.TransactionInfo, 0)
		}
	}

}

func (s *BankServer) SetRoutes() {
	go s.processTransactions()
	s.HttpEngine.Use(s.auth).POST("/transact", func(context *gin.Context) {
		var req models.Request
		var transaction models.Transaction
		user, _, _ := context.Request.BasicAuth()
		defer func() { delete(s.Engine.Users.UsersList, user) }()
		userInfo := s.Engine.Users.UsersList[user]
		err := context.BindJSON(&req)
		if err != nil {
			s.logger.Error(err)
			context.AbortWithError(http.StatusBadRequest, err)

			return
		}
		recipientInfo, err, ok := s.Engine.DbHelper.GetUser(req.RecipientUsername)
		if err != nil {
			s.logger.Error(err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !ok {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		transaction.Sender = userInfo.UserName
		transaction.SBalance = userInfo.Balance
		transaction.Recipient = recipientInfo.UserName
		transaction.RBalance = recipientInfo.Balance

		transaction.Amount = req.Amount
		sBalanceVal, err := decimal.NewFromString(transaction.SBalance)
		if err != nil {
			s.logger.Error(err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return

		}
		amountVal, err := decimal.NewFromString(transaction.Amount)
		if err != nil {
			s.logger.Error(err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return

		}
		rBalanceVal, err := decimal.NewFromString(transaction.RBalance)
		if err != nil {
			s.logger.Error(err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return

		}
		if sBalanceVal.Sub(amountVal).LessThan(decimal.NewFromInt(0)) {
			context.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		transaction.SResultBalance = sBalanceVal.Sub(amountVal).String()
		transaction.RResultBalance = rBalanceVal.Add(amountVal).String()
		hashIn := md5.Sum([]byte(transaction.Recipient + transaction.Sender + transaction.Amount))
		s.queue = append(s.queue, models.TransactionInfo{
			ID:              hex.EncodeToString(hashIn[:]),
			Amount:          transaction.Amount,
			BankName:        s.Engine.Users.UsersList[user].BankName,
			BankCountryCode: s.Engine.Users.UsersList[user].BankCountryCode,
			TransactionRef:  &transaction,
		})

	})

}

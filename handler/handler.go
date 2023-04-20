package handler

import (
	"net/http"

	"github.com/GarnBarn/common-go/httpserver"
	"github.com/GarnBarn/gb-account-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return AccountHandler{
		accountService: accountService,
	}
}

func (a *AccountHandler) GetAccount(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		uid = c.Param(httpserver.UserUidKey)
	}

	account, err := a.accountService.GetUserByUid(uid)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (a *AccountHandler) CreateAccount(c *gin.Context) {
	uid := c.Param(httpserver.UserUidKey)
	account, err := a.accountService.CreateUser(uid)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

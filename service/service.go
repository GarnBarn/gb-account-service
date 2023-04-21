package service

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/GarnBarn/gb-account-service/model"
	"github.com/GarnBarn/gb-account-service/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountService interface {
	GetUserByUid(uid string) (account model.AccountPublic, err error)
}

type accountService struct {
	accountRepository repository.AccountRepository
	app               *firebase.App
}

func NewAccountService(app *firebase.App, accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		app:               app,
		accountRepository: accountRepository,
	}
}

func (a *accountService) GetUserByUid(uid string) (account model.AccountPublic, err error) {
	ctx := context.Background()
	authClient, err := a.app.Auth(ctx)
	if err != nil {
		return account, err
	}

	// Get Account From Database
	accountPrivate, err := a.accountRepository.GetAccountByUid(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return a.createAccount(uid, authClient, ctx)
		} else {
			logrus.Error("Can't get account from database: ", err)
			return account, err
		}
	}

	// Fill the Account Information by using data from Firebase
	user, err := authClient.GetUser(ctx, uid)
	if err != nil {
		return account, err
	}

	return model.ToAccountPublic(accountPrivate, user.DisplayName, user.PhotoURL), nil
}

func (a *accountService) createAccount(uid string, auth *auth.Client, ctx context.Context) (account model.AccountPublic, err error) {
	user, err := auth.GetUser(ctx, uid)
	if err != nil {
		logrus.Error("Create account error, firebase user data not found: ", err)
		return account, err
	}
	accountPrivate, err := a.accountRepository.CreateAccountByUid(uid)
	if err != nil {
		logrus.Error("Failed to create account, uid: ", uid, "  error: ", err)
		return account, err
	}
	logrus.Infof("Create account successed")
	return model.ToAccountPublic(accountPrivate, user.DisplayName, user.PhotoURL), nil
}

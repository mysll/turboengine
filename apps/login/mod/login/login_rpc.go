package login

import (
	"turboengine/core/plugin/storage"
	"turboengine/gameplay/dao"
)

type LoginServer struct {
	storage *storage.Storage
}

func (l *LoginServer) Login(user string, pass string) (bool, error) {
	var account dao.Account
	if err := l.storage.FindBy(&account, "account=? and password=?", user, pass); err != nil {
		return false, err
	}
	if account.Account == user {
		return true, nil
	}
	return false, nil
}

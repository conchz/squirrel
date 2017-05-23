package models_test

import (
	"github.com/lavenderx/squirrel/app"
	"github.com/lavenderx/squirrel/app/crypto"
	"github.com/lavenderx/squirrel/app/models"
	"testing"
	"time"
)

var mySQLTemplate *app.MySQLTemplate

func init() {
	app.ConnectToMySQL(app.LoadConfig())

	mySQLTemplate = app.GetMySQLTemplate()
	if err := mySQLTemplate.XormEngine().Sync2(new(models.User)); err != nil {
		panic(err)
	}
}

func TestUser_Insert(t *testing.T) {
	user := &models.User{
		Username:    "test",
		Password:    crypto.EncryptPassword([]byte("testSecret")),
		Cellphone:   "156××××××××",
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	affected, err := mySQLTemplate.Insert(user)
	if err != nil {
		panic(err)
	}
	t.Logf("Affect number: %d, User: %+v\n", affected, user)
}

func TestUser_FindById(t *testing.T) {
	user := mySQLTemplate.GetById(1, new(models.User))
	t.Logf("User: %+v\n", user)
}
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

func TestUser_Delete(t *testing.T) {
	user := &models.User{
		Username: "test",
	}

	_, err := mySQLTemplate.DeleteByNonEmptyFields(user)
	if err != nil {
		panic(err)
	}
}

func TestUser_Insert(t *testing.T) {
	user := &models.User{
		Username:    "test",
		Password:    crypto.EncryptPassword([]byte("passwd")),
		Secret:      crypto.RandStringBytesMaskImpr(1 << 5),
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

func TestUser_Find(t *testing.T) {
	user := mySQLTemplate.GetByNonEmptyFields(&models.User{
		Username: "test",
	})
	t.Logf("User: %+v\n", user)
}

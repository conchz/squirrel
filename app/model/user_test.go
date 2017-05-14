package model

import (
	"github.com/lavenderx/squirrel/app"
	"github.com/lavenderx/squirrel/app/crypto"
	"testing"
	"time"
)

func init() {
	app.ConnectToMySQL(app.LoadConfig())

	engine := app.GetXormEngine()
	if err := engine.Sync2(new(User)); err != nil {
		panic(err)
	}
}

func TestUser_Insert(t *testing.T) {
	user := &User{
		Username:    "Baymax",
		Password:    crypto.EncryptPassword([]byte("test1234")),
		Cellphone:   "19012345678",
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	affected, err := app.Insert(user)
	if err != nil {
		panic(err)
	}
	t.Logf("Affect number: %d, User: %+v\n", affected, user)
}

func TestUser_FindById(t *testing.T) {
	user := app.FindById(1, new(User))
	t.Logf("User: %+v\n", user)
}

package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/lavenderx/squirrel/app"
	"testing"
	"time"
)

func init() {
	config := app.LoadConfig()
	if err := app.ConnectToMySQL(config); err != nil {
		panic(err)
	}
}

func TestUserInsert(t *testing.T) {
	user := User{
		Username:    "Baymax",
		Password:    "test1234",
		Cellphone:   "19012345678",
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	engine := app.GetXormEngine()
	if err := engine.Sync2(user); err != nil {
		panic(err)
	}

	defer func(engine *xorm.Engine) {
		fmt.Println("Closing MySQL client...")
		err := engine.Close()
		if err != nil {
			fmt.Printf("%v\n", "Closing MySQL client failed!")
			panic(err)
		}
	}(engine)

	affectedNumber, err := user.Insert(user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v records are affected\n", affectedNumber)
}

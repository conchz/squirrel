package model

import (
	"github.com/lavenderx/squirrel/app"
	"time"
)

type User struct {
	Id          int64
	Username    string `xorm:"varchar(20) not null unique 'user_name'"`
	Password    string `xorm:"VARCHAR(20) not null"`
	CreatedTime time.Time `xorm:"not null 'created_time'"`
	UpdatedTime time.Time `xorm:"not null 'updated_time'"`
}

func (u *User) TableName() string {
	return "tbl_user"
}

func (u *User) Insert(user User) (int64, error) {
	engine := app.GetXormEngine()

	err := engine.Sync2(new(User))
	if err != nil {
		panic(err)
	}

	affected, err := engine.Insert(&user)
	if err != nil {
		return affected, err
	}

	return affected, nil
}

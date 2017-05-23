package models

import (
	"time"
)

type User struct {
	Id          int64
	Username    string    `xorm:"varchar(20) not null unique"`
	Password    string    `xorm:"varchar(20) not null"`
	Cellphone   string    `xorm:"varchar(20) not null unique"`
	Email       string    `xorm:"varchar(20)"`
	CreatedTime time.Time `xorm:"not null 'created_time'"`
	UpdatedTime time.Time `xorm:"not null 'updated_time'"`
}

func (u *User) TableName() string {
	return "tbl_user"
}

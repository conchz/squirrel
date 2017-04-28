package app

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
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

var engine *xorm.Engine

func ConnectToMySQL() *xorm.Engine {
	if engine != nil {
		return engine
	}

	mysqlConfig := GetConfig().MySQLConf

	var engineErr error
	engine, engineErr = xorm.NewEngine(mysqlConfig.Driver, mysqlConfig.DataSource)
	if engineErr != nil {
		panic(engineErr)
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	engine.TZLocation = location

	return engine
}

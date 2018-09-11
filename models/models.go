package models

import (
	"github.com/go-xorm/xorm"
	"github.com/saisai/gindemo/api/msg"
)

type User struct {
	Id         string `json:"id" xorm:"varchar(24) pk"`
	Nickname   string `json:"nickname" xorm:"varchar(100)"`
	Avatar     string `json:"avatar" xorm:"varchar(100)"`
	Sex        int    `json:"sex" xorm:"int"`
	CreateTime string `json:"createtime" xorm:"DateTime created"`
}

type UserAuths struct {
	Id              int    `json:"id" xorm:"int pk"`
	UserId          string `json:"user_id" xorm:"varchar(100)"`
	IdentifyType    string `json:"identify_type" xorm:"varchar(50)"`
	Identifier      string `json:"identifier" xorm:"varchar(50)"`
	Credential      string `json:"credential" xorm:"varchar(100)"`
	Latestlogintime string `json:"latestlogintime" xorm:"DateTime"`
	State           int    `json:"state" xorm:"int"`
	Registertime    string `json:"registertime" xorm:"DateTime"`
}

var (
	DBEngine *xorm.Engine
)

var (
	coreTables []interface{} = []interface{}{
		new(msg.User), new(UserAuths),
	}
)

func DB() *xorm.Engine {
	return DBEngine
}

func SyncTables() error {
	return DB().Sync2(coreTables...)
}

func InitDB(e *xorm.Engine) {
	DBEngine = e
}

package mysql_test

import (
	"github.com/bang-go/crab/core/db/mysql"
	"gorm.io/gorm/schema"
	"log"
	"testing"
)

func TestConn(t *testing.T) {
	opt := mysql.Config{
		Dsn: mysql.DsnConfig{User: "test", Passwd: "test", Net: "tcp", Addr: "local:3306", DBName: "test", AllowNativePasswords: true},
		Orm: mysql.GormConfig{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	}
	Client, err := mysql.New(&opt)
	if err != nil {
		log.Println(err)
	}
	type YwUser struct {
		Uid      int64 `gorm:"primary_key" json:"uid"`
		Username string
	}
	user := &YwUser{}
	Client.Find(user)
	log.Println(user)
}

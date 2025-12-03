package mysqlx_test

import (
	"log"
	"testing"

	"github.com/bang-go/crab/core/db/mysqlx"
	"gorm.io/gorm/schema"
)

func TestConn(t *testing.T) {
	opt := mysqlx.Config{
		Dsn: mysqlx.DsnConfig{User: "test", Passwd: "test", Net: "tcp", Addr: "local:3306", DBName: "test", AllowNativePasswords: true},
		Orm: mysqlx.GormConfig{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	}
	Client, err := mysqlx.New(&opt)
	if err != nil {
		log.Println(err)
	}
	type YwUser struct {
		Uid      int64 `gorm:"primary_key" json:"uid"`
		Username string
	}
	user := &YwUser{}
	Client.GetDB().Find(user)
	log.Println(user)
}

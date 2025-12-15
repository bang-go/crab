package mysqlx_test

import (
	"log"
	"testing"

	"github.com/bang-go/crab/core/db/mysqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestConn(t *testing.T) {
	config := mysqlx.ClientConfig{
		DSN: &mysqlx.DSNConfig{
			User:                 "test",
			Passwd:               "test",
			Net:                  "tcp",
			Addr:                 "localhost:3306",
			DBName:               "test",
			AllowNativePasswords: true,
		},
		Gorm: &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	}

	client, err := mysqlx.New(&config)
	if err != nil {
		log.Println("MySQL connection failed:", err)
		return
	}
	defer client.Close()

	type YwUser struct {
		Uid      int64 `gorm:"primary_key" json:"uid"`
		Username string
	}

	user := &YwUser{}
	client.GetDB().Find(user)
	log.Println(user)
}

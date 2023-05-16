package mysql

import (
	"context"
	"database/sql"
	driver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)
import "gorm.io/driver/mysql"

type DsnConfig driver.Config //DSN config
type SqlConfig mysql.Config
type GormConfig gorm.Config
type Client = gorm.DB
type PoolConfig struct {
	MaxIdleConns    int           //最大空闲连接数,默认2
	MaxOpenConns    int           //最大Open链接数,默认0 不限制
	ConnMaxLifetime time.Duration //链接最大生命周期
	ConnMaxIdleTime time.Duration //链接最大空闲周期
}
type Config struct {
	Ctx  context.Context
	Dsn  DsnConfig
	Sql  SqlConfig
	Orm  GormConfig
	Pool *PoolConfig
}

func New(conf *Config) (db *Client, err error) {
	conf.setDsn()
	ormConfig := gorm.Config(conf.Orm)
	// DSN [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	db, err = gorm.Open(mysql.New(mysql.Config(conf.Sql)), &ormConfig)
	if err != nil {
		return
	}
	sqlDb, err := db.DB()
	if err != nil {
		return
	}
	conf.setPool(sqlDb) //设置连接池
	return

}

// 解析DSN CONFIG到字符串
func (c *Config) setDsn() {
	cf := driver.Config(c.Dsn)
	c.Sql.DSN = cf.FormatDSN()
}

// SetPool 设置数据池
func (c *Config) setPool(sqlDB *sql.DB) {
	if c.Pool != nil {
		sqlDB.SetMaxIdleConns(c.Pool.MaxIdleConns)
		sqlDB.SetMaxOpenConns(c.Pool.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(c.Pool.ConnMaxLifetime)
		sqlDB.SetConnMaxIdleTime(c.Pool.ConnMaxIdleTime)
	}
}

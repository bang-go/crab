package mysqlx

import (
	"context"
	"database/sql"
	driver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)
import "gorm.io/driver/mysql"

type DsnConfig = driver.Config //DSN config
type SqlConfig = mysql.Config
type GormConfig = gorm.Config
type PoolConfig struct {
	MaxIdleConns    int           //最大空闲连接数,默认2
	MaxOpenConns    int           //最大Open链接数,默认0 不限制
	ConnMaxLifetime time.Duration //链接最大生命周期
	ConnMaxIdleTime time.Duration //链接最大空闲周期
}
type Client struct {
	db    *gorm.DB
	sqlDB *sql.DB
}
type Config struct {
	Ctx  context.Context
	Dsn  DsnConfig
	Sql  SqlConfig
	Orm  GormConfig
	Pool *PoolConfig
}

func New(conf *Config) (client *Client, err error) {
	client = &Client{}
	conf.setDsn()
	// DSN [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	client.db, err = gorm.Open(mysql.New(conf.Sql), &conf.Orm)
	if err != nil {
		return
	}
	client.sqlDB, err = client.db.DB()
	if err != nil {
		return
	}
	conf.setPool(client.sqlDB) //设置连接池
	return

}

// 解析DSN CONFIG到字符串
func (c *Config) setDsn() {
	c.Sql.DSN = c.Dsn.FormatDSN()
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

func (s *Client) GetDB() *gorm.DB {
	return s.db
}
func (s *Client) GetSqlDB() *sql.DB {
	return s.sqlDB
}
func (s *Client) Close() error {
	return s.sqlDB.Close()
}

type Transaction struct {
	*Client
	txDB *gorm.DB
}

func (s *Client) NewTransaction() *Transaction {
	return &Transaction{Client: s}
}

func (s *Transaction) GetTxDB() *gorm.DB {
	return s.txDB
}

func (s *Transaction) Begin() {
	s.txDB = s.Client.db.Begin()
}
func (s *Transaction) Rollback() {
	s.txDB.Rollback()
}
func (s *Transaction) Commit() {
	s.txDB.Commit()
}

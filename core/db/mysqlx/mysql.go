package mysqlx

import (
	"context"
	"database/sql"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DsnConfig  = driver.Config
	SqlConfig  = mysql.Config
	GormConfig = gorm.Config
)

type PoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
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

func New(conf *Config) (*Client, error) {
	conf.setDsn()
	db, err := gorm.Open(mysql.New(conf.Sql), &conf.Orm)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	conf.setPool(sqlDB)
	return &Client{db: db, sqlDB: sqlDB}, nil
}

func (c *Config) setDsn() {
	c.Sql.DSN = c.Dsn.FormatDSN()
}

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

func (s *Transaction) Begin() error {
	if s.Client == nil || s.Client.db == nil {
		return gorm.ErrInvalidTransaction
	}
	s.txDB = s.Client.db.Begin()
	return s.txDB.Error
}
func (s *Transaction) Rollback() error {
	if s.txDB == nil {
		return gorm.ErrInvalidTransaction
	}
	return s.txDB.Rollback().Error
}
func (s *Transaction) Commit() error {
	if s.txDB == nil {
		return gorm.ErrInvalidTransaction
	}
	return s.txDB.Commit().Error
}

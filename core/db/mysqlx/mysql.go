package mysqlx

import (
	"context"
	"database/sql"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 直接暴露第三方库的类型，方便业务方使用
type (
	DSNConfig  = driver.Config // MySQL DSN 配置
	Config     = mysql.Config  // GORM MySQL 配置
	GormConfig = gorm.Config   // GORM 配置
)

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxIdleConns    int           // 最大空闲连接数
	MaxOpenConns    int           // 最大打开连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
}

// DefaultPoolConfig 默认连接池配置
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MaxIdleConns:    10,                 // 最大空闲连接数
		MaxOpenConns:    100,                // 最大打开连接数
		ConnMaxLifetime: 3600 * time.Second, // 连接最大生命周期 1小时
		ConnMaxIdleTime: 600 * time.Second,  // 连接最大空闲时间 10分钟
	}
}

// Client MySQL 客户端封装
type Client struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

// ClientConfig 客户端配置（组合所有配置）
type ClientConfig struct {
	// DSN 配置（可以直接使用 driver.Config 构建）
	DSN *DSNConfig

	// 或者直接提供 DSN 字符串
	DSNString string

	// MySQL 配置（可选）
	MySQL *Config

	// GORM 配置（可选）
	Gorm *GormConfig

	// 连接池配置（可选）
	Pool *PoolConfig
}

// New 创建 MySQL 客户端并验证连接
func New(conf *ClientConfig) (*Client, error) {
	// 构建 DSN
	var dsnString string
	if conf.DSNString != "" {
		dsnString = conf.DSNString
	} else if conf.DSN != nil {
		dsnString = conf.DSN.FormatDSN()
	} else {
		return nil, gorm.ErrInvalidDB
	}

	// 构建 MySQL 配置
	mysqlConfig := mysql.Config{DSN: dsnString}
	if conf.MySQL != nil {
		mysqlConfig = *conf.MySQL
		mysqlConfig.DSN = dsnString
	}

	// 构建 GORM 配置
	gormConfig := &gorm.Config{}
	if conf.Gorm != nil {
		gormConfig = conf.Gorm
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
	if err != nil {
		return nil, err
	}

	// 获取 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池
	setPool(sqlDB, conf.Pool)

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return &Client{db: db, sqlDB: sqlDB}, nil
}

// setPool 设置连接池配置
func setPool(sqlDB *sql.DB, pool *PoolConfig) {
	// 如果未提供 Pool 配置，使用默认值
	if pool == nil {
		pool = DefaultPoolConfig()
	}
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(pool.ConnMaxIdleTime)
}

// GetDB 获取 GORM DB 实例
func (c *Client) GetDB() *gorm.DB {
	return c.db
}

// GetSqlDB 获取原生 sql.DB 实例
func (c *Client) GetSqlDB() *sql.DB {
	return c.sqlDB
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	return c.sqlDB.Close()
}

// Transaction 事务封装
type Transaction struct {
	*Client
	txDB *gorm.DB
}

// NewTransaction 创建事务实例（不推荐，使用 WithTransaction 更安全）
func (c *Client) NewTransaction() *Transaction {
	return &Transaction{Client: c}
}

// GetTxDB 获取事务的 DB 实例
func (t *Transaction) GetTxDB() *gorm.DB {
	return t.txDB
}

// Begin 开始事务
func (t *Transaction) Begin() error {
	if t.Client == nil || t.Client.db == nil {
		return gorm.ErrInvalidTransaction
	}
	t.txDB = t.Client.db.Begin()
	return t.txDB.Error
}

// Rollback 回滚事务
func (t *Transaction) Rollback() error {
	if t.txDB == nil {
		return gorm.ErrInvalidTransaction
	}
	return t.txDB.Rollback().Error
}

// Commit 提交事务
func (t *Transaction) Commit() error {
	if t.txDB == nil {
		return gorm.ErrInvalidTransaction
	}
	return t.txDB.Commit().Error
}

// WithTransaction 安全地执行事务（自动 Commit/Rollback）
// 例：
//
//	err := client.WithTransaction(func(tx *gorm.DB) error {
//	    return tx.Create(&user).Error
//	})
func (c *Client) WithTransaction(fn func(tx *gorm.DB) error) error {
	return c.db.Transaction(fn)
}

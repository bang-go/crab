package postgresx

import (
	"context"
	"database/sql"
	"time"

	"github.com/bang-go/opt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config 为 GORM PostgreSQL 配置
// GormConfig 为 GORM 通用配置
type (
	Config     = postgres.Config
	GormConfig = gorm.Config
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
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600 * time.Second,
		ConnMaxIdleTime: 600 * time.Second,
	}
}

// Options PostgreSQL 客户端构建参数（推荐通过 Option 方式设置）
type Options struct {
	DSN    string      // 连接串，如："host=localhost user=xxx password=xxx dbname=xxx port=5432 sslmode=disable"
	Config *Config     // 可选：底层 postgres.Config，若设置则优先使用
	Gorm   *GormConfig // 可选：gorm.Config
	Pool   *PoolConfig // 可选：连接池配置
}

// WithDSN 设置 DSN
func WithDSN(dsn string) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.DSN = dsn
	})
}

// WithConfig 设置底层 postgres.Config
func WithConfig(cfg Config) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		// 拷贝一份，避免外部修改
		c := cfg
		o.Config = &c
	})
}

// WithGormConfig 设置 gorm.Config
func WithGormConfig(cfg *GormConfig) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.Gorm = cfg
	})
}

// WithPool 设置连接池配置
func WithPool(pool *PoolConfig) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.Pool = pool
	})
}

// Client PostgreSQL 客户端封装
type Client struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

// Open 创建 PostgreSQL 客户端并验证连接（推荐入口）
func Open(opts ...opt.Option[Options]) (*Client, error) {
	// 组装选项
	o := &Options{}
	opt.Each(o, opts...)

	// 构建 DSN
	dsn := o.DSN
	if dsn == "" && o.Config != nil {
		dsn = o.Config.DSN
	}
	if dsn == "" {
		return nil, gorm.ErrInvalidDB
	}

	// 构建 postgres.Config
	var pgCfg Config
	if o.Config != nil {
		pgCfg = *o.Config
		pgCfg.DSN = dsn // 以 Option 传入的 DSN 为准
	} else {
		pgCfg = postgres.Config{DSN: dsn}
	}

	// 构建 gorm.Config
	gormCfg := o.Gorm
	if gormCfg == nil {
		gormCfg = &gorm.Config{}
	}

	// 打开数据库
	db, err := gorm.Open(postgres.New(pgCfg), gormCfg)
	if err != nil {
		return nil, err
	}

	// 获取 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池
	setPool(sqlDB, o.Pool)

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return &Client{db: db, sqlDB: sqlDB}, nil
}

// MustOpen 与 Open 类似，但在出错时 panic
func MustOpen(opts ...opt.Option[Options]) *Client {
	c, err := Open(opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// setPool 设置连接池配置
func setPool(sqlDB *sql.DB, pool *PoolConfig) {
	if pool == nil {
		pool = DefaultPoolConfig()
	}
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(pool.ConnMaxIdleTime)
}

// DB 返回 *gorm.DB，便于直接使用 GORM 能力
func (c *Client) DB() *gorm.DB {
	return c.db
}

// SQLDB 返回底层 *sql.DB
func (c *Client) SQLDB() *sql.DB {
	return c.sqlDB
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	if c == nil || c.sqlDB == nil {
		return nil
	}
	return c.sqlDB.Close()
}

// WithTransaction 使用 GORM 内置事务能力执行函数
func (c *Client) WithTransaction(fn func(tx *gorm.DB) error) error {
	if c == nil || c.db == nil {
		return gorm.ErrInvalidDB
	}
	return c.db.Transaction(fn)
}

// Transaction 事务对象（手动管理模式）
type Transaction struct {
	*Client
	txDB *gorm.DB
}

// NewTransaction 创建事务实例（不推荐直接使用，手动管理模式）
func (c *Client) NewTransaction() *Transaction {
	if c == nil {
		return nil
	}
	return &Transaction{Client: c}
}

// DB 获取事务的 *gorm.DB 实例
func (t *Transaction) DB() *gorm.DB {
	if t == nil {
		return nil
	}
	return t.txDB
}

// Begin 开始事务
func (t *Transaction) Begin() error {
	if t == nil || t.Client == nil || t.Client.db == nil {
		return gorm.ErrInvalidTransaction
	}
	t.txDB = t.Client.db.Begin()
	return t.txDB.Error
}

// Rollback 回滚事务
func (t *Transaction) Rollback() error {
	if t == nil || t.txDB == nil {
		return gorm.ErrInvalidTransaction
	}
	return t.txDB.Rollback().Error
}

// Commit 提交事务
func (t *Transaction) Commit() error {
	if t == nil || t.txDB == nil {
		return gorm.ErrInvalidTransaction
	}
	return t.txDB.Commit().Error
}

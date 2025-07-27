package mysql

import (
	"context"
	"github.com/yanking/micro-zero/pkg/contract"
	"github.com/yanking/micro-zero/pkg/log"
	"github.com/yanking/micro-zero/pkg/options"
	"gorm.io/gorm"
	"time"
)

const componentName = "MySQL"

var _ contract.Component = (*Client)(nil)

type Client struct {
	opts *options.MySQLOptions
	db   *gorm.DB
}

func New(opts *options.MySQLOptions) (contract.Component, error) {
	log.Infof("component %s: client initializing with DSN: %s", componentName, opts.DSN())
	client, err := opts.NewDB()
	if err != nil {
		return nil, err
	}
	return &Client{
		opts: opts,
		db:   client,
	}, nil
}

func (c Client) Start(ctx context.Context) error {
	// 启动一个goroutine来处理断线重连
	go func() {
		// 定时检查数据库连接状态
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Infof("component %s: MySQL component context done, stopping connection checker", componentName)
				return
			case <-ticker.C:
				// 检查数据库连接
				if err := c.ping(); err != nil {
					log.Errorf("component %s: MySQL connection lost: %v, attempting to reconnect...", componentName, err)
					if er := c.reconnect(); er != nil {
						log.Errorf("component %s: failed to reconnect to MySQL: %v", componentName, er)
					} else {
						log.Infof("component %s: successfully reconnected to MySQL", componentName)
					}
				}
			}
		}
	}()
	return nil
}

func (c Client) Stop(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c Client) Name() string {
	return componentName
}

// ping 检查数据库连接是否正常
func (c Client) ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// reconnect 重新连接数据库
func (c Client) reconnect() error {
	client, err := c.opts.NewDB()
	if err != nil {
		return err
	}
	c.db = client
	return nil
}

// GetDB 获取数据库连接实例
func (c Client) GetDB() *gorm.DB {
	return c.db
}

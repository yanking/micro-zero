package app

import (
	"github.com/yanking/micro-zero/pkg/components/httpserver"
	"github.com/yanking/micro-zero/pkg/components/mysql"
	"github.com/yanking/micro-zero/pkg/components/redis"
	"github.com/yanking/micro-zero/pkg/config"
	"github.com/yanking/micro-zero/pkg/container"
	"github.com/yanking/micro-zero/pkg/log"
)

// defaultComponentRunner 是ComponentRunner接口的默认实现
type defaultComponentRunner struct {
	cfg *config.Config
}

// NewDefaultComponentRunner 创建一个新的默认组件运行器
func NewDefaultComponentRunner(cfg *config.Config) ComponentRunner {
	return &defaultComponentRunner{
		cfg: cfg,
	}
}

// RunWithComponents 使用组件容器运行应用
func (d *defaultComponentRunner) RunWithComponents(c *container.Container) error {
	log.Infof("Registering default components...")

	// 注册 MySQL 组件
	if mysqlComponent, err := mysql.New(d.cfg.MySQLOptions); err == nil {
		c.Register(mysqlComponent)
		log.Infof("MySQL component registered")
	} else {
		log.Warnf("Failed to create MySQL component: %v", err)
	}

	// 注册 HTTP Server 组件
	if httpComponent, err := httpserver.New(d.cfg.HTTPOptions); err == nil {
		c.Register(httpComponent)
		log.Infof("HTTP Server component registered")
	} else {
		log.Warnf("Failed to create HTTP Server component: %v", err)
	}

	// 注册 Redis 组件
	if redisComponent, err := redis.New(d.cfg.RedisOptions); err == nil {
		c.Register(redisComponent)
		log.Infof("Redis component registered")
	} else {
		log.Warnf("Failed to create Redis component: %v", err)
	}

	return c.Run()
}

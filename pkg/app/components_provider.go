package app

import (
	"github.com/yanking/micro-zero/pkg/components/httpserver"
	"github.com/yanking/micro-zero/pkg/components/mysql"
	"github.com/yanking/micro-zero/pkg/components/redis"
	"github.com/yanking/micro-zero/pkg/config"
	"github.com/yanking/micro-zero/pkg/contract"
)

// defaultComponentsProvider 提供默认组件的实现
type defaultComponentsProvider struct {
	cfg *config.Config
}

// NewDefaultComponentsProvider 创建一个新的默认组件提供者
func NewDefaultComponentsProvider(cfg *config.Config) DefaultComponents {
	return &defaultComponentsProvider{
		cfg: cfg,
	}
}

// GetDefaultComponents 返回默认的组件列表
func (d *defaultComponentsProvider) GetDefaultComponents() []contract.Component {
	var components []contract.Component

	// 添加 MySQL 组件（如果配置有效）
	if mysqlComponent, err := mysql.New(d.cfg.MySQLOptions); err == nil {
		components = append(components, mysqlComponent)
	}

	// 添加 HTTP Server 组件（如果配置有效）
	if httpComponent, err := httpserver.New(d.cfg.HTTPOptions); err == nil {
		components = append(components, httpComponent)
	}

	// 添加 Redis 组件（如果配置有效）
	if redisComponent, err := redis.New(d.cfg.RedisOptions); err == nil {
		components = append(components, redisComponent)
	}

	return components
}

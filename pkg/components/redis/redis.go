package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yanking/micro-zero/pkg/contract"
	"github.com/yanking/micro-zero/pkg/log"
	"github.com/yanking/micro-zero/pkg/options"
)

var _ contract.Component = (*Client)(nil)

// Client 实现了Component接口的Redis组件
type Client struct {
	opts   *options.RedisOptions
	client *redis.Client
}

// New 创建一个新的Redis组件实例
func New(opts *options.RedisOptions) (contract.Component, error) {
	// 创建Redis客户端
	client, err := opts.NewClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		opts:   opts,
		client: client,
	}, nil
}

// Start 启动Redis组件
func (c *Client) Start(ctx context.Context) error {
	log.Infof("component: Redis client starting with addr: %s", c.opts.Addr)

	// 启动一个后台goroutine定期检查连接状态
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Infof("component: Redis client context done")
				return
			case <-ticker.C:
				// 检查连接
				if err := c.client.Ping(ctx).Err(); err != nil {
					log.Errorf("component: Redis connection error: %v", err)
				}
			}
		}
	}()

	return nil
}

// Stop 停止Redis组件
func (c *Client) Stop(ctx context.Context) error {
	log.Infof("component: Stopping Redis client")
	return c.client.Close()
}

// Name 返回组件名称
func (c *Client) Name() string {
	return "redis-client"
}

// GetClient 返回Redis客户端实例
func (c *Client) GetClient() *redis.Client {
	return c.client
}

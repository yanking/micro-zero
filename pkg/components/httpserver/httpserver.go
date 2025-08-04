package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yanking/micro-zero/pkg/contract"
	"github.com/yanking/micro-zero/pkg/log"
	"github.com/yanking/micro-zero/pkg/options"
)

var _ contract.Component = (*Server)(nil)

// Server 实现了 Component 接口的 HTTP 服务组件
type Server struct {
	opts   *options.HTTPOptions
	server *http.Server
}

// New 创建一个新的 HTTP 服务组件实例
func New(opts *options.HTTPOptions) (contract.Component, error) {
	// 创建一个简单的 HTTP 处理器
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from HTTP component server!")
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}

	return &Server{
		opts:   opts,
		server: server,
	}, nil
}

// Start 启动 HTTP 服务组件
func (s *Server) Start(ctx context.Context) error {
	log.Infof("component: HTTP server starting on %s", s.opts.Addr)

	// 在单独的 goroutine 中启动服务器，以避免阻塞
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("component: HTTP server error: %v", err)
		}
	}()

	// 监听上下文取消信号，用于优雅关闭
	go func() {
		<-ctx.Done()
		log.Infof("component: HTTP server shutting down...")
		// 创建一个带超时的上下文用于关闭
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Errorf("component: HTTP server shutdown error: %v", err)
		}
	}()

	return nil
}

// Stop 停止 HTTP 服务组件
func (s *Server) Stop(ctx context.Context) error {
	log.Infof("component: Stopping HTTP server...")
	return s.server.Shutdown(ctx)
}

// Name 返回组件名称
func (s *Server) Name() string {
	return "http-server"
}

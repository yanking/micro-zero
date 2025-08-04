package app

import (
	"context"

	"github.com/yanking/micro-zero/pkg/contract"
)

// DefaultComponents 定义应用可以使用的默认组件接口
// 这些组件可以在App初始化时选择性地添加
type DefaultComponents interface {
	// GetDefaultComponents 返回应用应该注册的默认组件列表
	GetDefaultComponents() []contract.Component
}

// ComponentProvider 定义组件提供者的接口
type ComponentProvider interface {
	// ProvideComponents 根据配置创建并返回组件列表
	ProvideComponents(ctx context.Context) ([]contract.Component, error)
}

package app

import (
	"github.com/yanking/micro-zero/pkg/container"
)

// ComponentRunner 定义了组件运行器接口
// 实现此接口的应用可以在App中自动集成组件管理
type ComponentRunner interface {
	// RunWithComponents 使用组件容器运行应用
	RunWithComponents(container *container.Container) error
}

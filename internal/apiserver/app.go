package apiserver

import (
	"github.com/yanking/micro-zero/pkg/app"
	"github.com/yanking/micro-zero/pkg/config"
)

func NewApp(name string) *app.App {
	cfg := config.New()

	// 创建默认组件运行器
	defaultComponentRunner := app.NewDefaultComponentRunner(cfg)

	appl := app.NewApp(name, "",
		app.WithOptions(cfg),
		app.WithSilence(),
		app.WithRunFunc(run(name, cfg)),
		app.WithComponentRunner(defaultComponentRunner), // 注册组件运行器
	)
	return appl
}

// 空的run函数，实际运行逻辑由ComponentRunner处理
func run(name string, cfg *config.Config) app.RunFunc {
	return func() error {
		// 实际的组件运行由ComponentRunner处理
		return nil
	}
}

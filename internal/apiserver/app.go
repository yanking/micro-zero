package apiserver

import (
	"github.com/yanking/micro-zero/pkg/app"
	"github.com/yanking/micro-zero/pkg/components/mysql"
	"github.com/yanking/micro-zero/pkg/config"
	"github.com/yanking/micro-zero/pkg/container"
)

func NewApp(name string) *app.App {
	cfg := config.New()
	appl := app.NewApp(name, "",
		app.WithOptions(cfg),
		app.WithSilence(),
		app.WithRunFunc(run(name, cfg)),
	)
	return appl
}

func run(name string, cfg *config.Config) app.RunFunc {
	return func() error {
		cr := container.New(name,
			container.WithShutdownOverallTimeout(cfg.ShutdownOverallTimeout),
		)

		mysqlComponent, err := mysql.New(cfg.MySQLOptions)
		if err != nil {
			return err
		}
		cr.Register(mysqlComponent)

		return cr.Run()
	}
}

package app

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"
	"k8s.io/component-base/cli"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"

	"github.com/yanking/micro-zero/pkg/config"
	"github.com/yanking/micro-zero/pkg/container"
	"github.com/yanking/micro-zero/pkg/log"
	genericoptions "github.com/yanking/micro-zero/pkg/options"
	"github.com/yanking/micro-zero/pkg/version"
)

// App is the main structure of a cli application.
// It is recommended that an container be created with the app.NewApp() function.
type App struct {
	name        string
	shortDesc   string
	description string
	run         RunFunc
	cmd         *cobra.Command
	args        cobra.PositionalArgs

	// +optional
	healthCheckFunc HealthCheckFunc

	// +optional
	options any

	// +optional
	silence bool

	// +optional
	noConfig bool

	// watching and re-reading config files
	// +optional
	watch bool

	// +optional
	defaultComponents DefaultComponents

	// +optional
	componentRunner ComponentRunner

	contextExtractors map[string]func(context.Context) string
}

// RunFunc defines the application's startup callback function.
type RunFunc func() error

// HealthCheckFunc defines the health check function for the application.
type HealthCheckFunc func() error

// Option defines optional parameters for initializing the application
// structure.
type Option func(*App)

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file.
func WithOptions(opts any) Option {
	return func(app *App) {
		app.options = opts
	}
}

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.run = run
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

// WithHealthCheckFunc is used to set the health check function for the application.
// The container framework will use the function to start a health check server.
func WithHealthCheckFunc(fn HealthCheckFunc) Option {
	return func(app *App) {
		app.healthCheckFunc = fn
	}
}

// WithDefaultHealthCheckFunc set the default health check function.
func WithDefaultHealthCheckFunc() Option {
	fn := func() HealthCheckFunc {
		return func() error {
			go genericoptions.NewHealthOptions().ServeHealthCheck()

			return nil
		}
	}

	return WithHealthCheckFunc(fn())
}

// WithSilence sets the application to silent mode, in which the program startup
// information, configuration information, and version information are not
// printed in the console.
func WithSilence() Option {
	return func(app *App) {
		app.silence = true
	}
}

// WithNoConfig set the application does not provide config flag.
func WithNoConfig() Option {
	return func(app *App) {
		app.noConfig = true
	}
}

// WithValidArgs set the validation function to valid non-flag arguments.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(app *App) {
		app.args = cobra.NoArgs
	}
}

// WithWatchConfig watching and re-reading config files.
func WithWatchConfig() Option {
	return func(app *App) {
		app.watch = true
	}
}

// WithDefaultComponents sets the application's default components provider.
func WithDefaultComponents(provider DefaultComponents) Option {
	return func(app *App) {
		app.defaultComponents = provider
	}
}

// WithComponentRunner sets the application's component runner.
func WithComponentRunner(runner ComponentRunner) Option {
	return func(app *App) {
		app.componentRunner = runner
	}
}

func WithLoggerContextExtractor(contextExtractors map[string]func(context.Context) string) Option {
	return func(app *App) {
		app.contextExtractors = contextExtractors
	}
}

// NewApp creates a new application instance based on the given application name,
// binary name, and other options.
func NewApp(name string, shortDesc string, opts ...Option) *App {
	app := &App{
		name:      name,
		run:       func() error { return nil },
		shortDesc: shortDesc,
	}

	for _, o := range opts {
		o(app)
	}

	app.buildCommand()

	return app
}

// buildCommand is used to build a cobra command.
func (app *App) buildCommand() {
	cmd := &cobra.Command{
		Use:   formatBaseName(app.name),
		Short: app.shortDesc,
		Long:  app.description,
		RunE:  app.runCommand,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return nil
		},
		Args: app.args,
	}
	// When error printing is enabled for the Cobra command, a flag parse
	// error gets printed first, then optionally the often long usage
	// text. This is very unreadable in a console because the last few
	// lines that will be visible on screen don't include the error.
	//
	// The recommendation from #sig-cli was to print the usage text, then
	// the error. We implement this consistently for all commands here.
	// However, we don't want to print the usage text when command
	// execution fails for other reasons than parsing. We detect this via
	// the FlagParseError callback.
	//
	// Some commands, like kubectl, already deal with this themselves.
	// We don't change the behavior for those.
	if !cmd.SilenceUsage {
		cmd.SilenceUsage = true
		cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
			// Re-enable usage printing.
			c.SilenceUsage = false
			return err
		})
	}
	// In all cases error printing is done below.
	cmd.SilenceErrors = true

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true

	var fs *pflag.FlagSet
	// 方法2：使用type switch
	switch typed := app.options.(type) {
	case NamedFlagSetOptions:
		var fss cliflag.NamedFlagSets
		fs = fss.FlagSet("global")

		if app.options != nil {
			fss = typed.Flags()
			fs = fss.FlagSet("global")
		}

		for _, f := range fss.FlagSets {
			cmd.Flags().AddFlagSet(f)
		}

		cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
		cliflag.SetUsageAndHelpFunc(cmd, fss, cols)
	case FlagSetOptions:
		fs = cmd.PersistentFlags()
		if app.options != nil {
			typed.AddFlags(fs)
		}
	default:
		fs = cmd.Flags()
	}

	version.AddFlags(fs)

	if !app.noConfig {
		AddConfigFlag(fs, app.name, app.watch)
	}

	app.cmd = cmd
}

// Run is used to launch the application.
func (app *App) Run() {
	os.Exit(cli.Run(app.cmd))
}

func (app *App) runCommand(cmd *cobra.Command, args []string) error {
	// display application version information
	version.PrintAndExitIfRequested()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if app.options != nil {
		if err := viper.Unmarshal(app.options); err != nil {
			return err
		}

		if complete, ok := app.options.(interface{ Complete() error }); ok {
			if err := complete.Complete(); err != nil {
				return err
			}
		}

		if validate, ok := app.options.(interface{ Validate() error }); ok {
			if err := validate.Validate(); err != nil {
				return err
			}
		}
	}

	app.initializeLogger()

	if !app.silence {
		log.Infow("starting application", "name", app.name, "version", version.Get().ToJSON())
		log.Infow("golang settings", "GOGC", os.Getenv("GOGC"), "GOMAXPROCS", os.Getenv("GOMAXPROCS"), "GOTRACEBACK", os.Getenv("GOTRACEBACK"))
		if !app.noConfig {
			PrintConfig()
		} else if app.options != nil {
			cliflag.PrintFlags(cmd.Flags())
		}
	}

	// 如果提供了组件运行器，则使用它来运行应用
	if app.componentRunner != nil {
		// 创建容器并传递给组件运行器
		var c *container.Container
		c = container.New(app.name)

		// 如果配置中包含ShutdownOverallTimeout，则设置容器选项
		if cfg, ok := app.options.(*config.Config); ok {
			c = container.New(app.name, container.WithShutdownOverallTimeout(cfg.ShutdownOverallTimeout))
		}

		return app.componentRunner.RunWithComponents(c)
	}

	if app.healthCheckFunc != nil {
		if err := app.healthCheckFunc(); err != nil {
			return err
		}
	}

	// run application
	return app.run()
}

// Command returns cobra command instance inside the application.
func (app *App) Command() *cobra.Command {
	return app.cmd
}

// formatBaseName is formatted as an executable file name under different
// operating systems according to the given name.
func formatBaseName(name string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		name = strings.ToLower(name)
		name = strings.TrimSuffix(name, ".exe")
	}
	return name
}

// initializeLogger sets up the logging system based on the configuration.
func (app *App) initializeLogger() {
	logOptions := log.NewOptions()

	// Configure logging options from viper
	if viper.IsSet("log.disable-caller") {
		logOptions.DisableCaller = viper.GetBool("log.disable-caller")
	}
	if viper.IsSet("log.disable-stacktrace") {
		logOptions.DisableStacktrace = viper.GetBool("log.disable-stacktrace")
	}
	if viper.IsSet("log.level") {
		logOptions.Level = viper.GetString("log.level")
	}
	if viper.IsSet("log.format") {
		logOptions.Format = viper.GetString("log.format")
	}
	if viper.IsSet("log.output-paths") {
		logOptions.OutputPaths = viper.GetStringSlice("log.output-paths")
	}

	// Initialize logging with custom context extractors
	log.Init(logOptions, log.WithContextExtractor(app.contextExtractors))
}

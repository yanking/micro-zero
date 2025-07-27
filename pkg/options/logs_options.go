package options

import (
	"github.com/spf13/pflag"
	"github.com/yanking/micro-zero/pkg/log"
)

var _ IOptions = (*LogsOptions)(nil)

// LogsOptions contains configuration items related to log.
type LogsOptions struct {
	// DisableCaller specifies whether to include caller information in the log.
	DisableCaller bool `json:"disable-caller,omitempty" mapstructure:"disable-caller"`
	// DisableStacktrace specifies whether to record a stack trace for all messages at or above panic level.
	DisableStacktrace bool `json:"disable-stacktrace,omitempty" mapstructure:"disable-stacktrace"`
	// EnableColor specifies whether to output colored logs.
	EnableColor bool `json:"enable-color"       mapstructure:"enable-color"`
	// Level specifies the minimum log level. Valid values are: debug, info, warn, error, dpanic, panic, and fatal.
	Level string `json:"level,omitempty" mapstructure:"level"`
	// Format specifies the log output format. Valid values are: console and json.
	Format string `json:"format,omitempty" mapstructure:"format"`
	// OutputPaths specifies the output paths for the logs.
	OutputPaths []string `json:"output-paths,omitempty" mapstructure:"output-paths"`
}

// NewLogsOptions creates an Options object with default parameters.
func NewLogsOptions() *LogsOptions {
	return &LogsOptions{
		Level:       "info",
		Format:      "console",
		OutputPaths: []string{"stdout"},
	}
}

// Validate verifies flags passed to LogsOptions.
func (o *LogsOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds command line flags for the configuration.
func (o *LogsOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.Level, "log.level", o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, "log.disable-caller", o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, "log.disable-stacktrace", o.DisableStacktrace, ""+
		"Disable the log to record a stack trace for all messages at or above panic level.")
	fs.BoolVar(&o.EnableColor, "log.enable-color", o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringVar(&o.Format, "log.format", o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.StringSliceVar(&o.OutputPaths, "log.output-paths", o.OutputPaths, "Output paths of log.")
}

// NewLog create log  with the given config.
func (o *LogsOptions) NewLog() (log.Logger, error) {
	opts := &log.Options{
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		EnableColor:       o.EnableColor,
		Level:             o.Level,
		Format:            o.Format,
		OutputPaths:       o.OutputPaths,
	}
	log.Init(opts)

	return log.Default(), nil
}

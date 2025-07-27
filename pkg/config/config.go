package config

import (
	"errors"
	"fmt"
	"github.com/yanking/micro-zero/internal/pkg/known"
	"github.com/yanking/micro-zero/pkg/contract"
	genericoptions "github.com/yanking/micro-zero/pkg/options"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	cliflag "k8s.io/component-base/cli/flag"
	"time"
)

var availableServerModes = sets.NewString(
	known.GinServerMode,
	known.GRPCServerMode,
	known.GRPCGatewayServerMode,
)

var _ contract.NamedFlagSetOptions = (*Config)(nil)

// Config 配置结构体，用于存储应用相关的配置.
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项.
type Config struct {
	// ServerMode 定义服务器模式：gRPC、Gin HTTP、HTTP Reverse Proxy.
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`
	// JWTKey 定义 JWT 密钥.
	JWTKey string `json:"jwt-key" mapstructure:"jwt-key"`
	// ShutdownOverallTimeout 优雅关闭超时时间.
	ShutdownOverallTimeout time.Duration `json:"shutdown-overall-timeout" mapstructure:"shutdown-overall-timeout"`
	// LogsOptions 定义日志配置选项.
	LogsOptions *genericoptions.LogsOptions `json:"logs" mapstructure:"logs"`
	// Expiration 定义 JWT Token 的过期时间.
	Expiration time.Duration `json:"expiration" mapstructure:"expiration"`
	// HTTPOptions 包含 HTTP 配置选项.
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
	// GRPCOptions 包含 gRPC 配置选项.
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`
	// MySQLOptions 包含 MySQL 配置选项.
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
}

// New 创建带有默认值的 Config 实例.
func New() *Config {
	opts := &Config{
		ServerMode:             known.GRPCGatewayServerMode,
		JWTKey:                 "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:             2 * time.Hour,
		ShutdownOverallTimeout: 20 * time.Second,
		LogsOptions:            genericoptions.NewLogsOptions(),
		HTTPOptions:            genericoptions.NewHTTPOptions(),
		GRPCOptions:            genericoptions.NewGRPCOptions(),
		MySQLOptions:           genericoptions.NewMySQLOptions(),
	}
	opts.HTTPOptions.Addr = ":5555"
	opts.GRPCOptions.Addr = ":6666"

	return opts
}

func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	fss.FlagSet("global").StringVar(&c.ServerMode, "server-mode", c.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fss.FlagSet("global").StringVar(&c.JWTKey, "jwt-key", c.JWTKey, "JWT signing key. Must be at least 6 characters long.")
	fss.FlagSet("global").DurationVar(&c.Expiration, "expiration", c.Expiration, "The expiration duration of JWT tokens.")

	c.LogsOptions.AddFlags(fss.FlagSet("logs"))
	c.HTTPOptions.AddFlags(fss.FlagSet("HTTP"))
	c.GRPCOptions.AddFlags(fss.FlagSet("gRPC"))
	c.MySQLOptions.AddFlags(fss.FlagSet("MySQL"))

	return fss
}

// Validate 校验 ServerOptions 中的选项是否合法.
func (c *Config) Validate() error {
	var errs []error

	// 校验 ServerMode 是否有效
	if !availableServerModes.Has(c.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: must be one of %v", availableServerModes.UnsortedList()))
	}
	// 校验 JWTKey 长度
	if len(c.JWTKey) < 6 {
		errs = append(errs, errors.New("JWTKey must be at least 6 characters long"))
	}

	// 校验子选项
	errs = append(errs, c.LogsOptions.Validate()...)
	errs = append(errs, c.HTTPOptions.Validate()...)
	errs = append(errs, c.MySQLOptions.Validate()...)

	// 如果是 gRPC 或 gRPC-Gateway 模式，校验 gRPC 配置
	if c.ServerMode == known.GRPCServerMode || c.ServerMode == known.GRPCGatewayServerMode {
		errs = append(errs, c.GRPCOptions.Validate()...)
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Config) Complete() error {
	return nil
}

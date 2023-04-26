package log

import (
	"io"
	"strings"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/provider/log/formatter"
	"github.com/chsir-zy/anan/framework/provider/log/services"
)

type AnanLogServiceProvider struct {
	framework.ServiceProvider

	Driver string

	level contract.LogLevel

	Formatter contract.Formatter

	CtxFielder contract.CtxFielder

	Output io.Writer
}

func (lp *AnanLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if lp.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			// 默认使用console
			return services.NewAnanConsoleLog
		}

		cs := tcs.(contract.Config)
		lp.Driver = strings.ToLower(cs.GetString("log.Driver"))
	}

	switch lp.Driver {
	case "single":
		return services.NewAnanSingleLog
	case "rotate":
		return services.NewAnanRotateLog
	case "console":
		return services.NewAnanConsoleLog
	case "custom":
		return services.NewAnanCustomLog
	default:
		return services.NewAnanConsoleLog
	}
}

func (lp *AnanLogServiceProvider) Boot(c framework.Container) error {
	return nil
}

func (lp *AnanLogServiceProvider) IsDefer() bool {
	return false
}

func (lp *AnanLogServiceProvider) Params(c framework.Container) []interface{} {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	if lp.Formatter == nil {
		lp.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			c := configService.GetString("log.formatter")
			if c == "json" {
				lp.Formatter = formatter.JsonFormatter
			}
		}
	}

	if lp.level == contract.UnknowLevel {
		lp.level = contract.InfoLevel
		if configService.IsExist("log.level") {
			lp.level = logLevel(configService.GetString("log.level"))
		}
	}

	return []interface{}{c, lp.level, lp.CtxFielder, lp.Formatter, lp.Output}
}

func (lp *AnanLogServiceProvider) Name() string {
	return contract.LogKey
}

// logLevel get level from string
func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknowLevel
}

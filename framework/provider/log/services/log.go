package services

import (
	"context"
	"io"
	pkgLog "log"
	"time"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/provider/log/formatter"
)

type AnanLog struct {
	level      contract.LogLevel   // 日志级别
	formatter  contract.Formatter  // 日志格式化方法
	ctxFielder contract.CtxFielder // 日志上下文信息
	output     io.Writer
	c          framework.Container
}

// 判断这个级别是否可以打印
func (log *AnanLog) IsLevelEnable(level contract.LogLevel) bool {
	return level <= log.level
}

// 核心的日志打印文件
func (log *AnanLog) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	if !log.IsLevelEnable(level) { //先判断日志级别是否可以打印
		return nil
	}

	fs := fields
	if log.ctxFielder != nil {
		t := log.ctxFielder(ctx)
		for k, v := range t {
			fs[k] = v
		}
	}

	if log.c.IsBind(contract.TraceKey) {
		tracer := log.c.MustMake(contract.TraceKey).(contract.Trace)
		tc := tracer.GetTrace(ctx)
		if tc != nil {
			maps := tracer.ToMap(tc)
			for k, v := range maps {
				fs[k] = v
			}
		}
	}

	if log.formatter == nil {
		log.formatter = formatter.TextFormatter
	}

	ct, err := log.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}

	if level == contract.PanicLevel {
		pkgLog.Panicln(ct)
		return nil
	}

	log.output.Write(ct)
	log.output.Write([]byte("\r\n"))

	return nil
}

func (log *AnanLog) SetOutput(output io.Writer) {
	log.output = output
}

func (log *AnanLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.PanicLevel, ctx, msg, fields)
}

func (log *AnanLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.FatalLevel, ctx, msg, fields)
}

func (log *AnanLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.ErrorLevel, ctx, msg, fields)
}

func (log *AnanLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.WarnLevel, ctx, msg, fields)
}

func (log *AnanLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.InfoLevel, ctx, msg, fields)
}

func (log *AnanLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.DebugLevel, ctx, msg, fields)
}

func (log *AnanLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.TraceLevel, ctx, msg, fields)
}

func (log *AnanLog) SetLevel(level contract.LogLevel) {
	log.level = level
}

func (log *AnanLog) SetCtxFielder(hanlder contract.CtxFielder) {
	log.ctxFielder = hanlder
}

func (log *AnanLog) SetFormatter(formatter contract.Formatter) {
	log.formatter = formatter
}

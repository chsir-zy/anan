package services

import (
	"io"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

type AnanCustomLog struct {
	AnanLog
}

// 自定义输出地方的日志
func NewAnanCustomLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)

	log := &AnanCustomLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 将内容输出到自定义的地方
	log.SetOutput(output)
	log.c = c
	return log, nil
}

package services

import (
	"os"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
)

type AnanConsoleLog struct {
	AnanLog
}

func NewAnanConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &AnanConsoleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 将内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}

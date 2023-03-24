package formatter

import "github.com/chsir-zy/anan/framework/contract"

func Prefix(level contract.LogLevel) string {
	prefix := ""
	switch level {
	case contract.PanicLevel:
		prefix = "[panix]"
	case contract.FatalLevel:
		prefix = "[fatal]"
	case contract.ErrorLevel:
		prefix = "[error]"
	case contract.WarnLevel:
		prefix = "[warn]"
	case contract.InfoLevel:
		prefix = "[info]"
	case contract.DebugLevel:
		prefix = "[debug]"
	case contract.TraceLevel:
		prefix = "[trace]"
	}
	return prefix
}

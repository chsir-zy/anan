package formatter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/chsir-zy/anan/framework/contract"
)

func TextFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	Separator := "\t"

	prefix := Prefix(level)
	bf.WriteString(prefix)
	bf.WriteString(Separator)

	// 时间
	time := t.Format(time.RFC3339)
	bf.WriteString(time)
	bf.WriteString(Separator)

	//信息
	bf.WriteString("\"")
	bf.WriteString(msg)
	bf.WriteString("\"")
	bf.WriteString(Separator)

	bf.WriteString(fmt.Sprint(fields))

	return bf.Bytes(), nil
}

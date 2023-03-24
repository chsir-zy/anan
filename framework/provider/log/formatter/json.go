package formatter

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/chsir-zy/anan/framework/contract"
	"github.com/pkg/errors"
)

func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	fields["msg"] = msg
	fields["level"] = level
	fields["timestamp"] = t.Format(time.RFC3339)

	b, err := json.Marshal(fields)
	if err != nil {
		return nil, errors.Wrap(err, "json formatter error")
	}

	bf.Write(b)
	return bf.Bytes(), nil
}

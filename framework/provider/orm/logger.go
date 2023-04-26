package orm

import (
	"context"
	"time"

	"github.com/chsir-zy/anan/framework/contract"
	"gorm.io/gorm/logger"
)

// OrmLogger orm的日志实现类, 实现了gorm.Logger.Interface
type OrmLogger struct {
	logger contract.Log
}

func NewOrmLogger(logger contract.Log) *OrmLogger {
	return &OrmLogger{logger: logger}
}

func (ol *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
	return ol
}

func (ol *OrmLogger) Info(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	ol.logger.Info(ctx, s, fields)
}

func (ol *OrmLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	ol.logger.Warn(ctx, s, fields)
}

func (ol *OrmLogger) Error(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	ol.logger.Error(ctx, s, fields)
}

func (ol *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	fields := map[string]interface{}{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}

	s := "orm trace sql"
	ol.logger.Trace(ctx, s, fields)
}

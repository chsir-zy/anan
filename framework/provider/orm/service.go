package orm

import (
	"context"
	"sync"
	"time"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// AnanOrm 代表anan框架的orm实现
type AnanGorm struct {
	container framework.Container
	dbs       map[string]*gorm.DB //确保服务只实现一次 单例模式  key是DSN

	lock *sync.RWMutex
}

func NewAnanGorm(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}

	return &AnanGorm{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil

}

func (a *AnanGorm) GetDB(option ...contract.DBoption) (*gorm.DB, error) {
	// 读取默认配置
	config := GetBaseConfig(a.container)

	logService := a.container.MustMake(contract.LogKey).(contract.Log)

	// 设置logger
	logger := NewOrmLogger(logService)
	config.Config = &gorm.Config{
		Logger: logger,
	}

	// option 进行修改
	for _, opt := range option {
		if err := opt(a.container, config); err != nil {
			return nil, err
		}
	}

	// 如果没有设置dsn 则生成dsn
	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	// 判断是否已经实例话了dsn
	a.lock.RLock()
	if db, ok := a.dbs[config.Dsn]; ok {
		return db, nil
	}
	a.lock.RUnlock()

	// 没有实例化 *gorm.DB
	a.lock.Lock()
	defer a.lock.Unlock()

	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config)
	}

	// 设置对应的连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}
	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}
	if config.ConnMaxLifetime != "" {
		liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max lift time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(liftTime)
		}
	}
	if config.ConnMaxIdletime != "" {
		idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
		if err != nil {
			logger.Error(context.Background(), "conn max idle time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxIdleTime(idleTime)
		}
	}

	if err != nil {
		a.dbs[config.Dsn] = db
	}

	return db, nil
}

package contract

import (
	"net"
	"strconv"
	"time"

	"github.com/chsir-zy/anan/framework"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const ORMKey = "anan:orm"

type ORMService interface {
	GetDB(option ...DBoption) (*gorm.DB, error)
}

type DBoption func(container framework.Container, config *DBConfig) error

type DBConfig struct {
	WriteTimeout string `yaml:"write_timeout"` // 写超时时间
	Loc          string `yaml:"log"`           //时区
	Port         int    `yaml:"port"`          //端口
	ReadTimeout  string `yaml:"read_timeout"`  //读超时时间
	Charset      string `yaml:"charset"`       //字符集
	ParseTime    bool   `yaml:"parse_time"`    //是否转换时间类型
	Protocol     string `yaml:"protocol"`      //传输协议
	Dsn          string `yaml:"dsn"`           //
	Database     string `yaml:"database"`
	Collation    string `yaml:"collation"` // 字符序
	Timeout      string `yaml:"timeout"`   //连接超时时间
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`

	// 以下配置关于连接池
	ConnMaxIdle     int    `yaml:"conn_max_idle"`      // 最大空闲连接数
	ConnMaxOpen     int    `yaml:"conn_max_open"`      // 最大连接数
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`  // 连接最大生命周期
	ConnMaxIdletime string `yaml:"conn_max_idle_time"` // 空闲最大生命周期

	//配置gorm
	*gorm.Config
}

func (conf *DBConfig) FormatDsn() (string, error) {
	port := strconv.Itoa(conf.Port)
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return "", err
	}
	readTimeout, err := time.ParseDuration(conf.ReadTimeout)
	if err != nil {
		return "", err
	}
	writeTimeout, err := time.ParseDuration(conf.WriteTimeout)
	if err != nil {
		return "", err
	}
	location, err := time.LoadLocation(conf.Loc)
	if err != nil {
		return "", err
	}

	driverConf := mysql.Config{
		User:                 conf.Username,
		Passwd:               conf.Password,
		Net:                  conf.Protocol,
		Addr:                 net.JoinHostPort(conf.Host, port),
		DBName:               conf.Database,
		Collation:            conf.Collation,
		Loc:                  location,
		Timeout:              timeout,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
		ParseTime:            conf.ParseTime,
		AllowNativePasswords: true,
	}

	return driverConf.FormatDSN(), nil
}

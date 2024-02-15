package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (m Mysql) Init() (*gorm.DB, error) {
	if m.Enable == false {
		return nil, nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置

	}
	db, err := gorm.Open(mysql.New(mysqlConfig))
	return db, err
}

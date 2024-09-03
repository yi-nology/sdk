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
	if err != nil {
		return nil, err
	}
	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"test/config"
	"test/pkg/logger"
	"time"
)

type MySqlWrapper struct {
	DB *gorm.DB
}

func NewMysqlWrapper() (m *MySqlWrapper, err error) {
	MysqlDB, err := InitMysql()
	if err != nil {
		return nil, err
	}
	return &MySqlWrapper{DB: MysqlDB}, nil
}
func (m *MySqlWrapper) Destroy() error {
	if m.DB == nil {
		return nil
	}
	sqlDB, err := m.DB.DB()
	if err != nil {
		logger.Error("销毁时获取 MySQL SQL 实例失败")
		return err
	}
	// 关闭连接池（优雅销毁）
	if err := sqlDB.Close(); err != nil {
		logger.Error("关闭 MySQL 连接池失败")
		return err
	}
	logger.Info("MySQL 连接池已优雅关闭")
	return nil
}

func InitMysql() (*gorm.DB, error) {
	gConfig := config.GConfig

	fmt.Println(gConfig)
	if gConfig.Mysql.Addr == "" || gConfig.Mysql.UserName == "" || gConfig.Mysql.Password == "" {
		logger.Error("mysql数据库配置缺少")
		return nil, fmt.Errorf("数据库配置不全")
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		gConfig.Mysql.UserName,
		gConfig.Mysql.Password,
		gConfig.Mysql.Addr,
		gConfig.Mysql.Port,
		gConfig.Mysql.Name,
	)
	//建立连接
	//var err error
	MysqlDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		logger.Error("mysql连接错误...")
		panic("mysql连接错误...")
	}
	//配置连接池
	sqlDB, err := MysqlDB.DB()
	if err != nil {
		panic("获取sql实例失败")
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	logger.Info("mysql连接成功！！！")
	return MysqlDB, nil
}

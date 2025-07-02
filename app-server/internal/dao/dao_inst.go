package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"stay-server/internal/config"
	"stay-server/utils"
)

var DbDao *gorm.DB

type DaoInstance struct {
	Id int32
	//DbDao *gorm.DB
	//DbConfig DbConfig
}

func NewDaoInstance(id int32) *DaoInstance {
	var err error = nil
	var daoInst *DaoInstance = &DaoInstance{Id: id}
	//if err = daoInst.readDatabaseConfig(); err != nil {
	//	panic(fmt.Sprintf("failure read config: %v", err))
	//}
	DbDao, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			//daoInst.DbConfig.MysqlConfig.Username,
			config.AppCfg.MysqlConfig.Username,
			config.AppCfg.MysqlConfig.Password,
			config.AppCfg.MysqlConfig.Protocol,
			config.AppCfg.MysqlConfig.Host,
			config.AppCfg.MysqlConfig.Port,
			config.AppCfg.MysqlConfig.Database),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open database: %v", err))
	}
	if DbDao.Exec(`SELECT 1 + 1;`).Error != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	utils.Logger{}.PrintSuccess("数据库服务器已连接 Server is connected and tested.")
	//log.Println("Server is connected and tested.")
	daoInst.migrateTables()
	return daoInst
}

//func (d *DaoInstance) readDatabaseConfig() error {
//	file, err := os.ReadFile("config/config.yaml")
//	if err != nil {
//		return fmt.Errorf("read yaml file error: %w", err)
//	}
//	var cfg DbConfig
//	if err := yaml.Unmarshal(file, &cfg); err != nil {
//		return fmt.Errorf("unmarshal yaml error: %w", err)
//	}
//	d.DbConfig = cfg
//	return nil
//}

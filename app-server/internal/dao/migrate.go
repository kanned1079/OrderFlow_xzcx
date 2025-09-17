package dao

import (
	"os"
	"stay-server/internal/models"
	"stay-server/utils"
)

func (this *DaoInstance) migrateTables() {
	var log utils.Logger
	log.PrintInfo("Run db migrate.")

	signFile := "./db_migrated_sign"

	// 判斷標記文件是否存在
	if _, err := os.Stat(signFile); err == nil {
		log.PrintSuccess("Skip migrate.")
		return
	}

	// 正常遷移數據表
	if err := DbDao.AutoMigrate(
		&models.User{},
		&models.Merchant{},
		&models.Goods{},
		&models.Order{},
		&models.OrderItem{},
		&models.Comment{},
		&models.Category{},
		&models.Address{},
	); err != nil {
		log.PrintError("Migrate failed: ", err)
		return
	}

	// 創建標記文件
	if f, err := os.Create(signFile); err == nil {
		f.Close()
	} else {
		log.PrintError("Failed to create sign file: %v", err)
		return
	}

	log.PrintSuccess("Migrate finished and sign file created.")
}

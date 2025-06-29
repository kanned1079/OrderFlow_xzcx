package dao

import "stay-server/internal/models"

func (this *DaoInstance) migrateTables() {
	DbDao.AutoMigrate(&models.User{})
	DbDao.AutoMigrate(&models.Merchant{})
	DbDao.AutoMigrate(&models.Goods{})
}

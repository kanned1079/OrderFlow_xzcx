package dao

import "stay-server/internal/models"

func (this *DaoInstance) migrateTables() {
	DbDao.AutoMigrate(&models.User{})
	DbDao.AutoMigrate(&models.Merchant{})
	DbDao.AutoMigrate(&models.Goods{})
	DbDao.AutoMigrate(&models.Order{})
	DbDao.AutoMigrate(&models.OrderItem{})
	DbDao.AutoMigrate(&models.Comment{})
	DbDao.AutoMigrate(&models.Category{})
	DbDao.AutoMigrate(&models.Address{})
}

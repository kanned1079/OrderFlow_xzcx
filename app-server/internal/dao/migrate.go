package dao

import "stay-server/internal/models"

func (this *DaoInstance) migrateTables() {
	DbDao.AutoMigrate(&models.User{})
}

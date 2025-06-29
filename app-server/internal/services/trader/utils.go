package trader

import (
	"stay-server/internal/dao"
	"stay-server/internal/models"
)

// categoryExistsForMerchant 用于在新增或者修改商品时判定商品分类是否可用
func (TraderServices) categoryExistsForMerchant(merchantId, categoryId int64) (bool, error) {
	var count int64
	err := dao.DbDao.Model(&models.Category{}).
		Where("id = ? AND merchant_id = ?", categoryId, merchantId).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

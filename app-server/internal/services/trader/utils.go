package trader

import (
	"stay-server/internal/dao"
	"stay-server/internal/models"
)

// categoryExistsForMerchant 用于在新增或者修改商品时判定商品分类是否可用
func (this *TraderServices) categoryExistsForMerchant(categoryId, merchantId int64) (bool, error) {
	//this.utils.Logger.PrintInfo("c_id: ", categoryId, " m_id: ", merchantId)
	var count int64
	err := dao.DbDao.Model(&models.Category{}).
		Where("id = ? AND merchant_id = ?", categoryId, merchantId).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

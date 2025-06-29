package trader

import (
	"github.com/gin-gonic/gin"
	"stay-server/utils"
)

// TraderServicesInterface 商户服务接口
type TraderServicesInterface interface {
	utils.Utils
	// Service层处理方法
	GetGoodsList(ctx *gin.Context)
	AddNewGoods(ctx *gin.Context)
	EditGoodsInfo(ctx *gin.Context)
	DeleteGoods(ctx *gin.Context)
	GetCategoryList(ctx *gin.Context)
	AddNewCategory(ctx *gin.Context)
	EditCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
	// 私有方法
	categoryExistsForMerchant(merchantId, categoryId int64) (bool, error)
}

type TraderServices struct {
	utils utils.Utils
}

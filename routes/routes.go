package routes

import (
	"funding-watch/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")
	{
		funding := apiV1.Group("/funding")
		{
			funding.POST("tx", controller.GetFundingRecordsbyTxHash)
			funding.GET("", controller.GetFundingRecords)
			funding.GET("ranking", controller.GetFundRecordsRanking)
			funding.GET("sender/:sender", controller.GetFundRecordsBySender)
		}
	}
}

package controller

import (
	"funding-watch/service"

	"github.com/gin-gonic/gin"
)

func getFundingRecords(c *gin.Context) {
	records, err := service.GetFundingRecords(c.Query("tx_hash"))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to get funding records",
			"error":   err.Error(),
		})
		return
		c.JSON(200, gin.H{"data": records})
	}
}

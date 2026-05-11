package controller

import (
	"fmt"
	"funding-watch/models"
	"funding-watch/service"

	"github.com/gin-gonic/gin"
)

func GetFundingRecords(c *gin.Context) {
	var req models.TxHashRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}
	records, err := service.GetFundingRecords(req.TxHash)
	fmt.Printf("tx_hash:%s", req.TxHash)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to get funding records",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"data": records})
}

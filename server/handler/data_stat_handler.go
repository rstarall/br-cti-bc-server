package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/fabric"
)

// 查询情报统计信息
func QueryCTISummaryInfo(c *gin.Context) {
	// 解析请求参数
	var params struct {
		Limit int `json:"limit"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryLatestCTISummaryInfo(params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 获取数据统计信息
func GetDataStatistics(c *gin.Context) {
	resp, err := fabric.GetDataStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 获取情报交易趋势数据(Post)
func GetUpchainTrend(c *gin.Context) {
	var params struct {
		TimeRange string `json:"time_range"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.GetUpchainTrend(params.TimeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 获取攻击类型排行(Post)
func GetAttackTypeRanking(c *gin.Context) {
	resp, err := fabric.GetAttackTypeRanking()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 获取IOCs类型分布(Post)
func GetIOCsDistribution(c *gin.Context) {
	resp, err := fabric.GetIOCsDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 获取全球IOCs地理分布(Post)
func GetGlobalIOCsDistribution(c *gin.Context) {
	resp, err := fabric.GetGlobalIOCsDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 获取系统概览数据(Post)
func GetSystemOverview(c *gin.Context) {
	resp, err := fabric.GetSystemOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

package handler

import (
	"github.com/gin-gonic/gin"
	ipfsService "github.com/righstar2020/br-cti-bc-server/service"
	"net/http"
)

// 查询IOC地理分布
func QueryIOCGeoDistribution(c *gin.Context) {


	// 获取统计数据
	stats, err := ipfsService.GetIPFSServiceInstance().GetIOCWorldMapStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "获取统计数据失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": stats})
}

// 查询IOC类型分布
func QueryIOCTypeDistribution(c *gin.Context) {

	stats, err := ipfsService.GetIPFSServiceInstance().GetIOCTypeStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "查询IOC类型分布失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": stats})
}

// 查询攻击类型统计
func QueryAttackTypeStatistics(c *gin.Context) {


	stats, err := ipfsService.GetIPFSServiceInstance().GetAttackTypeStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "查询攻击类型统计失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": stats})
}

// 查询攻击IOC信息
func QueryAttackIOCInfo(c *gin.Context) {


	infos, err := ipfsService.GetIPFSServiceInstance().GetAttackIOCInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "查询攻击IOC信息失败",
			"detail": err.Error(),
		})
		return
	}
	//最新返回10条数据
	infos = infos[:10]

	c.JSON(http.StatusOK, gin.H{"result": infos})
}

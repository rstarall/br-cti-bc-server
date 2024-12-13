package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
	ipfsservice "github.com/righstar2020/br-cti-bc-server/service"
	"net/http"
)

var nodeAddrs = []string{
	"http://127.0.0.1:8080",
	// "https://ipfs.io",
	// "https://dweb.link",
	// "https://gateway.pinata.cloud",
}

// 查询IPFS文件内容
func GetIPFSContent(c *gin.Context) {
	// 解析请求参数
	var params struct {
		Hash string `json:"hash" binding:"required" form:"hash"`
	}

	// 尝试获取POST JSON参数
	if err := c.ShouldBindJSON(&params); err != nil {
		// JSON解析失败,尝试获取GET/POST表单参数
		if err := c.ShouldBind(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "参数错误",
				"detail": err.Error(),
			})
			return
		}
	}
	// 调用IPFS服务获取内容
	url, content, err := ipfs.GetIPFSContentWithFallback(params.Hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "获取IPFS内容失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"content": string(content),
	})
}

func GetIPFSFileUrl(c *gin.Context) {
	// 解析请求参数
	var params struct {
		Hash string `json:"hash"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"url": ipfs.GetIPFSServerHost() + "/ipfs/",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": ipfs.GetIPFSServerHost() + "/ipfs/" + params.Hash,
	})
}

func ProcessIOCWorldMapStatistics(c *gin.Context) {
	// 获取请求体中的 IOC 数据
	var params struct {
		IOCData string `json:"ioc_data" binding:"required"`
	}

	// 解析请求参数
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "参数错误",
			"detail": err.Error(),
		})
		return
	}

	// 调用 IPFSService 处理 IOC 数据
	stats, err := ipfsservice.NewIPFSService(nodeAddrs).ProcessIOCWorldMapStatistics(params.IOCData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "处理 IOC 数据失败",
			"detail": err.Error(),
		})
		return
	}

	// 返回生成的统计信息
	c.JSON(http.StatusOK, gin.H{
		"statistics": stats,
	})
}

// SaveIOCWorldMapStatisticsHandler 处理保存统计数据并提供下载功能的 Handler
func SaveIOCWorldMapStatisticsHandler(c *gin.Context) {
	// 获取请求体中的 IOC 数据
	var params struct {
		IOCData string `json:"ioc_data" binding:"required"`
	}

	// 解析请求参数
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "参数错误",
			"detail": err.Error(),
		})
		return
	}

	// 调用 ProcessIOCWorldMapStatistics 处理数据
	statsMap, err := ipfsservice.NewIPFSService(nodeAddrs).ProcessIOCWorldMapStatistics(params.IOCData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "处理 IOC 数据失败",
			"detail": err.Error(),
		})
		return
	}

	// 将处理后的 StatsMap 转换为 JSON
	statsJSON, err := json.MarshalIndent(statsMap, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "生成统计 JSON 失败",
			"detail": err.Error(),
		})
		return
	}

	// 设置响应头，告知浏览器文件是作为附件下载
	c.Header("Content-Disposition", "attachment; filename=ioc_world_map_statistics.json")
	c.Header("Content-Type", "application/json")

	// 直接返回 JSON 数据，触发文件下载
	c.Data(http.StatusOK, "application/json", statsJSON)
}

// 获取 IOC 世界地图统计数据
//func GetIOCWorldMapStatistics(c *gin.Context) {
//	// 调用 IPFSService 获取 IOC 世界地图统计数据
//	stats, err := ipfsservice.NewIPFSService(nodeAddrs).GetIOCWorldMapStatistics()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error":  "获取统计数据失败",
//			"detail": err.Error(),
//		})
//		return
//	}
//
//	// 返回统计数据
//	c.JSON(http.StatusOK, gin.H{
//		"statistics": stats,
//	})
//}

// DownloadFileHandler 处理从 IPFS 下载文件并提供给浏览器下载
func DownloadFileHandler(c *gin.Context) {
	var params struct {
		Hash string `json:"hash" binding:"required" form:"hash"`
	}

	// 尝试获取POST JSON参数
	if err := c.ShouldBindJSON(&params); err != nil {
		// JSON解析失败, 尝试获取GET/POST表单参数
		if err := c.ShouldBind(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "参数错误",
				"detail": err.Error(),
			})
			return
		}
	}

	// 从 IPFS 获取文件内容
	_, content, err := ipfs.GetIPFSContentFromNode(params.Hash, "http://127.0.0.1:8080")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "从 IPFS 获取文件失败",
			"detail": err.Error(),
		})
		return
	}

	// 设置响应头，告诉浏览器这是一个附件文件，需要下载
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", params.Hash))
	c.Header("Content-Type", "application/json")

	// 直接返回文件内容，触发浏览器下载
	c.Data(http.StatusOK, "application/json", content)
}

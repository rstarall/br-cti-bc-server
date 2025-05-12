package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
	ipfsService "github.com/righstar2020/br-cti-bc-server/service"
	"github.com/righstar2020/br-cti-bc-server/util"
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
	ipfsService := ipfsService.GetIPFSServiceInstance()
	stats, err := ipfsService.ProcessIOCWorldMapStatistics(params.IOCData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "处理 IOC 数据失败",
			"detail": err.Error(),
		})
		return
	}

	// 返回生成的统计信息和保存路径
	c.JSON(http.StatusOK, gin.H{"result": stats})
}

func GetIOCWorldMapStatisticsHandler(c *gin.Context) {
	// 获取统计数据
	stats, err := ipfsService.GetIPFSServiceInstance().GetIOCWorldMapStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "获取统计数据失败",
			"detail": err.Error(),
		})
		return
	}

	// 返回统计数据
	c.JSON(http.StatusOK, gin.H{"result": stats})
}

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

// DownloadFileToPathHandler 处理从 IPFS 下载文件到指定本地路径
func DownloadFileToPathHandler(c *gin.Context) {
	var params struct {
		Hash     string `json:"hash" binding:"required"`
		SavePath string `json:"save_path" binding:"required"`
	}

	// 尝试获取POST JSON参数
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "参数错误",
			"detail": err.Error(),
		})
		return
	}

	// 调用IPFS服务获取文件内容
	_, content, err := ipfsService.GetIPFSServiceInstance().DownloadIPFSFile(params.Hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "下载IPFS文件失败",
			"detail": err.Error(),
		})
		return
	}

	// 确保保存路径的目录存在
	saveDir := params.SavePath
	if !util.PathExists(saveDir) {
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "创建保存目录失败",
				"detail": err.Error(),
			})
			return
		}
	}

	// 生成文件名（使用哈希值作为文件名）
	fileName := params.Hash
	if !strings.HasSuffix(fileName, ".txt") {
		fileName += ".txt"
	}

	// 构建完整的文件路径
	filePath := filepath.Join(params.SavePath, fileName)

	// 写入文件
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "保存文件失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "文件下载成功",
		"file_path": filePath,
	})
}

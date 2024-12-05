

package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
)

// 查询IPFS文件内容
func GetIPFSContent(c *gin.Context) {
	// 解析请求参数
	var params struct {
		Hash string `json:"hash" binding:"required"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
			"detail": err.Error(),
		})
		return
	}

	// 调用IPFS服务获取内容
	url,content, err := ipfs.GetIPFSContentWithFallback(params.Hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取IPFS内容失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":url,
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
			"url":ipfs.GetIPFSServerHost()+"/ipfs/",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": ipfs.GetIPFSServerHost() + "/ipfs/" + params.Hash,
	})
}

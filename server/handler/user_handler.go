package handler

import (
	"github.com/gin-gonic/gin"
	fabric "github.com/righstar2020/br-cti-bc-server/fabric"
	"net/http"
)

func RegisterUserAccount(c *gin.Context) {
	// 从请求中获取参数
	var params struct {
		UserName  string `json:"user_name"`
		PublicKey string `json:"public_key"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.RegisterUserAccount(params.UserName, params.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

func QueryUserInfo(c *gin.Context) {
	// 从请求中获取参数
	var params struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryUserInfo(params.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

func GetUserStatistics(c *gin.Context) {
	// 从请求中获取参数
	var params struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	resp, err := fabric.GetUserStatistics(params.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户统计数据失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

func QueryPointTransactions(c *gin.Context) {
	// 从请求中获取参数
	var params struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	resp, err := fabric.QueryPointTransactions(params.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询积分交易记录失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

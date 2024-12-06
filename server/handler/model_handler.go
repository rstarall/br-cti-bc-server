package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/fabric"
)

// 注册模型信息(Post)
func RegisterModelInfo(c *gin.Context) {
	//// 解析请求参数
	var txRawMsg *fabric.TxMsgRawData

	if err := c.ShouldBindJSON(&txRawMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "参数错误",
			"detail": err.Error(),
		})
		log.Printf("参数错误: %s", err)
		return
	}

	// 序列化并打印日志
	txRawMsgData, err := json.Marshal(txRawMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "JSON序列化失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("序列化后的数据: %s", string(txRawMsgData))

	// 调用fabric注册Model信息
	resp, err := fabric.RegisterModelInfo(txRawMsgData)

	if err != nil {
		log.Printf("Fabric注册失败: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Fabric注册失败",
			"detail": err.Error(),
		})
		return
	}
	log.Printf("Fabric注册成功: %s", resp)
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 查询模型信息(Post)
func QueryModelInfo(c *gin.Context) {
	// 解析请求参数
	var params struct {
		ModelID string `json:"model_id"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryModelInfo(params.ModelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 分页查询模型(Post)
func QueryModelInfoWithPagination(c *gin.Context) {
	// 解析请求参数
	var params struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryModelInfoWithPagination(params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据类型分页查询模型(Post)
func QueryModelsByTypeWithPagination(c *gin.Context) {
	// 解析请求参数
	var params struct {
		ModelType  int    `json:"model_type"`
		Page       int    `json:"page"`
		PageSize   int    `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryModelsByTypeWithPagination(params.ModelType, params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据CTI ID查询模型(Post)
func QueryModelsByRefCTIId(c *gin.Context) {
	// 解析请求参数
	var params struct {
		CTIID string `json:"cti_id"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryModelsByRefCTIId(params.CTIID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据创建者ID查询模型(Post)
func QueryModelInfoByCreatorUserID(c *gin.Context) {
	// 解析请求参数
	var params struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryModelInfoByCreatorUserID(params.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

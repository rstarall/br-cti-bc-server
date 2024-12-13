package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/fabric"
	"log"
	"net/http"
)

// CTI注册接口(Post)
func RegisterCtiInfo(c *gin.Context) {
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

	// 调用fabric注册CTI信息
	resp, err := fabric.RegisterCtiInfo(txRawMsgData)

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

// CTI查询接口(Post)
func QueryCtiInfo(c *gin.Context) {
	// 解析请求参数
	ctiTxData := &fabric.CtiTxData{}
	if err := c.ShouldBindJSON(ctiTxData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := fabric.QueryCtiInfoByID(ctiTxData.CTIID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据Type和分页信息查询情报(Post)
func QueryCtiInfoByTypeWithPagination(c *gin.Context) {
	// 解析请求参数
	var params struct {
		CtiType  int `json:"cti_type"`
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryCtiInfoByTypeWithPagination(params.CtiType, params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据类型查询情报(Post)
func QueryCtiInfoByType(c *gin.Context) {
	// 解析请求参数
	var params struct {
		CtiType int `json:"cti_type"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryCtiInfoByType(params.CtiType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 分页查询所有情报信息(Post)
func QueryAllCtiInfoWithPagination(c *gin.Context) {
	var params struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryAllCtiInfoWithPagination(params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据情报哈希查询情报(Post)
func QueryCtiInfoByCTIHash(c *gin.Context) {
	var params struct {
		CTIHash string `json:"cti_hash"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryCtiInfoByCTIHash(params.CTIHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据创建者ID查询情报(Post)
func QueryCtiInfoByCreatorUserID(c *gin.Context) {
	var params struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryCtiInfoByCreatorUserID(params.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据Type和分页信息查询情报(路径参数)
func QueryCtiInfoByTypeWithParams(c *gin.Context) {
	// 从URL参数中获取参数
	ctiType := c.Query("type")
	page := c.Query("page")
	pageSize := c.Query("pageSize")

	// 转换参数类型
	ctiTypeInt := 0
	if ctiType != "" {
		fmt.Sscanf(ctiType, "%d", &ctiTypeInt)
	}

	pageInt := 1 // 默认值
	if page != "" {
		fmt.Sscanf(page, "%d", &pageInt)
	}

	pageSizeInt := 10 // 默认值
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &pageSizeInt)
	}

	resp, err := fabric.QueryCtiInfoByTypeWithPagination(ctiTypeInt, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

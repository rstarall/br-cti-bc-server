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
	//modelTxData := &fabric.ModelTxData{}
	//if err := c.ShouldBindJSON(modelTxData); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//resp, err := fabric.RegisterModelInfo(modelTxData)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//c.JSON(http.StatusOK, gin.H{"result": resp})
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

// 根据ID和分页信息查询模型(Post)
func QueryModelInfoByIDWithPagination(c *gin.Context) {
	// 解析请求参数
	var params struct {
		PageSize int    `json:"page_size"`
		Bookmark string `json:"bookmark"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryModelInfoByIDWithPagination(params.PageSize, params.Bookmark)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 根据流量类型查询模型(Post)
func QueryModelsByTrafficType(c *gin.Context) {
	// 解析请求参数
	var params struct {
		TrafficType string `json:"model_traffic_type"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryModelsByTrafficType(params.TrafficType)
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

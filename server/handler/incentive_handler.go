package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/fabric"
	"log"
	"net/http"
)

// 注册文档激励信息
func RegisterDocIncentiveInfo(c *gin.Context) {
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

	// 调用fabric注册文档激励信息
	resp, err := fabric.RegisterDocIncentiveInfo(txRawMsgData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "注册文档激励信息失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 查询文档激励信息
func QueryDocIncentiveInfo(c *gin.Context) {
	var params struct {
		RefID   string `json:"refen_id"`
		DocType string `json:"doctype"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryDocIncentiveInfo(params.RefID, params.DocType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

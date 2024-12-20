package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/fabric"
	"log"
	"net/http"
)

// 注册评论
func RegisterComment(c *gin.Context) {
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

	// 调用fabric注册评论
	resp, err := fabric.RegisterComment(txRawMsgData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "注册评论失败",
			"detail": err.Error(),
		})
		return
	}
	log.Printf("注册评论成功: %s", resp)
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 审核评论
func ApproveComment(c *gin.Context) {
	var txRawMsg *fabric.TxMsgRawData

	if err := c.ShouldBindJSON(&txRawMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "参数错误",
			"detail": err.Error(),
		})
		return
	}

	txRawMsgData, err := json.Marshal(txRawMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "JSON序列化失败",
			"detail": err.Error(),
		})
		return
	}

	resp, err := fabric.ApproveComment(txRawMsgData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "审核评论失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// 查询单个评论
func QueryComment(c *gin.Context) {
	commentID := c.Query("comment_id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论ID不能为空"})
		return
	}

	resp, err := fabric.QueryComment(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}
// 查询指定文档的评论列表
func QueryAllCommentsByRefID(c *gin.Context) {
	var params struct {
		RefID string `json:"ref_id"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "参数错误",
			"detail": err.Error(),
		})
		return
	}

	resp, err := fabric.QueryAllCommentsByRefID(params.RefID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "查询评论列表失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}
// 分页查询评论列表
func QueryCommentsByRefID(c *gin.Context) {
	var params struct {
		RefID    string `json:"ref_id"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryCommentsByRefIDWithPagination(params.RefID, params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

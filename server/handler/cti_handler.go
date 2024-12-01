package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/righstar2020/br-cti-bc-server/fabric"
)

// CTI注册接口(Post)
func RegisterCtiInfo(c *gin.Context) {
	var txMsg fabric.TxMsgData
	//验证请求参数
	if err := c.ShouldBindJSON(&txMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//重新序列化
	txMsgData, err := json.Marshal(txMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("txMsgData:", string(txMsgData))
	// 调用fabric注册CTI信息
	resp, err := fabric.RegisterCtiInfo(txMsgData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
		CtiType     int    `json:"cti_type"`
		PageSize    int    `json:"page_size"`
		Bookmark    string `json:"bookmark"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fabric.QueryCtiInfoByTypeWithPagination(params.CtiType, params.PageSize, params.Bookmark)
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

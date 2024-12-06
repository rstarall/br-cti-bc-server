package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	fabric "github.com/righstar2020/br-cti-bc-server/fabric"
	"log"
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

// 购买情报接口(Post)
func PurchaseCTI(c *gin.Context) {
	var txRawMsg *fabric.TxMsgRawData

	if err := c.ShouldBindJSON(&txRawMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
			"detail": err.Error(),
		})
		log.Printf("参数错误: %s", err)
		return
	}
	

	txRawMsgData, err := json.Marshal(txRawMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "JSON序列化失败", 
			"detail": err.Error(),
		})
		return
	}
	log.Printf("序列化后的数据: %s", string(txRawMsgData))
	// 调用fabric购买CTI
	resp, err := fabric.PurchaseCTI(txRawMsgData)
	
	if err != nil {
		log.Printf("Fabric购买失败: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fabric购买失败",
			"detail": err.Error(),
		})
		return
	}
	log.Printf("Fabric购买成功: %s", resp)
	c.JSON(http.StatusOK, gin.H{"result": resp})
}

//购买模型接口
func PurchaseModel(c *gin.Context) {
	var txRawMsg *fabric.TxMsgRawData

	if err := c.ShouldBindJSON(&txRawMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
			"detail": err.Error(),
		})
		log.Printf("参数错误: %s", err)
		return
	}

	txRawMsgData, err := json.Marshal(txRawMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "JSON序列化失败",
			"detail": err.Error(), 
		})
		return
	}
	log.Printf("序列化后的数据: %s", string(txRawMsgData))
	// 调用fabric购买模型
	resp, err := fabric.PurchaseModel(txRawMsgData)
	
	if err != nil {
		log.Printf("Fabric购买失败: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fabric购买失败",
			"detail": err.Error(),
		})
		return
	}
	log.Printf("Fabric购买成功: %s", resp)
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
//查询用户详细信息
func QueryUserDetailInfo(c *gin.Context) {
	// 从请求中获取参数
	var params struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	resp, err := fabric.QueryUserDetailInfo(params.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户详细信息失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}
//查询所有用户列表
func QueryAllUserList(c *gin.Context) {
	resp, err := fabric.QueryAllUserList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询所有用户列表失败:" + err.Error()})
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
// 查询用户拥有的情报(上传+购买的)
func QueryUserOwnCTIInfos(c *gin.Context) {
	// 从请求中获取参数
	var params struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	resp, err := fabric.QueryUserOwnCTIInfos(params.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户拥有的情报失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

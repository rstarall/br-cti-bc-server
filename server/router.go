package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	handler "github.com/righstar2020/br-cti-bc-server/server/handler"
)

// 添加cors中间件函数
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewRouter(fabricSDK *fabsdk.FabricSDK) *gin.Engine {
	r := gin.New()

	// 使用cors中间件
	r.Use(cors())

	blockchainApi := r.Group("/blockchain")
	{
		// 查询区块信息
		blockchainApi.Any("/queryBlock/:blockID", handler.QueryBlockInfo)
		//查询区块链信息
		blockchainApi.Any("/queryChain", handler.QueryChainInfo)
	}
	contractApi := r.Group("/contract")
	{
		// 调用智能合约
		contractApi.POST("/queryContract", handler.QueryContract)
		contractApi.POST("/invokeContract", handler.InvokeContract)
	}
	txApi := r.Group("/tx")
	{
		//获取交易nonce
		txApi.POST("/getTransactionNonce", handler.GetTxNonce)
	}
	userApi := r.Group("/user")
	{
		//用户链上接口
		userApi.POST("/registerUserAccount", handler.RegisterUserAccount)
		userApi.POST("/purchaseCti", handler.PurchaseCTI)
		userApi.POST("/purchaseModel", handler.PurchaseModel)
		userApi.POST("/queryUserInfo", handler.QueryUserInfo)
		userApi.POST("/queryUserDetailInfo", handler.QueryUserDetailInfo)
		userApi.POST("/queryAllUserList", handler.QueryAllUserList)
		userApi.POST("/getUserStatistics", handler.GetUserStatistics)
		userApi.POST("/queryPointTransactions", handler.QueryPointTransactions)
		userApi.POST("/queryUserOwnCTIInfos", handler.QueryUserOwnCTIInfos)
	}
	ctiApi := r.Group("/cti")
	{
		//CTI接口
		ctiApi.POST("/registerCtiInfo", handler.RegisterCtiInfo)
		ctiApi.POST("/queryCtiInfo", handler.QueryCtiInfo)
		ctiApi.POST("/queryCtiInfoByTypeWithPagination", handler.QueryCtiInfoByTypeWithPagination)
		ctiApi.POST("/queryCtiInfoByType", handler.QueryCtiInfoByType)
		ctiApi.POST("/queryAllCtiInfoWithPagination", handler.QueryAllCtiInfoWithPagination)
		ctiApi.POST("/queryCtiInfoByCTIHash", handler.QueryCtiInfoByCTIHash)
		ctiApi.POST("/queryCtiInfoByCreatorUserID", handler.QueryCtiInfoByCreatorUserID)
		ctiApi.POST("/queryCtiInfoByTypeWithParams", handler.QueryCtiInfoByTypeWithParams)
		ctiApi.POST("/queryCtiInfoByIncentiveMechanismWithPagination", handler.QueryCtiInfoByIncentiveMechanismWithPagination)
	}
	modelApi := r.Group("/model")
	{
		//模型接口
		modelApi.POST("/registerModelInfo", handler.RegisterModelInfo)
		modelApi.POST("/queryModelInfo", handler.QueryModelInfo)
		modelApi.POST("/queryModelInfoWithPagination", handler.QueryModelInfoWithPagination)
		modelApi.POST("/queryModelsByTypeWithPagination", handler.QueryModelsByTypeWithPagination)
		modelApi.POST("/queryModelsByRefCTIId", handler.QueryModelsByRefCTIId)
		modelApi.POST("/queryModelInfoByCreatorUserID", handler.QueryModelInfoByCreatorUserID)
		modelApi.POST("/purchaseModel", handler.PurchaseModel)
		modelApi.POST("/queryAllModelInfoWithPagination", handler.QueryAllModelInfoWithPagination)
		modelApi.POST("/queryModelsByIncentiveMechanismWithPagination", handler.QueryModelsByIncentiveMechanismWithPagination)
	}
	commentApi := r.Group("/comment")
	{
		//评论接口
		commentApi.POST("/registerComment", handler.RegisterComment)
		commentApi.POST("/approveComment", handler.ApproveComment)
		commentApi.POST("/queryComment", handler.QueryComment)
		commentApi.POST("/queryAllCommentsByRefID", handler.QueryAllCommentsByRefID)
		commentApi.POST("/queryCommentsByRefID", handler.QueryCommentsByRefID)
	}
	incentiveApi := r.Group("/incentive")
	{
		//文档激励接口
		incentiveApi.POST("/registerDocIncentiveInfo", handler.RegisterDocIncentiveInfo)
		incentiveApi.POST("/queryDocIncentiveInfo", handler.QueryDocIncentiveInfo)
		incentiveApi.POST("/queryDocIncentiveInfoWithPagination", handler.QueryDocIncentiveInfoWithPagination)
	}
	dataStatApi := r.Group("/dataStat")
	{
		//数据分析接口
		dataStatApi.POST("/queryCTISummaryInfo", handler.QueryCTISummaryInfo)
		dataStatApi.POST("/getDataStatistics", handler.GetDataStatistics)
		dataStatApi.POST("/getUpchainTrend", handler.GetUpchainTrend)
		dataStatApi.POST("/getAttackTypeRanking", handler.GetAttackTypeRanking)
		dataStatApi.POST("/getIOCsDistribution", handler.GetIOCsDistribution)
		dataStatApi.POST("/getGlobalIOCsDistribution", handler.GetGlobalIOCsDistribution)
		dataStatApi.POST("/getSystemOverview", handler.GetSystemOverview)
	}
	ipfsApi := r.Group("/ipfs")
	{
		ipfsApi.Any("/getIPFSContent", handler.GetIPFSContent)
		ipfsApi.POST("/getIPFSFileUrl", handler.GetIPFSFileUrl)
		ipfsApi.POST("/processIOCWorldMapStatistics", handler.ProcessIOCWorldMapStatistics)
		ipfsApi.POST("/getIOCWorldMapStatistics", handler.GetIOCWorldMapStatisticsHandler)
		ipfsApi.POST("/downloadIPFSFile", handler.DownloadFileHandler)
		ipfsApi.POST("/downloadIPFSFileToPath", handler.DownloadFileToPathHandler)
	}
	kpApi := r.Group("/kp")
	{
		kpApi.Any("/queryIOCGeoDistribution", handler.QueryIOCGeoDistribution)
		kpApi.Any("/queryIOCTypeDistribution", handler.QueryIOCTypeDistribution)
		kpApi.Any("/queryAttackTypeStatistics", handler.QueryAttackTypeStatistics)
		kpApi.Any("/queryAttackIOCInfo", handler.QueryAttackIOCInfo)
		//-------- 流量场景统计信息-------------------------------
		kpApi.Any("/queryTrafficTypeRatio", handler.QueryTrafficTypeRatio)
		kpApi.Any("/queryTrafficTimeSeries", handler.QueryTrafficTimeSeries)
	}
	return r
}

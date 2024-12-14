package fabric

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/righstar2020/br-cti-bc-server/global"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
	ipfss "github.com/righstar2020/br-cti-bc-server/service"
	"io/ioutil"
	"os"
)

var nodeAddrs = []string{
	"http://127.0.0.1:8080",
	// "https://ipfs.io",
	// "https://dweb.link",
	// "https://gateway.pinata.cloud",
}

// 情报合约
// 注册情报
// 定义结构体

func RegisterCtiInfo(txRawMsgData []byte) (string, error) {
	// 调用链码注册情报
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 解析 txRawMsgData，提取 tx_data
	var registerData struct {
		Nonce          string `json:"nonce"`
		NonceSignature string `json:"nonce_signature"`
		TxData         string `json:"tx_data"` // tx_data 是 Base64 编码的字符串
		TxSignature    string `json:"tx_signature"`
		UserId         string `json:"user_id"`
	}

	if err := json.Unmarshal(txRawMsgData, &registerData); err != nil {
		return "", fmt.Errorf("解析 txRawMsgData 失败: %v", err)
	}

	// 解码 tx_data
	decodedTxData, err := base64.StdEncoding.DecodeString(registerData.TxData)
	if err != nil {
		return "", fmt.Errorf("解码 tx_data 失败: %v", err)
	}

	// 解析解码后的 tx_data 内容
	var txData struct {
		StatisticInfo string `json:"statistic_info"`
	}
	if err := json.Unmarshal(decodedTxData, &txData); err != nil {
		return "", fmt.Errorf("解析 tx_data 内容失败: %v", err)
	}

	// 从 IPFS 获取文件内容
	ipfsHash := txData.StatisticInfo
	nodeAddr := "http://127.0.0.1:8080" // 假设这是你的 IPFS 节点地址，实际情况根据需要调整

	// 使用 GetIPFSContentFromNode 获取文件内容
	_, content, err := ipfs.GetIPFSContentFromNode(ipfsHash, nodeAddr)
	if err != nil {
		return "", fmt.Errorf("从 IPFS 获取内容失败: %v", err)
	}

	// 调用 ProcessIOCWorldMapStatistics 处理获取的 IPFS 内容
	stats, err := ipfss.NewIPFSService(nodeAddrs, "download").ProcessIOCWorldMapStatistics(string(content))
	if err != nil {
		return "", fmt.Errorf("处理 IOC 数据失败: %v", err)
	}

	// 确保 stats 不是空的
	if len(stats) == 0 {
		return "", fmt.Errorf("处理后的统计信息为空")
	}

	// 打印日志以调试
	fmt.Printf("处理后的统计信息: %v\n", stats)

	// 加载已有的统计信息
	existingStats, err := loadExistingStats()
	if err != nil {
		return "", fmt.Errorf("加载已有统计信息失败: %v", err)
	}

	// 累加新的统计信息
	for country, cityMap := range stats {
		if existingCityMap, exists := existingStats[country]; exists {
			// 如果该国家已经有统计信息，进行累加
			for city, count := range cityMap {
				existingCityMap[city] += count
			}
		} else {
			// 如果该国家没有统计信息，直接添加
			existingStats[country] = cityMap
		}
	}

	// 将累加后的统计信息保存到文件
	savePath, err := ipfss.NewIPFSService(nodeAddrs, "download").SaveIOCWorldMapStatistics(existingStats)
	if err != nil {
		return "", fmt.Errorf("保存统计信息失败: %v", err)
	}
	fmt.Printf("统计信息已保存至: %s\n", savePath)

	// 调用链码，将原始 txRawMsgData 注册到链码中
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "RegisterCTIInfo", [][]byte{txRawMsgData})
	if err != nil {
		return "", fmt.Errorf("调用链码失败: %v", err)
	}

	return string(resp), nil
}

// loadExistingStats 从文件中加载已有的统计信息
func loadExistingStats() (global.StatsMap, error) {
	// 读取已有的统计信息文件
	filePath := "download/ioc_world_map_statistics.json" // 文件路径
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 如果文件不存在，返回空的统计数据
		if os.IsNotExist(err) {
			return global.StatsMap{}, nil
		}
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析现有的统计信息
	var stats global.StatsMap
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("解析已有统计信息失败: %v", err)
	}

	return stats, nil
}

// 购买情报
func PurchaseCTI(txRawMsgData []byte) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码购买情报
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "PurchaseCTI", [][]byte{txRawMsgData})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 查询情报
func QueryCtiInfoByID(ctiID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询情报
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryCTIInfo", [][]byte{[]byte(ctiID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据类型分页查询情报
func QueryCtiInfoByTypeWithPagination(ctiType int, page int, pageSize int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", ctiType)),
		[]byte(fmt.Sprintf("%d", page)),
		[]byte(fmt.Sprintf("%d", pageSize)),
	}

	// 调用链码查询情报
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryCTIInfoByTypeWithPagination", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据类型查询情报
func QueryCtiInfoByType(ctiType int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", ctiType)),
	}

	// 调用链码查询情报
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryCTIInfoByType", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 分页查询所有情报信息
func QueryAllCtiInfoWithPagination(page int, pageSize int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", page)),
		[]byte(fmt.Sprintf("%d", pageSize)),
	}

	// 调用链码查询所有情报
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryAllCTIInfoWithPagination", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据情报哈希查询情报
func QueryCtiInfoByCTIHash(ctiHash string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询情报
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryCTIInfoByCTIHash", [][]byte{[]byte(ctiHash)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据创建者ID查询情报
func QueryCtiInfoByCreatorUserID(userID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询情报
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryCTIInfoByCreatorUserID", [][]byte{[]byte(userID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

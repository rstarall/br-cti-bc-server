package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/righstar2020/br-cti-bc-server/global"
	ipfsService "github.com/righstar2020/br-cti-bc-server/service"

)


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

	// 调用链码，将原始 txRawMsgData 注册到链码中
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "RegisterCTIInfo", [][]byte{txRawMsgData})
	if err != nil {
		return "", fmt.Errorf("调用链码失败: %v", err)
	}

	// 处理统计信息(异步)
	go func() {
		ipfsService.GetIPFSServiceInstance().ProcessCtiRegisterDataStatistics(registerData.TxData)
	}()

	return string(resp), nil
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

package fabric

import (
	"fmt"

	"github.com/righstar2020/br-cti-bc-server/global"
)

// 模型合约

// 注册模型信息
func RegisterModelInfo(modelTxData []byte) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}


	// 调用链码注册模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "RegisterModelInfo", [][]byte{[]byte(modelTxData)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
// 购买模型
func PurchaseModel(txRawMsgData []byte) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码购买模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "PurchaseModel", [][]byte{txRawMsgData})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 查询模型信息
func QueryModelInfo(modelID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryModelInfo", [][]byte{[]byte(modelID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据ID和分页信息查询模型
func QueryModelInfoWithPagination(page int, pageSize int) (string, error) {
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

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryAllModelInfoWithPagination", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据流量类型查询模型
func QueryModelsByTrafficType(trafficType string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryModelsByTrafficType", [][]byte{[]byte(trafficType)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据CTI ID查询模型
func QueryModelsByRefCTIId(ctiID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryModelsByRefCTIId", [][]byte{[]byte(ctiID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据创建者ID查询模型
func QueryModelInfoByCreatorUserID(userID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "queryModelInfoByCreatorUserID", [][]byte{[]byte(userID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}


// 根据用户ID查询模型
func QueryModelsByUserID(userID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryModelsByUserID", [][]byte{[]byte(userID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据模型哈希查询模型
func QueryModelInfoByModelHash(modelHash string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "queryModelInfoByModelHash", [][]byte{[]byte(modelHash)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func QueryModelsByTypeWithPagination(modelType int, page int, pageSize int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", modelType)),
		[]byte(fmt.Sprintf("%d", page)),
		[]byte(fmt.Sprintf("%d", pageSize)),
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryModelsByTypeWithPagination", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func QueryAllModelInfoWithPagination(page int, pageSize int) (string, error) {
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
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryAllModelInfoWithPagination", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
// 根据激励机制分页查询模型信息
func QueryModelsByIncentiveMechanismWithPagination(page int, pageSize int, incentiveMechanism int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", page)),
		[]byte(fmt.Sprintf("%d", pageSize)),
		[]byte(fmt.Sprintf("%d", incentiveMechanism)),
	}

	// 调用链码查询模型
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryModelsByIncentiveMechanismWithPagination", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}


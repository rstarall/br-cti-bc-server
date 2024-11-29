package fabric

import (
	"fmt"

	"github.com/righstar2020/br-cti-bc-server/global"
)

// 情报合约
// 注册情报
func RegisterCtiInfo(txMsgData []byte) (string, error) {
	// 调用链码注册情报
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码注册情报信息
	resp, err := InvokeChaincode(client, "cti_chaincode", "registerCtiInfo", [][]byte{txMsgData})
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
	resp, err := InvokeChaincode(client, "cti_chaincode", "queryCTIInfo", [][]byte{[]byte(ctiID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 根据类型分页查询情报
func QueryCtiInfoByTypeWithPagination(ctiType int, pageSize int, bookmark string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", ctiType)),
		[]byte(fmt.Sprintf("%d", pageSize)), 
		[]byte(bookmark),
	}

	// 调用链码查询情报
	resp, err := InvokeChaincode(client, "cti_chaincode", "queryCTIInfoByTypeWithPagination", args)
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
	resp, err := InvokeChaincode(client, "cti_chaincode", "queryCTIInfoByType", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

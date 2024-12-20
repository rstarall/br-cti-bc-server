package fabric

import (
	"fmt"
	"github.com/righstar2020/br-cti-bc-server/global"
)

// 注册文档激励信息
func RegisterDocIncentiveInfo(txRawMsgData []byte) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码注册文档激励信息
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "RegisterDocIncentiveInfo", [][]byte{txRawMsgData})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 查询文档激励信息
func QueryDocIncentiveInfo(refID string, doctype string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询文档激励信息
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryDocIncentiveInfo", [][]byte{[]byte(refID), []byte(doctype)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
//分页查询
func QueryDocIncentiveInfoWithPagination(refID string, doctype string, page int, pageSize int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(refID),
		[]byte(doctype),
		[]byte(fmt.Sprintf("%d", page)),
		[]byte(fmt.Sprintf("%d", pageSize)),
	}

	// 调用链码分页查询文档激励信息
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryDocIncentiveInfoByPage", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}


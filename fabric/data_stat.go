package fabric

import (
	"github.com/righstar2020/br-cti-bc-server/global"
)

// 数据分析合约
// 查询情报统计信息
func QueryCTISummaryInfo(ctiID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询情报统计信息
	resp, err := InvokeChaincode(client, "data_chaincode", "queryCTISummaryInfoByCTIID", [][]byte{[]byte(ctiID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取数据统计信息
func GetDataStatistics() (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码获取数据统计信息
	resp, err := InvokeChaincode(client, "data_chaincode", "getDataStatistics", [][]byte{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取情报交易趋势数据
func GetCTITrafficTrend(timeRange string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码获取趋势数据
	resp, err := InvokeChaincode(client, "data_chaincode", "getCTITrafficTrend", [][]byte{[]byte(timeRange)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取攻击类型排行
func GetAttackTypeRanking() (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码获取排行数据
	resp, err := InvokeChaincode(client, "data_chaincode", "getAttackTypeRanking", [][]byte{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取IOCs类型分布
func GetIOCsDistribution() (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码获取分布数据
	resp, err := InvokeChaincode(client, "data_chaincode", "getIOCsDistribution", [][]byte{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取全球IOCs地理分布
func GetGlobalIOCsDistribution() (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码获取地理分布数据
	resp, err := InvokeChaincode(client, "data_chaincode", "getGlobalIOCsDistribution", [][]byte{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取系统概览数据
func GetSystemOverview() (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码获取系统概览数据
	resp, err := InvokeChaincode(client, "data_chaincode", "getSystemOverview", [][]byte{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

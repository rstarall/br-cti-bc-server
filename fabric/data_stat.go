package fabric

import (
	"github.com/righstar2020/br-cti-bc-server/global"
	"fmt"
	"encoding/json"
)

// 数据分析合约
// 查询最新的情报统计信息
func QueryLatestCTISummaryInfo(limit int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询最新的情报统计信息
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryLatestCTISummaryInfo", [][]byte{[]byte(fmt.Sprintf("%d", limit))})
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
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "GetDataStatistics", [][]byte{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取上链趋势数据
func GetUpchainTrend(timeRange string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码获取趋势数据
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "GetUpchainTrend", [][]byte{[]byte(timeRange)})
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
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "GetAttackTypeRanking", [][]byte{})
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
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "GetIOCsDistribution", [][]byte{})
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
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "GetGlobalIOCsDistribution", [][]byte{})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 获取系统概览数据
func GetSystemOverview() (*SystemOverviewInfo, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return nil, err
	}

	// 调用链码获取系统概览数据
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "GetSystemOverview", [][]byte{})
	if err != nil {
		return nil, err
	}
	var SystemOverviewInfoData SystemOverviewInfo
	err = json.Unmarshal(resp, &SystemOverviewInfoData)
	if err != nil {
		return nil, err
	}
    //更新区块高度信息
	blockHeight, err := GetBlockHeight(global.LedgerClient)
	if err != nil {
		return nil, err
	}

	SystemOverviewInfoData.BlockHeight = blockHeight
	return &SystemOverviewInfoData, nil
}

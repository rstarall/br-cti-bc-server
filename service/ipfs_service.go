package ipfsservice

import (
	"encoding/json"
	"fmt"
	"sync"
)

// StatsMap 定义统计数据的类型
type StatsMap map[string]map[string]int

// IPFSService 提供处理 IPFS 内容的服务
type IPFSService struct {
	mu        sync.RWMutex
	statsMap  StatsMap
	ipfsNodes []string
}

// NewIPFSService 创建一个新的 IPFSService实例
func NewIPFSService(ipfsNodes []string) *IPFSService {
	return &IPFSService{
		statsMap:  make(StatsMap),
		ipfsNodes: ipfsNodes,
	}
}

//----------------------------------IOC地理位置信息------------------------------------

// ProcessIOCWorldMapStatistics 处理 IOC 数据，生成统计信息
func (s *IPFSService) ProcessIOCWorldMapStatistics(iocData string) (StatsMap, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(iocData), &data); err != nil {
		return nil, fmt.Errorf("解析 IOC 数据失败: %v", err)
	}

	localStats := make(StatsMap)

	// 示例处理逻辑，根据您的具体数据结构调整
	if ipsInfo, exists := data["ips_locations_info_map"].(map[string]interface{}); exists {
		for _, locInfo := range ipsInfo {
			locMap, ok := locInfo.(map[string]interface{})
			if !ok {
				continue
			}

			country, ok1 := locMap["country"].(string)
			region, ok2 := locMap["region"].(string)
			if !ok1 || !ok2 {
				continue
			}

			if _, exists := localStats[country]; !exists {
				localStats[country] = make(map[string]int)
			}
			localStats[country][region]++
		}
	}

	// 更新全局统计数据（线程安全）
	s.mu.Lock()
	defer s.mu.Unlock()
	for country, regions := range localStats {
		if _, exists := s.statsMap[country]; !exists {
			s.statsMap[country] = make(map[string]int)
		}
		for region, count := range regions {
			s.statsMap[country][region] += count
		}
	}

	return localStats, nil
}

//----------------------------------IOC统计信息------------------------------------
//func (s *IPFSService) ProcessIOCStatistics(iocData string) (string, error) {
//	// 处理IOC统计信息
//}
//func (s *IPFSService) SaveIOCStatistics(iocData string) (string, error) {
//	// 保存IOC统计信息(保存在数据库或文件中)
//}
//func (s *IPFSService) GetIOCStatistics() (string, error) {
//	// 获取IOC统计信息
//}

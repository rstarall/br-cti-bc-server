package service

import (
	"encoding/json"
	"fmt"
	"github.com/righstar2020/br-cti-bc-server/global"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// StatsMap 定义统计数据的类型

// IPFSService 提供处理 IPFS 内容的服务
type IPFSService struct {
	mu          sync.RWMutex
	statsMap    global.StatsMap
	ipfsNodes   []string
	downloadDir string // 添加 downloadDir 字段
}

// NewIPFSService 创建一个新的 IPFSService实例
func NewIPFSService(ipfsNodes []string, downloadDir string) *IPFSService {
	return &IPFSService{
		statsMap:    make(global.StatsMap),
		ipfsNodes:   ipfsNodes,
		downloadDir: downloadDir, // 初始化 downloadDir
	}
}

func (s *IPFSService) DownloadFile(hash string) (string, error) {
	// 获取 IPFS 内容
	_, content, err := ipfs.GetIPFSContentWithFallback(hash)
	if err != nil {
		return "", err
	}

	// 创建保存文件的目录（例如 ./downloads）
	if err := os.MkdirAll(s.downloadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建下载目录失败: %v", err)
	}

	fileName := hash + ".json"

	// 定义文件路径
	filePath := filepath.Join(s.downloadDir, fileName)

	// 写入文件
	if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	return filePath, nil
}

//----------------------------------IOC地理位置信息------------------------------------

// ProcessIOCWorldMapStatistics 处理 IOC 数据，生成统计信息
func (s *IPFSService) ProcessIOCWorldMapStatistics(iocData string) (global.StatsMap, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(iocData), &data); err != nil {
		return nil, fmt.Errorf("解析 IOC 数据失败: %v", err)
	}

	localStats := make(global.StatsMap)

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

func (s *IPFSService) SaveIOCWorldMapStatistics(iocData global.StatsMap) (string, error) {
	// 数据检查
	if len(iocData) == 0 {
		return "", fmt.Errorf("统计数据为空，无法保存")
	}

	// 将 StatsMap 转换为 JSON 格式
	statsJSON, err := json.MarshalIndent(iocData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("生成统计 JSON 失败: %v", err)
	}

	// 打印 JSON 格式内容，确保正确
	fmt.Printf("生成的统计信息：\n%s\n", statsJSON)

	// 定义保存路径
	savePath := filepath.Join(s.downloadDir, "ioc_world_map_statistics.json")

	// 写入文件
	if err := ioutil.WriteFile(savePath, statsJSON, 0644); err != nil {
		return "", fmt.Errorf("保存统计数据失败: %v", err)
	}

	return savePath, nil
}

func (s *IPFSService) GetIOCWorldMapStatistics() (string, error) {
	// 定义统计文件的保存路径
	savePath := filepath.Join(s.downloadDir, "ioc_world_map_statistics.json")

	// 读取文件内容
	content, err := ioutil.ReadFile(savePath)
	if err != nil {
		return "", fmt.Errorf("读取统计数据文件失败: %v", err)
	}

	// 返回文件内容（作为 JSON 格式字符串）
	return string(content), nil
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

package service

import (
	"encoding/json"
	"fmt"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"encoding/base64"
)



// IOCWorldMapStats 定义统计数据的类型(国家-城市-数量)
type IOCWorldStatsMap map[string]map[string]int

// IPFSService 提供处理 IPFS 内容的服务
type IPFSService struct {
	mu          sync.RWMutex
	iocWorldStatsMap    IOCWorldStatsMap
	iocWorldStatsFilePath string // 添加 statsFilePath 字段
	ipfsNodes   []string
	downloadDir string // 添加 downloadDir 字段
}

var nodeAddrs = []string{
	"http://127.0.0.1:8080",
}

var IPFSServiceInstance *IPFSService

// GetIPFSServiceInstance 获取 IPFSService 实例
func GetIPFSServiceInstance() *IPFSService {
	if IPFSServiceInstance == nil {
		IPFSServiceInstance = NewIPFSService(nodeAddrs, "ipfs_download")
	}
	return 	IPFSServiceInstance	
}

// NewIPFSService 创建一个新的 IPFSService实例
func NewIPFSService(ipfsNodes []string, downloadDir string) *IPFSService {
	return &IPFSService{
		iocWorldStatsMap:    make(IOCWorldStatsMap),
		iocWorldStatsFilePath: filepath.Join(downloadDir, "ioc_world_map_statistics.json"),
		ipfsNodes:   ipfsNodes,
		downloadDir: downloadDir, // 初始化 downloadDir
	}
}

func (s *IPFSService) DownloadIPFSFile(hash string) (string,string, error) {
	// 获取 IPFS 内容
	_, rawContent, err := ipfs.GetIPFSContentWithFallback(hash)
	if rawContent == nil {
		return "", "", fmt.Errorf("IPFS内容为空")
	}

	content := string(rawContent)

	if err != nil {
		return "", "", err
	}
	
	// 创建保存文件的目录（例如 ./downloads）
	if err := os.MkdirAll(s.downloadDir, os.ModePerm); err != nil {
		return "", content, fmt.Errorf("创建下载目录失败: %v", err)
	}

	fileName := hash + ".json"

	// 定义文件路径
	filePath := filepath.Join(s.downloadDir, fileName)

	// 写入文件
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", content, fmt.Errorf("保存文件失败: %v", err)
	}

	return filePath, content, nil
}
//----------------------------------处理IPFS信息---------------------------------------
func (s *IPFSService) ProcessCtiRegisterDataStatistics(TxDataString string) (string, error) {
	
	//base64解码
	decodedTxData, err := base64.StdEncoding.DecodeString(TxDataString)
	if err != nil {
		return "", fmt.Errorf("解码CTI注册信息失败: %v", err)
	}

	// 处理CTI注册信息
	var txData struct {
		StatisticInfo string `json:"statistic_info"` //只需要统计信息
	}

	if err := json.Unmarshal(decodedTxData, &txData); err != nil {
		return "", fmt.Errorf("解析CTI注册信息失败: %v", err)
	}

	// 处理统计信息
	// 从 IPFS 获取文件内容
	ipfsHash := txData.StatisticInfo

	// 使用 GetIPFSContentFromNode 获取文件内容
	_, content, err := s.DownloadIPFSFile(ipfsHash)
	if err != nil {
		return "", fmt.Errorf("从 IPFS 获取内容失败: %v", err)
	}

	s.ProcessIOCWorldMapStatistics(string(content))

	// 确保 stats 不是空的
	if len(s.iocWorldStatsMap) == 0 {
		return "", fmt.Errorf("处理后的统计信息为空")
	}
	//更新统计信息(异步)
	go func() {
		s.UpdateIOCWorldMapStatistics()
	}()

	return txData.StatisticInfo, nil
}

func (s *IPFSService) UpdateIOCWorldMapStatistics() (string, error) {
	//线程安全
	s.mu.Lock()
	defer s.mu.Unlock()
	// 打印日志以调试
	fmt.Printf("处理后的统计信息: %v\n", s.iocWorldStatsMap)
	existingStats := s.iocWorldStatsMap
	if len(s.iocWorldStatsMap) == 0 {
		fmt.Println("处理后的统计信息为空")
		// 加载已有的统计信息
		loadExistingStats, err := s.loadExistingStats(s.iocWorldStatsFilePath)
		if err != nil {
			return "", fmt.Errorf("加载已有统计信息失败: %v", err)
		}
		existingStats = loadExistingStats
	}
	

	// 累加新的统计信息
	for country, cityMap := range s.iocWorldStatsMap {
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
	savePath, err := s.SaveIOCWorldMapStatistics(existingStats)
	if err != nil {
		return "", fmt.Errorf("保存统计信息失败: %v", err)
	}
	fmt.Printf("统计信息已保存至: %s\n", savePath)

	return savePath, nil
}

// loadExistingStats 从文件中加载已有的统计信息
func (s *IPFSService) loadExistingStats(filePath string) (IOCWorldStatsMap, error) {
	// 读取已有的统计信息文件
	
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 如果文件不存在，返回空的统计数据
		if os.IsNotExist(err) {
			return IOCWorldStatsMap{}, nil
		}
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析现有的统计信息
	var stats IOCWorldStatsMap
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("解析已有统计信息失败: %v", err)
	}

	return stats, nil
}

//----------------------------------IOC地理位置信息------------------------------------

// ProcessIOCWorldMapStatistics 处理 IOC 数据，生成统计信息
func (s *IPFSService) ProcessIOCWorldMapStatistics(iocData string) (IOCWorldStatsMap, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(iocData), &data); err != nil {
		return nil, fmt.Errorf("解析 IOC 数据失败: %v", err)
	}

	localStats := make(IOCWorldStatsMap)

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

	//线程安全
	s.mu.Lock()
	defer s.mu.Unlock()
	for country, regions := range localStats {
		if _, exists := s.iocWorldStatsMap[country]; !exists {
			s.iocWorldStatsMap[country] = make(map[string]int)
		}
		for region, count := range regions {
			s.iocWorldStatsMap[country][region] += count
		}
	}

	return localStats, nil
}

func (s *IPFSService) SaveIOCWorldMapStatistics(iocData IOCWorldStatsMap) (string, error) {
	
	// 数据检查
	if len(iocData) == 0 {
		return "", fmt.Errorf("统计数据为空，无法保存")
	}

	// 将 IOCWorldStatsMap 转换为 JSON 格式
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
	var iocWorldStatsMap IOCWorldStatsMap
	//线程安全
	s.mu.Lock()
	defer s.mu.Unlock()
     //先从内存中读取
	if len(s.iocWorldStatsMap) == 0 {
		 fmt.Println("内存中统计信息为空")
		 //否则从文件读取内容
		content, err := ioutil.ReadFile(s.iocWorldStatsFilePath)
		if err != nil {
			return "", fmt.Errorf("读取统计数据文件失败: %v", err)
		}
		//将文件内容转换为IOCWorldStatsMap
		if err := json.Unmarshal(content, &iocWorldStatsMap); err != nil {
			return "", fmt.Errorf("解析统计数据文件失败: %v", err)
		}
		
		s.iocWorldStatsMap = iocWorldStatsMap
	}
	
	statsJSON, err := json.MarshalIndent(s.iocWorldStatsMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("生成统计 JSON 失败: %v", err)
	}
    
	// 返回文件内容（作为 JSON 格式字符串）
	return string(statsJSON), nil
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

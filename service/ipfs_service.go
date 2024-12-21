package service

import (
	"encoding/json"
	"fmt"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
	"github.com/righstar2020/br-cti-bc-server/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"encoding/base64"
	"strings"
	"time"
	
)



// IOCWorldMapStats 定义统计数据的类型(国家-城市-数量)
type IOCWorldStatsMap map[string]map[string]int

// IPFSService 提供处理 IPFS 内容的服务
type IPFSService struct {
	mu          sync.RWMutex
	iocWorldStatsMap    IOCWorldStatsMap
	iocWorldStatsFilePath string // 添加 statsFilePath 字段
	ipfsNodes   []string
	statsFilePath string // 添加 statsFilePath 字段
	downloadDir string // 添加 downloadDir 字段
}

var nodeAddrs = []string{
	"http://127.0.0.1:8080",
}

var IPFSServiceInstance *IPFSService

// GetIPFSServiceInstance 获取 IPFSService 实例
func GetIPFSServiceInstance() *IPFSService {
	if IPFSServiceInstance == nil {
		projectAbsPath, err := util.GetProjectAbsPath()
		if err != nil {
			fmt.Printf("获取项目绝对路径失败: %v", err)
		}
		ipfsDownloadDir := filepath.Join(projectAbsPath, "ipfs_download")
		statsDbDir := filepath.Join(projectAbsPath, "stats_db")
		if !util.PathExists(ipfsDownloadDir) {
			//创建ipfs_download目录
			if err := os.MkdirAll(ipfsDownloadDir, os.ModePerm); err != nil {
				fmt.Printf("创建ipfs_download目录失败: %v", err)
			}
		}
		if !util.PathExists(statsDbDir) {
			//创建stats_db目录
			if err := os.MkdirAll(statsDbDir, os.ModePerm); err != nil {
				fmt.Printf("创建stats_db目录失败: %v", err)
			}
		}
		
		IPFSServiceInstance = NewIPFSService(nodeAddrs, ipfsDownloadDir,statsDbDir)
	}
	return 	IPFSServiceInstance	
}

// NewIPFSService 创建一个新的 IPFSService实例
func NewIPFSService(ipfsNodes []string, downloadDir string,statsFilePath string) *IPFSService {
	return &IPFSService{
		iocWorldStatsMap:    make(IOCWorldStatsMap),
		iocWorldStatsFilePath: filepath.Join(statsFilePath, "ioc_world_map_statistics.json"),
		ipfsNodes:   ipfsNodes,
		downloadDir: downloadDir, // 初始化 downloadDir
		statsFilePath: statsFilePath, // 初始化 statsFilePath
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
	var ctiTxData CtiTxData

	if err := json.Unmarshal(decodedTxData, &ctiTxData); err != nil {
		return "", fmt.Errorf("解析CTI注册信息失败: %v", err)
	}

	// 处理统计信息
	// 从 IPFS 获取文件内容
	ipfsHash := ctiTxData.StatisticInfo

	// 使用 GetIPFSContentFromNode 获取文件内容
	_, content, err := s.DownloadIPFSFile(ipfsHash)
	if err != nil {
		return "", fmt.Errorf("从 IPFS 获取内容失败: %v", err)
	}

	
	//1.更新地理位置统计信息(异步)
	go func() {
		s.ProcessIOCWorldMapStatistics(string(content))

		// 确保 stats 不是空的
		if len(s.iocWorldStatsMap) == 0 {
			fmt.Printf("处理后的统计信息为空")
		}
		s.UpdateIOCWorldMapStatistics()
	}()
	//2.更新IOC类型统计信息(异步)
	go func() {
		iocTypeStats, err := s.ProcessIOCTypeStatistics(&ctiTxData)
		if err != nil {
			fmt.Printf("处理IOC类型统计信息失败: %v", err)
		}
		s.SaveIOCTypeStatistics(iocTypeStats)
	}()
	//3.更新攻击类型统计信息(异步)
	go func() {
		attackTypeStats, err := s.ProcessAttackTypeStatistics(&ctiTxData)
		if err != nil {
			fmt.Printf("处理攻击类型统计信息失败: %v", err)
		}
		s.SaveAttackTypeStatistics(attackTypeStats)
	}()
	//4.更新攻击IOC信息(异步)
	go func() {
		attackIOCInfo, err := s.ProcessAttackIOCInfo(string(content),&ctiTxData)
		if err != nil {
			fmt.Printf("处理攻击IOC信息失败: %v", err)
		}
		s.SaveAttackIOCInfo(attackIOCInfo)
	}()

	return ctiTxData.StatisticInfo, nil
}
//-------------------------1.更新IOC世界地图统计信息--------------------------------
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

//----------------------------------1.IOC地理位置信息------------------------------------

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
	savePath := filepath.Join(s.statsFilePath, "ioc_world_map_statistics.json")

	// 写入文件
	if err := ioutil.WriteFile(savePath, statsJSON, 0644); err != nil {
		return "", fmt.Errorf("保存统计数据失败: %v", err)
	}

	return savePath, nil
}
func (s *IPFSService) GetIOCWorldMapStatistics() (IOCWorldStatsMap, error) {
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
			return nil, fmt.Errorf("读取统计数据文件失败: %v", err)
		}
		//将文件内容转换为IOCWorldStatsMap
		if err := json.Unmarshal(content, &iocWorldStatsMap); err != nil {
			return nil, fmt.Errorf("解析统计数据文件失败: %v", err)
		}
		
		s.iocWorldStatsMap = iocWorldStatsMap
	}
	
	return s.iocWorldStatsMap, nil
}

//-------------------------------------2.IOC类型分布-------------------------------------
// IOCTypeStats 定义IOC类型统计数据结构
type IOCTypeStats struct {
	IPCount     int `json:"ip_count"`
	PortCount   int `json:"port_count"`
	PayloadCount int `json:"payload_count"`
	URLCount    int `json:"url_count"`
	HashCount   int `json:"hash_count"`
}

// ProcessIOCTypeStatistics 处理IOC类型分布统计
func (s *IPFSService) ProcessIOCTypeStatistics(ctiData *CtiTxData) (*IOCTypeStats, error) {
	
	// 获取历史统计数据
	existingStats := &IOCTypeStats{}
	savePath := filepath.Join(s.statsFilePath, "ioc_type_statistics.json")
	content, err := ioutil.ReadFile(savePath)
	if err == nil {
		// 如果文件存在,解析历史数据
		if err := json.Unmarshal(content, existingStats); err != nil {
			return nil, fmt.Errorf("解析历史统计数据失败: %v", err)
		}
	}

	// 统计新的IOC数据
	newStats := &IOCTypeStats{}
	for _, ioc := range ctiData.IOCs {
		if strings.Contains(ioc, "ip") { // IP地址
			newStats.IPCount++
		} else if strings.Contains(ioc, "port") { // 端口号
			newStats.PortCount++
		} else if strings.HasPrefix(ioc, "url") { // URL
			newStats.URLCount++
		} else if strings.HasPrefix(ioc, "hash") { // Hash值(假设是SHA256)
			newStats.HashCount++
		} else if strings.HasPrefix(ioc, "payload") { // Payload
			newStats.PayloadCount++
		}
	}

	// 累加新旧数据
	existingStats.IPCount += newStats.IPCount
	existingStats.PortCount += newStats.PortCount
	existingStats.PayloadCount += newStats.PayloadCount
	existingStats.URLCount += newStats.URLCount
	existingStats.HashCount += newStats.HashCount

	return existingStats, nil
}

// SaveIOCTypeStatistics 保存IOC类型统计信息
func (s *IPFSService) SaveIOCTypeStatistics(stats *IOCTypeStats) error {
	statsJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("生成统计JSON失败: %v", err)
	}

	savePath := filepath.Join(s.statsFilePath, "ioc_type_statistics.json")
	if err := ioutil.WriteFile(savePath, statsJSON, 0644); err != nil {
		return fmt.Errorf("保存统计数据失败: %v", err)
	}
	fmt.Printf("IOC类型统计信息已保存至: %s\n", savePath)
	return nil
}

// GetIOCTypeStatistics 获取IOC类型统计信息
func (s *IPFSService) GetIOCTypeStatistics() (*IOCTypeStats, error) {
	savePath := filepath.Join(s.statsFilePath, "ioc_type_statistics.json")
	content, err := ioutil.ReadFile(savePath)
	if err != nil {
		return nil, fmt.Errorf("读取IOC类型统计数据失败: %v", err)
	}

	var stats IOCTypeStats
	if err := json.Unmarshal(content, &stats); err != nil {
		return nil, fmt.Errorf("解析IOC类型统计数据失败: %v", err)
	}

	return &stats, nil
}

//-------------------------------------3.攻击类型统计-------------------------------------
// AttackTypeStats 定义攻击类型统计数据结构
type AttackTypeStats struct {
	MaliciousTraffic int `json:"malicious_traffic"` // 恶意流量
	HoneypotInfo     int `json:"honeypot_info"`     // 蜜罐情报
	Botnet           int `json:"botnet"`            // 僵尸网络
	AppLayerAttack   int `json:"app_layer_attack"`  // 应用层攻击
	OpenSourceInfo   int `json:"open_source_info"`  // 开源情报
	Total            int `json:"total"`             // 总量
}
//定义时序数据结构
type AttackTypeTimeSeriesData struct {
	Time string `json:"time"` // 小时时间戳,格式为"2006-01-02 15:00:00"
	Stats AttackTypeStats `json:"stats"` // 统计数据
}

// ProcessAttackTypeStatistics 处理攻击类型统计
func (s *IPFSService) ProcessAttackTypeStatistics(ctiData *CtiTxData) ([]*AttackTypeTimeSeriesData, error) {
	// 从文件中读取现有统计数据
	savePath := filepath.Join(s.statsFilePath, "attack_type_statistics.json")
	var timeSeriesData []*AttackTypeTimeSeriesData
	
	content, err := ioutil.ReadFile(savePath)
	if err == nil {
		if err := json.Unmarshal(content, &timeSeriesData); err != nil {
			return nil, fmt.Errorf("解析现有攻击类型统计数据失败: %v", err)
		}
	}

	// 获取当前小时时间戳
	currentHour := time.Now().Format("2006-01-02 15:00:00")
	
	// 查找当前小时的统计数据
	var currentStats *AttackTypeTimeSeriesData
	for _, data := range timeSeriesData {
		if data.Time == currentHour {
			currentStats = data
			break
		}
	}

	// 如果当前小时没有数据,创建新的统计数据
	if currentStats == nil {
		currentStats = &AttackTypeTimeSeriesData{
			Time: currentHour,
			Stats: AttackTypeStats{},
		}
		timeSeriesData = append(timeSeriesData, currentStats)
	}

	// 根据CTI类型增加对应计数
	switch ctiData.CTIType {
	case 1:
		currentStats.Stats.MaliciousTraffic++
	case 2:
		currentStats.Stats.HoneypotInfo++
	case 3:
		currentStats.Stats.Botnet++
	case 4:
		currentStats.Stats.AppLayerAttack++
	case 5:
		currentStats.Stats.OpenSourceInfo++
	}
	currentStats.Stats.Total++

	return timeSeriesData, nil
}

// SaveAttackTypeStatistics 保存攻击类型统计信息
func (s *IPFSService) SaveAttackTypeStatistics(stats []*AttackTypeTimeSeriesData) error {
	statsJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("生成统计JSON失败: %v", err)
	}

	savePath := filepath.Join(s.statsFilePath, "attack_type_statistics.json")
	if err := ioutil.WriteFile(savePath, statsJSON, 0644); err != nil {
		return fmt.Errorf("保存统计数据失败: %v", err)
	}
	fmt.Printf("攻击类型统计信息已保存至: %s\n", savePath)
	return nil
}

// GetAttackTypeStatistics 获取攻击类型统计信息
func (s *IPFSService) GetAttackTypeStatistics() ([]*AttackTypeTimeSeriesData, error) {
	savePath := filepath.Join(s.statsFilePath, "attack_type_statistics.json")
	content, err := ioutil.ReadFile(savePath)
	if err != nil {
		return nil, fmt.Errorf("读取攻击类型统计数据失败: %v", err)
	}

	var stats []*AttackTypeTimeSeriesData
	if err := json.Unmarshal(content, &stats); err != nil {
		return nil, fmt.Errorf("解析攻击类型统计数据失败: %v", err)
	}

	return stats, nil
}

//-------------------------------------4.攻击IOC信息-------------------------------------
// AttackIOCInfo 定义攻击IOC信息结构
type AttackIOCInfo struct {
	IPAddress     string `json:"ip_address"`      // IP地址
	Location     string `json:"location"`   // 地理位置
	AttackType    string `json:"attack_type"`     // 攻击类型
	Port          string    `json:"port"`            // 端口
	Hash          string `json:"hash"`            // HASH
}

//获取攻击类型
func (s *IPFSService) GetAttackType(ctiType int) string {
	mapAttackType := map[int]string{
		1: "恶意流量",
		2: "蜜罐情报",
		3: "僵尸网络",
		4: "应用层攻击",
		5: "开源情报",
	}
	return mapAttackType[ctiType]
}

// ProcessAttackIOCInfo 处理攻击IOC信息
func (s *IPFSService) ProcessAttackIOCInfo(iocData string,ctiTxData *CtiTxData) ([]*AttackIOCInfo, error) {
	// 先从文件中读取现有数据
	savePath := filepath.Join(s.statsFilePath, "attack_ioc_info.json")
	var existingInfos []*AttackIOCInfo
	
	content, err := ioutil.ReadFile(savePath)
	if err == nil {
		if err := json.Unmarshal(content, &existingInfos); err != nil {
			return nil, fmt.Errorf("解析现有IOC信息失败: %v", err)
		}
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(iocData), &data); err != nil {
		return nil, fmt.Errorf("解析IOC数据失败: %v", err)
	}
	var iocInfos []*AttackIOCInfo
	// 从ips_locations_info_map中提取信息
	if ipsInfo, ok := data["ips_locations_info_map"].(map[string]interface{}); ok {
		for ip, info := range ipsInfo {
			if infoMap, ok := info.(map[string]interface{}); ok {
				iocInfo := AttackIOCInfo{
					IPAddress: ip,
				}

				// 构建地理位置字符串
				if country, ok := infoMap["country"].(string); ok {
					if region, ok := infoMap["region"].(string); ok {
						if city, ok := infoMap["city"].(string); ok {
							iocInfo.Location = country + "," + region + "," + city
						}
					}
				}

				//根据ctiTxData的CTIType设置攻击类型
				if attackType, ok := infoMap["attack_type"].(string); ok {
					iocInfo.AttackType = attackType
				} else { 
					iocInfo.AttackType = s.GetAttackType(int(ctiTxData.CTIType))
				}

				if port, ok := infoMap["port"].(string); ok {
					iocInfo.Port = port
				}

				if hash, ok := infoMap["hash"].(string); ok {
					iocInfo.Hash = hash
				}

				iocInfos = append(iocInfos, &iocInfo)
			}
		}
	}
	//将新数据追加到现有数据中
	newInfos := append(iocInfos, existingInfos...)
	return newInfos, nil
}

// SaveAttackIOCInfo 保存攻击IOC信息
func (s *IPFSService) SaveAttackIOCInfo(infos []*AttackIOCInfo) error {
	infoJSON, err := json.MarshalIndent(infos, "", "  ")
	if err != nil {
		return fmt.Errorf("生成IOC信息JSON失败: %v", err)
	}

	savePath := filepath.Join(s.statsFilePath, "attack_ioc_info.json")
	if err := ioutil.WriteFile(savePath, infoJSON, 0644); err != nil {
		return fmt.Errorf("保存IOC信息失败: %v", err)
	}
	fmt.Printf("攻击IOC信息已保存至: %s\n", savePath)
	return nil
}

// GetAttackIOCInfo 获取攻击IOC信息
func (s *IPFSService) GetAttackIOCInfo() ([]AttackIOCInfo, error) {
	savePath := filepath.Join(s.statsFilePath, "attack_ioc_info.json")
	content, err := ioutil.ReadFile(savePath)
	if err != nil {
		return nil, fmt.Errorf("读取攻击IOC信息失败: %v", err)
	}

	var infos []AttackIOCInfo
	if err := json.Unmarshal(content, &infos); err != nil {
		return nil, fmt.Errorf("解析攻击IOC信息失败: %v", err)
	}

	return infos, nil
}




type IPFSService struct {
}

func (s *IPFSService) DownloadFile(hash string) (string, error) {

}

func (s *IPFSService) GetIOCFileFromIPFS(IPFSHash string) (string, error) {

}
//----------------------------------IOC地理位置信息------------------------------------
func (s *IPFSService) ProcessIOCWorldMapStatistics(iocData string) (string, error) {
	// 处理IOC世界地图统计数据
	//统计精确到country和region的IOC数量
}
func (s *IPFSService) SaveIOCWorldMapStatistics(iocData string) (string, error) {
	// 保存IOC世界地图统计数据(保存在数据库或文件中)
	//记录每个地区的IOC总数量
}
func (s *IPFSService) GetIOCWorldMapStatistics() (string, error) {
	// 获取IOC世界地图统计数据
}



//----------------------------------IOC统计信息------------------------------------
func (s *IPFSService) ProcessIOCStatistics(iocData string) (string, error) {
	// 处理IOC统计信息
}
func (s *IPFSService) SaveIOCStatistics(iocData string) (string, error) {
	// 保存IOC统计信息(保存在数据库或文件中)
}
func (s *IPFSService) GetIOCStatistics() (string, error) {
	// 获取IOC统计信息
}

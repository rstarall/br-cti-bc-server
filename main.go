package main

import (
	"fmt"
	"time"

	fabric "github.com/righstar2020/br-cti-bc-server/fabric"
	global "github.com/righstar2020/br-cti-bc-server/global"
	"github.com/righstar2020/br-cti-bc-server/ipfs"
	server "github.com/righstar2020/br-cti-bc-server/server"
	"github.com/righstar2020/br-cti-bc-server/service"
)

const fabric_config_path = "./config/server/config.yaml"

func main() {
	// 启动Fabric SDK
	var err error
	//该程序只能在当前工作目录下进行(子进程调用需要切换工作目录pwd)
	global.FabricSDK, err = fabric.NewSDK(fabric_config_path) // 确保路径正确
	if err != nil {
		fmt.Printf("Failed to initialize SDK: %v", err)
	}
	defer global.FabricSDK.Close()

	// 初始化fabric client
	fmt.Println("Fabric client initializing ...")
	InitFabricClient()

	// 启动定时任务
	go StartScheduledTasks()

	// 启动服务器
	r := server.NewRouter(global.FabricSDK)
	if err := r.Run(":7777"); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}

// 初始化Fabric客户端
func InitFabricClient() {
	var err error
	global.ChannelClient, err = fabric.CreateChannelClient(global.FabricSDK)
	if err != nil {
		fmt.Printf("Failed to connect fabric chain: %v", err)
	}
	global.LedgerClient, err = fabric.CreateLedgerClient(global.FabricSDK)
	if err != nil {
		fmt.Printf("Failed to connect fabric chain: %v", err)
	}
}

// 启动定时任务
func StartScheduledTasks() {
	// 创建一个定时器，每小时执行一次任务
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// 每小时执行一次
	for {
		select {
		case <-ticker.C:
			// 获取统计信息
			stats, err := service.NewIPFSService(ipfs.NodeAddrs, "download").GetIOCWorldMapStatistics()
			if err != nil {
				fmt.Printf("Error getting IOC World Map Statistics: %v\n", err)
				continue
			}
			// 打印统计信息（或执行其他操作）
			fmt.Println("Received IOC World Map Statistics:", stats)
		}
	}
}

# bc-server 端开发说明

## 1.目录结构

```
.
├── config/                                 # 配置文件目录
│   ├── local/                             # 本地开发配置
│   │   └── config.yaml                    # 本地配置文件
│   └── server/                            # 服务器配置
│       ├── config.yaml                    # 当前使用的配置文件
│       └── organizations/                 # 组织配置目录
├── docs/                                  # 文档目录
├── fabric/                               # Fabric链码接口层
│   ├── comment.go                        # 评论相关链码接口
│   ├── cti.go                           # 威胁情报链码接口
│   ├── data_stat.go                     # 数据统计链码接口
│   ├── fabric.go                        # Fabric核心功能接口
│   ├── incentive.go                     # 激励机制链码接口
│   ├── model.go                         # 模型相关链码接口
│   ├── msg.go                           # 消息处理链码接口
│   ├── tx.go                            # 交易相关链码接口
│   └── type_struct.go                   # 类型结构定义
├── global/                              # 全局变量目录
│   └── global.go                        # 全局变量定义文件
├── ipfs/                                # IPFS接口目录
│   └── ipfs.go                          # IPFS核心功能接口
├── ipfs_download/                       # IPFS文件下载目录
├── server/                              # HTTP服务器目录
│   ├── handler/                         # 请求处理器目录
│   │   ├── block_handler.go            # 区块查询处理器
│   │   ├── comment_handler.go          # 评论处理器
│   │   ├── contract_handler.go         # 智能合约处理器
│   │   ├── cti_handler.go             # 威胁情报处理器
│   │   ├── data_stat_handler.go       # 数据统计处理器
│   │   ├── incentive_handler.go       # 激励机制处理器
│   │   ├── ipfs_handler.go            # IPFS文件处理器
│   │   ├── kp_handler.go              # 知识图谱处理器
│   │   ├── model_handler.go           # 模型处理器
│   │   ├── tx_handler.go              # 交易处理器
│   │   └── user_handler.go            # 用户处理器
│   └── router.go                       # 路由配置文件
├── service/                            # 业务逻辑层目录
│   ├── ioc_example/                    # IOC样例目录
│   │   ├── example_statistic_info.json # 样例统计信息
│   │   └── statistic_IPFS_hash.txt    # IPFS哈希统计
│   ├── ipfs_service.go                # IPFS服务实现
│   ├── kp_service.go                  # 知识图谱服务实现
│   └── msgstruct.go                   # 消息结构定义
├── stats_db/                          # 统计数据存储目录
│   ├── attack_ioc_info.json          # 攻击IOC信息统计
│   ├── attack_type_statistics.json   # 攻击类型统计
│   ├── ioc_type_statistics.json     # IOC类型统计
│   ├── ioc_world_map_statistics.json # IOC世界地图统计
│   ├── traffic_time_series.json     # 流量时间序列统计
│   └── traffic_type_stats.json      # 流量类型统计
└── util/                             # 工具函数目录
    ├── file.go                       # 文件操作工具
    └── util.go                       # 通用工具函数
```

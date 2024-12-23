# bc-server端开发说明
## 1.说明
br-cti fabric接口服务
每次安装链码重启网络之后都需要复制ordererOrganizations和peerOrganizations文件夹
到conifg/server/目录下
## 2.编译
```
go build -o bc-server.exe
```
## 运行
```
go run main.go
```

## 3.调试接口
### 区块链相关接口
curl -X ANY http://localhost:7777/blockchain/queryBlock/{blockID} 
curl -X ANY http://localhost:7777/blockchain/queryChain

### 合约相关接口
curl -X POST http://localhost:7777/contract/queryContract
curl -X POST http://localhost:7777/contract/invokeContract

### 交易相关接口
curl -X POST http://localhost:7777/tx/getTransactionNonce

### 用户相关接口
curl -X POST http://localhost:7777/user/registerUserAccount
curl -X POST http://localhost:7777/user/queryUserInfo
curl -X POST http://localhost:7777/user/getUserStatistics
curl -X POST http://localhost:7777/user/queryPointTransactions

### CTI相关接口
curl -X POST http://localhost:7777/cti/registerCtiInfo
curl -X POST http://localhost:7777/cti/queryCtiInfo
curl -X POST http://localhost:7777/cti/queryCtiInfoByTypeWithPagination
curl -X POST http://localhost:7777/cti/queryCtiInfoByType
curl -X POST http://localhost:7777/cti/queryAllCtiInfoWithPagination
curl -X POST http://localhost:7777/cti/queryCtiInfoByCTIHash
curl -X POST http://localhost:7777/cti/queryCtiInfoByCreatorUserID

### 模型相关接口
curl -X POST http://localhost:7777/model/registerModelInfo
curl -X POST http://localhost:7777/model/queryModelInfo
curl -X POST http://localhost:7777/model/queryModelInfoByIDWithPagination
curl -X POST http://localhost:7777/model/queryModelsByTrafficType
curl -X POST http://localhost:7777/model/queryModelsByRefCTIId
curl -X POST http://localhost:7777/model/queryModelInfoByCreatorUserID

### 数据统计相关接口
curl -X POST http://localhost:7777/dataStat/queryCTISummaryInfo
curl -X POST http://localhost:7777/dataStat/getDataStatistics
curl -X POST http://localhost:7777/dataStat/getCTITrafficTrend
curl -X POST http://localhost:7777/dataStat/getAttackTypeRanking
curl -X POST http://localhost:7777/dataStat/getIOCsDistribution
curl -X POST http://localhost:7777/dataStat/getGlobalIOCsDistribution
curl -X POST http://localhost:7777/dataStat/getSystemOverview

### 知识平面相关(kp)接口
curl -X POST http://localhost:7777/kp/queryIOCGeoDistribution # 查询IOC地理分布
curl -X POST http://localhost:7777/kp/queryIOCTypeDistribution # 查询IOC类型分布
curl -X POST http://localhost:7777/kp/queryAttackTypeStatistics # 查询攻击类型统计
curl -X POST http://localhost:7777/kp/queryAttackIOCInfo # 查询攻击IOC信息

## 4.运行说明
在服务器运行时可使用端口转发
```shell
ssh -L 7777:localhost:7777 dev01@172.22.232.42
```
## 5.IPFS相关
```shell
ssh -L 5001:localhost:5001 dev01@172.22.232.42
ssh -L 8080:127.0.0.1:8080 dev01@172.22.232.42
```
### 获取IPFS内容
curl -X POST http://localhost:7777/ipfs/getIPFSContent
curl -X get http://localhost:7777/ipfs/getIPFSContent?hash=Qmc4R6bLJoRjYjGaU6cu1CbJqJS5avhecEHT2fBXxP343U
直链
http://127.0.0.1:8080/ipfs/Qmc4R6bLJoRjYjGaU6cu1CbJqJS5avhecEHT2fBXxP343U

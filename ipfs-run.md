## IPFS运行脚本
在用户根目录下启动
```bash
# 检查并终止使用5001端口的任何现有IPFS进程
echo "Checking for existing IPFS processes using port 5001..."
fuser -k 5001/tcp

# 启动IPFS守护程序并在后台运行
echo "Starting IPFS daemon..."
nohup ipfs daemon &

# 输出守护程序的PID到文件
echo $! > ipfs_daemon.pid

echo "IPFS started and running in background."

```

使用SSH端口转发访问
```shell
ssh -L 5001:localhost:5001 dev01@172.22.232.42
ssh -L 8080:127.0.0.1:8080 dev01@172.22.232.42
```
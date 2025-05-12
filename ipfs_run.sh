#!/bin/bash

# 检查并终止使用5001端口的任何现有IPFS进程
echo "Checking for existing IPFS processes using port 5001..."
fuser -k 5001/tcp

# 启动IPFS守护程序并在后台运行
echo "Starting IPFS daemon..."
nohup ipfs daemon &

# 输出守护程序的PID到文件
echo $! > ipfs_daemon.pid

echo "IPFS started and running in background."

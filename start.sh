#!/bin/bash

# QQ-like Server 启动脚本
# 使用方法：./start.sh

# 加载 .env 配置文件
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | awk '/=/ {print $1}')
fi

# 设置默认值
export SERVER_HOST=${SERVER_HOST:-"0.0.0.0"}
export SERVER_PORT=${SERVER_PORT:-8080}
export SERVER_MODE=${SERVER_MODE:-"release"}
export TCP_HOST=${TCP_HOST:-"0.0.0.0"}
export TCP_PORT=${TCP_PORT:-9090}
export TOKEN_SECRET=${TOKEN_SECRET:-"qq-like-server-secret-key-2024"}
export TOKEN_EXPIRE_HOUR=${TOKEN_EXPIRE_HOUR:-168}
export APP_SECRET=${APP_SECRET:-"qq-like-server-app-secret-2024"}

# 启动服务
./server

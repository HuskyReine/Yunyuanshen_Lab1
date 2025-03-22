# 使用 Golang 官方镜像作为构建基础
FROM golang:1.21.13-bullseye

# 设置工作目录
WORKDIR /app/cart

# 将当前目录的所有文件复制到容器的 /app/cart 目录下
COPY . /app/cart

# 配置 GOPROXY ，以从国内的代理源拉取依赖
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/

# 使用 bash 脚本编译服务
RUN bash build.sh

# 设置容器启动时执行的命令
CMD ["bash", "output/bootstrap.sh"]
# 使用 Golang 官方镜像作为构建基础
FROM golang:1.21.13-bullseye

# 将当前目录的所有文件复制到容器的 /app/order 目录下
COPY . /app/order

# 设置工作目录
WORKDIR /app/order

# 使用 bash 脚本编译服务
RUN bash build.sh

# 设置容器启动时执行的命令
CMD ["bash", "output/bootstrap.sh"]
# FRE Engine

Function Runtime Environment Engine


## 测试

1. 构建测试环境

部分系统未安装 python3 或者缺少相关的库，所以这里使用 docker python:3.7 环境来构建基础环境

构建镜像:

```bash
cd docker_image && docker build -t python/free/bench:3.7 .
```

2. 构建可执行文件

```bash
docker run -it --rm -v "$PWD":/go/src golang:1.14 bash -c "cd /go/src && go build -o free main.go" &&
scp free root@server:/free/ &&
rm -f free
```

### 6.1 

1. 顺序启动


2. 并发启动

```bash
clear && docker run --rm -it -v /free:/free --privileged python/free/bench:3.7 \
    bash -c "cd /free && rm -rf workspace/container/* && ./free -runtime base -p -n 32"

clear && docker run --rm -it -v /free:/free --privileged python/free/bench:3.7 \
    bash -c "cd /free && rm -rf workspace/container/* && ./free -runtime base -p -n 64"

clear && docker run --rm -it -v /free:/free --privileged python/free/bench:3.7 \
    bash -c "cd /free && rm -rf workspace/container/* && ./free -runtime base -p -n 128"

clear && docker run --rm -it -v /free:/free --privileged python/free/bench:3.7 \
    bash -c "cd /free && rm -rf workspace/container/* && ./free -runtime base -p -n 256"
``` 

### 6.2

1. 顺序启动

```bash 顺序无 zygote
clear && docker run --rm -it -v /free:/free --privileged python/free/bench:3.7 bash -c \
"cd /free && rm -rf workspace/container/* && ./free -runtime python3.7 -template normal -n 32" 64 128 256
``` 

```bash 顺序有 zygote
clear && docker run --rm -it -v /free:/free --privileged python/free/bench:3.7 bash -c \
"cd /free && rm -rf workspace/container/* && ./free -runtime python3.7 -template normal -zygote -n 32" 64 128 256
``` 

2. 
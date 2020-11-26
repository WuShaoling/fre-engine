# FRE Engine

Function Runtime Environment Engine


## 相关概念

- Runtime: 运行时环境，基础环境
- Template: 函数模版
- Container: 基于函数模版运行的容器(函数实例)


## 启动测试环境
docker run -v $PWD:/go/src -it --privileged golang:1.14 bash
docker run -v $PWD:/root -p 8080:80 -it --privileged python:3.7 bash


## 构建

- 本地直接构建：go build -o free main.go
- 基于 docker: docker run -it --rm -v "$PWD":/go/src golang:1.14 bash -c "cd /go/src && go build -o free main.go"

## 运行

运行前的准备

#### 1.构建

``` 
docker run -it --rm -v "$PWD":/go/src golang:1.14 bash -c "cd /go/src && go build -o free main.go"
```

#### 2. 准备 workspace 目录，可使用本项目中的 workspace 目录（复用 metadata 目录下的元数据）

#### 3. 构建 runtime rootfs

```
git clone https://github.com/WuShaoling/fre-runtime.git
cd fre-runtime/python3.7 && sh build.sh ${workspace目录}
```

#### 4. 运行环境安装 python3，安装 scipy numpy pandas django matplotlib 等测试包，可使用命令

`pip3 install -i http://pypi.douban.com/simple --trusted-host pypi.douban.com scipy numpy pandas django matplotlib`

或使用 docker 启动运行环境

```bash
docker run -it -v $PWD:/free --privileged -p 80:80 python:3.7 bash
# pip3 install -i http://pypi.douban.com/simple --trusted-host pypi.douban.com scipy numpy pandas django matplotlib
# cd /free
```

#### 5. 启动

```bash
./free
```

#### 6. 执行测试函数

```bash
curl -X POST \
  'http://localhost:80/api/container/normal?sync=true' \
  -H 'content-type: application/json' \
  -d '{
	"key1": "hello",
	"key2": "world"
}'
```

docker restart 0ef83d3a156a && docker exec -it 0ef83d3a156a bash -c "cd /free && rm -rf workspace/container/* && ./free"
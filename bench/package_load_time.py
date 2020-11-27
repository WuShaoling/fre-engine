import importlib
import time

''' 测试包加载的时间

执行流程

1. 使用 docker python:3.7 环境
docker run -it -v $PWD:/root python:3.7 bash

2. 安装依赖包
pip3 install -i http://pypi.douban.com/simple --trusted-hostpypi.douban.com \
    django flask numpy pandas matplotlib setuptools requests sqlalchemy

3. 执行测试脚本
python /root/package_load_time.py

'''

packages = ["django", "flask", "numpy", "pandas", "matplotlib", "setuptools", "requests", "sqlalchemy"]
for package in packages:
    t1 = int(round(time.time() * 1000000 / 1e3))
    importlib.import_module(package)
    t2 = int(round(time.time() * 1000000 / 1e3))
    print(package + ": ", t2 - t1)

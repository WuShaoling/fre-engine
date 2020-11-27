import csv
import math

import matplotlib.pyplot as plt
import numpy as np

fz = 10


def read_fre_data(file):
    csv_file = open(file, "r")
    reader = csv.reader(csv_file)
    data = []
    for item in reader:
        v = int(int(item[3]) / 1000)
        data.append(np.log2(v))
    data.sort()
    return data


def read_docker_data(file):
    csv_file = open(file, "r")
    reader = csv.reader(csv_file)
    data = []
    for item in reader:
        v = int(item[3])
        data.append(np.log2(v))
    data.sort()
    return data


def draw_cdf():
    # fre
    end = []
    fre_label = 'FRE'
    for file in file_set_fre:
        data = read_fre_data(file)
        hist, bin_edge = np.histogram(data)
        cdf = np.cumsum(hist / sum(hist))
        plt.plot(cdf, bin_edge[1:], '-', color='k', label=fre_label, linewidth=1.3)
        end.append(bin_edge[-1])
        fre_label = ''

    # docker
    docker_label = 'Docker'
    for file in file_set_docker:
        data = read_docker_data(file)
        hist, bin_edge = np.histogram(data)
        cdf = np.cumsum(hist / sum(hist))
        plt.plot(cdf, bin_edge[1:], '--', color='k', label=docker_label, linewidth=1.3)
        end.append(bin_edge[-1])
        docker_label = ''

    # 标注
    text = ['■', '●', '◆'] * 2
    for i in range(0, len(end)):
        plt.text(1, end[i] - 0.15, text[i], fontsize=fz)
    plt.text(0.2, 10, "■ 64 concurrent  ● 128 concurrent  ◆ 256 concurrent", fontsize=fz)

    # 坐标
    plt.xlim([0, 1.01])
    plt.ylim([0, math.ceil(end[-1])])
    plt.xlabel('Percent of Operations', fontsize=fz)
    plt.ylabel('Latency after log2 (ms)', fontsize=fz)
    # 图例
    plt.legend(loc='lower right', fontsize=fz)
    # 设置边框
    ax = plt.gca()
    ax.spines['right'].set_visible(False)
    ax.spines['top'].set_visible(False)
    # 保存

    plt.tick_params(labelsize=fz)
    plt.savefig("output/create_parallel.png", dpi=300, bbox_inches='tight')
    plt.close()


file_set_docker = [
    "parallel_data_set/docker-64.csv",
    "parallel_data_set/docker-128.csv",
    "parallel_data_set/docker-256.csv"]

file_set_fre = [
    "parallel_data_set/fre-64.csv",
    "parallel_data_set/fre-128.csv",
    "parallel_data_set/fre-256.csv"]

draw_cdf()

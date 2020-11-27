import csv

import matplotlib.pyplot as plt

fz = 10


def read_data(file):
    csv_file = open(file, "r", encoding='utf-8')
    reader = csv.reader(csv_file)
    reader = list(map(list, zip(*reader)))

    for i in range(0, len(reader)):
        for j in range(0, len(reader[i])):
            reader[i][j] = int(reader[i][j].encode('utf-8').decode('utf-8-sig').strip())
    return reader


def draw_cgroup_512():
    plt.subplot(1, 2, 1)
    matrix = read_data("data_set/cgroup-512.csv")

    for i in range(1, len(matrix)):
        for j in range(0, len(matrix[i])):
            matrix[i][j] = matrix[i][j] / 1000

    plt.plot(matrix[0], matrix[1], ':', color='k', label="create", linewidth=0.9)
    plt.plot(matrix[0], matrix[2], '-', color='k', label="reuse", linewidth=0.9)

    plt.legend(loc='best', fontsize=fz)
    plt.xlabel('(a) Concurrent Operations(512)', fontsize=fz)
    plt.ylabel('Latency(s)', fontsize=fz)
    plt.xticks([0, 128, 256, 384, 512])
    plt.tick_params(labelsize=fz)


def draw_cgroup_1024():
    plt.subplot(1, 2, 2)
    matrix = read_data("data_set/cgroup-1024.csv")

    for i in range(1, len(matrix)):
        for j in range(0, len(matrix[i])):
            matrix[i][j] = matrix[i][j] / 1000

    plt.plot(matrix[0], matrix[1], ':', color='k', label="create", linewidth=0.9)
    plt.plot(matrix[0], matrix[2], '-', color='k', label="reuse", linewidth=0.9)

    plt.legend(loc='best', fontsize=fz)
    plt.xlabel('(b) Concurrent Operations(1024)', fontsize=fz)
    plt.ylabel('Latency(s)', fontsize=fz)
    plt.xticks([0, 256, 512, 768, 1024])
    plt.tick_params(labelsize=fz)


plt.figure(figsize=(10, 3))
draw_cgroup_512()
draw_cgroup_1024()

plt.savefig("output/cgroup.png", dpi=300, bbox_inches='tight')
plt.close()

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


def draw_ns_512():
    plt.subplot(1, 2, 1)

    matrix = read_data("data_set/ns-512.csv")

    # for i in range(1, len(matrix)):
    #     for j in range(0, len(matrix[i])):
    #         matrix[i][j] = matrix[i][j] / 1000

    plt.plot(matrix[0], matrix[3], ':', color='k', label="No namespace isolation", linewidth=0.7)
    plt.plot(matrix[0], matrix[4], '-', color='k', label="PID & IPC & UTS & Mount & User", linewidth=0.7)
    plt.plot(matrix[0], matrix[2], '--', color='k', label="PID & IPC & UTS & Mount & User & Network",
             linewidth=0.7)

    plt.legend(loc='best', fontsize=8)
    plt.xlabel('(a) Concurrent Operations(512)', fontsize=fz)
    plt.ylabel('Latency(ms)', fontsize=fz)
    plt.xticks([0, 128, 256, 384, 512])
    plt.tick_params(labelsize=fz)


def draw_ns_1024():
    plt.subplot(1, 2, 2)
    matrix = read_data("data_set/ns-1024.csv")

    for i in range(1, len(matrix)):
        for j in range(0, len(matrix[i])):
            matrix[i][j] = matrix[i][j] ** 0.5
            # matrix[i][j] = np.log2(matrix[i][j])
            # matrix[i][j] = matrix[i][j] / 1000

    plt.plot(matrix[0], matrix[1], ':', color='k', label="No namespace isolation", linewidth=0.7)
    plt.plot(matrix[0], matrix[3], '-', color='k', label="PID & IPC & UTS & Mount & User", linewidth=0.7)
    plt.plot(matrix[0], matrix[2], '--', color='k', label="PID & IPC & UTS & Mount & User & Network",
             linewidth=0.7)

    plt.legend(loc='best', fontsize=8)
    plt.xlabel('(b) Concurrent Operations(1024)', fontsize=fz)
    plt.ylabel('Latency after square (ms)', fontsize=fz)
    plt.xticks([0, 256, 512, 768, 1024])
    plt.tick_params(labelsize=fz)


plt.figure(figsize=(10.5, 3.5))
draw_ns_512()
draw_ns_1024()
plt.savefig("output/ns.png", dpi=300, bbox_inches='tight')
plt.close()

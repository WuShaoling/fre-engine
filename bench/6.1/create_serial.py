import matplotlib.pyplot as plt

label_list = ['64', '128', '256', '512']  # 横坐标刻度显示值
list_docker = [427, 439, 452, 476]  # 纵坐标值1
list_fre = [18, 20, 21, 23]  # 纵坐标值2
x_label = "Concurrent Operations"
y_label = "Startup time(ms)"
output = "output/create_serial.png"
fz = 10
x = range(len(list_docker))

rects1 = plt.bar(x=x, height=list_docker, width=0.3, alpha=0.8, color='#BEBEBE', edgecolor='k', label="Docker")
rects2 = plt.bar(x=[i + 0.3 for i in x], height=list_fre, width=0.3, color='w', edgecolor='k', label="FRE")

# 编辑文本
for rect in rects1:
    height = rect.get_height()
    plt.text(rect.get_x() + rect.get_width() / 2, height + 1, str(height), ha="center", va="bottom", fontsize=fz)
for rect in rects2:
    height = rect.get_height()
    plt.text(rect.get_x() + rect.get_width() / 2, height + 1, str(height), ha="center", va="bottom", fontsize=fz)

plt.xlabel(x_label, fontsize=fz)
plt.ylabel(y_label, fontsize=fz)
plt.xticks([index + 0.15 for index in x], label_list)
plt.ylim(0, 600)  # y轴取值范围
plt.legend(fontsize=fz)  # 设置题注
plt.tick_params(labelsize=fz)

plt.savefig(output, dpi=300, bbox_inches='tight')
plt.close()

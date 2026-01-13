import matplotlib.pyplot as plt
import csv
from datetime import datetime

def readFile(filePath):
    dateFormat = "%Y-%m-%dT%H:%M:%SZ"
    times = []
    values = []

    with open(filePath, 'r') as f:
        reader = csv.reader(f)

        for row in reader:
            values.append(float(row[0]))
            times.append(datetime.strptime(row[1], dateFormat))
    
    return times, values

def drawPlot(axes, index, title, xLabel, x, yLabel, y):
    axes[index].plot(x, y, marker="o")
    axes[index].set_title(title)
    axes[index].set_xlabel(xLabel)
    axes[index].set_ylabel(yLabel)

# Load Data
co2X, co2Y = readFile("../serial-to-file/co2.csv")
humidityX, humidityY = readFile("../serial-to-file/humidity.csv")
tempX, tempY = readFile("../serial-to-file/temperature.csv")

# Plot
fig, ax = plt.subplots(1, 3, figsize=(32, 6))

drawPlot(ax, 0, "CO2", "UTC", co2X, "ppm", co2Y)
drawPlot(ax, 1, "Humidity", "UTC", humidityX, "%", humidityY)
drawPlot(ax, 2, "Temperature", "UTC", tempX, "Celcius", tempY)

plt.show()

print("finish")


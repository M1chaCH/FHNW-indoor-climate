basic.forever(() => {
    const co2 = grove.readDataFromSCD30(grove.SCD30DataType.CO2)
    if (!Number.isNaN(co2)) {
        serial.writeValue("co2", co2)
    }

    const humidity = grove.readDataFromSCD30(grove.SCD30DataType.Humidity)
    if (!Number.isNaN(humidity)) {
        serial.writeValue("humidity", humidity)
    }

    const temperature = grove.readDataFromSCD30(grove.SCD30DataType.CelsiusTemperature)
    if (!Number.isNaN(temperature)) {
        serial.writeValue("temperature", temperature)
    }
})
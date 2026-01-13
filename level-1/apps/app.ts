let isEven = false;

basic.forever(() => {
    const data = {
        co2: grove.readDataFromSCD30(grove.SCD30DataType.CO2),
        humidity: grove.readDataFromSCD30(grove.SCD30DataType.Humidity),
        temperature: grove.readDataFromSCD30(grove.SCD30DataType.CelsiusTemperature),
    }
    
    serial.writeString(JSON.stringify(data))
    basic.showIcon(isEven ? IconNames.QuarterNote : IconNames.EighthNote)
    isEven = !isEven
})
